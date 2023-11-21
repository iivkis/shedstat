package services

import (
	"context"
	"fmt"
	"shedstat/internal/adapters/repository"
	"shedstat/internal/core/domain"
	shedevrumapi "shedstat/pkg/shedevrum-api"
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/exp/slog"
)

type ProfileService struct {
	logger               *slog.Logger
	repoProfile          *repository.ProfilePostgresRepository
	repoProfileCollector *repository.ProfileCollectorPostgresRepository
	repoMetrics          *repository.MetricsClickHouseRepository
	shedAPI              *shedevrumapi.ShedevrumAPI
	cache                *cache.Cache
}

func NewProfileService(
	logger *slog.Logger,
	repoProfile *repository.ProfilePostgresRepository,
	repoProfileCollector *repository.ProfileCollectorPostgresRepository,
	repoMetrics *repository.MetricsClickHouseRepository,
	shedAPI *shedevrumapi.ShedevrumAPI,
) *ProfileService {
	svc := &ProfileService{
		logger:               logger,
		repoProfile:          repoProfile,
		repoProfileCollector: repoProfileCollector,
		repoMetrics:          repoMetrics,
		shedAPI:              shedAPI,
		cache:                cache.New(time.Hour, time.Minute*30),
	}
	go svc.shedulerFeedTopDay()
	return svc
}

func (s *ProfileService) shedulerFeedTopDay() {
	for {
		profileCollectorFeedTopDay, err := s.repoProfileCollector.GetLastCollectedAt(domain.PROFILE_COLLECTOR_COLLECTOR_TYPE_FEED_TOP_DAY)
		if err != nil {
			s.logger.Error(err.Error())
			continue
		}
		if !profileCollectorFeedTopDay.Valid || profileCollectorFeedTopDay.Time.Add(time.Hour*24).Before(time.Now()) {
			s.collectProfileFromFeedTopDay()
			if err := s.repoProfileCollector.UpdateLastCollectedAt(domain.PROFILE_COLLECTOR_COLLECTOR_TYPE_FEED_TOP_DAY); err != nil {
				s.logger.Error(err.Error())
			}
		}

		time.Sleep(time.Minute * 10)
	}
}

func (s *ProfileService) collectProfileFromFeedTopDay() {
	const op = "services.ProfileService.collectProfileFromTopDay"
	s.logger.Info("run_profile_from_feed_top_day_collector", "op", op)
	for startFrom := ""; ; {
		feed, err := s.shedAPI.Feed.GetTop(shedevrumapi.FEED_TOP_PERIOD_DAY, 100, startFrom)
		if err != nil {
			s.logger.Error(err.Error(), "op", op)
			break
		}
		for _, post := range feed.Posts {
			s.logger.Info("check_if_profile_exists", "op", op, "shedevrum_id", post.User.ID)
			if exists, _ := s.repoProfile.ExistsByShedevrumID(context.Background(), post.User.ID); !exists {
				s.logger.Info("add_new_profile", "op", op, "shedevrum_id", post.User.ID)
				if err := s.repoProfile.Create(context.Background(), post.User.ID); err != nil {
					s.logger.Error(err.Error(), "op", op)
				}
			}
		}
		if startFrom = feed.Next; startFrom == "" {
			break
		}
	}
}

func (s *ProfileService) Create(ctx context.Context, shedevrumID string) error {
	return s.repoProfile.Create(ctx, shedevrumID)
}

func (s *ProfileService) Get(ctx context.Context, id int) (*domain.ProfileEnity, error) {
	return s.repoProfile.GetByID(ctx, id)
}

func (s *ProfileService) GetByShedevrumID(ctx context.Context, shedevrumID string) (*domain.ProfileEnity, error) {
	profile, err := s.repoProfile.GetByShedevrumID(ctx, shedevrumID)
	if err != nil {
		return nil, err
	}

	cachedProfileID := fmt.Sprintf("profile_%s", profile.ShedevrumID)

	if cachedProfile, ok := s.cache.Get(cachedProfileID); ok {
		p := cachedProfile.(*domain.ProfileEnity)
		profile.AvatarURL = p.AvatarURL
		profile.Name = p.Name
	} else {
		remoteProfile, err := s.shedAPI.Users.GetFeed(shedevrumID, 0, "")
		if err != nil {
			return nil, err
		}
		profile.Name = remoteProfile.User.DisplayName
		profile.AvatarURL = remoteProfile.User.AvatartURL

		remoteSocialStats, err := s.shedAPI.Users.GetSocialStats(shedevrumID)
		if err != nil {
			return nil, err
		}
		profile.Subscriptions = remoteSocialStats.Subscriptions
		profile.Subscribers = remoteSocialStats.Subscribers
		profile.Likes = remoteSocialStats.Likes

		s.cache.Set(cachedProfileID, profile, 0)
	}

	return profile, nil
}

func (s *ProfileService) GetList(ctx context.Context, startFromID int, amount int) ([]*domain.ProfileEnity, error) {
	return s.repoProfile.GetList(ctx, startFromID, amount)
}

func (s *ProfileService) GetMetrics(ctx context.Context, shedevrumID string) ([]*domain.MetricsChartEntity, error) {
	return s.repoMetrics.GetByShedevrumID(shedevrumID)
}
