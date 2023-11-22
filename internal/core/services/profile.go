package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shedstat/internal/adapters/repository"
	"shedstat/internal/core/domain"
	shedevrumapi "shedstat/pkg/shedevrum-api"
	"sync"
	"sync/atomic"
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/exp/slog"
)

type ProfileService struct {
	logger                      *slog.Logger
	repoProfile                 *repository.ProfilePostgresRepository
	repoProfileCollector        *repository.ProfileCollectorPostgresRepository
	repoProfileMetrics          *repository.MetricsClickHouseRepository
	repoProfileMetricsCollector *repository.MetricsCollectorPostgresRepository
	shedAPI                     *shedevrumapi.ShedevrumAPI
	cache                       *cache.Cache
}

func NewProfileService(
	logger *slog.Logger,
	repoProfile *repository.ProfilePostgresRepository,
	repoProfileCollector *repository.ProfileCollectorPostgresRepository,
	repoProfileMetrics *repository.MetricsClickHouseRepository,
	repoProfileMetricsCollector *repository.MetricsCollectorPostgresRepository,
	shedAPI *shedevrumapi.ShedevrumAPI,
) *ProfileService {
	svc := &ProfileService{
		logger:                      logger,
		repoProfile:                 repoProfile,
		repoProfileCollector:        repoProfileCollector,
		repoProfileMetrics:          repoProfileMetrics,
		repoProfileMetricsCollector: repoProfileMetricsCollector,
		shedAPI:                     shedAPI,
		cache:                       cache.New(time.Hour, time.Minute*30),
	}
	svc.runSchedulers()
	return svc
}

func (s *ProfileService) runSchedulers() {
	go s.profilesAssemblyScheduling()
	go s.profileMetricsAssemblyScheduling()
}

func (s *ProfileService) profileMetricsAssemblyScheduling() {
	const op = "services.MetricsService.sheduler"
	for {
		lastMetricShedule, err := s.repoProfileMetricsCollector.GetLast(context.Background())
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				s.logger.Error(err.Error(), "op", op)
				continue
			}
		}
		if errors.Is(err, sql.ErrNoRows) || lastMetricShedule.CreatedAt.Add(time.Hour*8).Before(time.Now()) {
			metricShedule, err := s.profileMetricsCollect(context.Background())
			if err != nil {
				s.logger.Error(err.Error(), "op", op)
			} else {
				if err := s.repoProfileMetricsCollector.Create(context.Background(), metricShedule); err != nil {
					s.logger.Error(err.Error(), "op", op)
				}
			}
		}

		time.Sleep(time.Minute * 1)
	}
}

func (s *ProfileService) profileMetricsCollect(ctx context.Context) (*domain.MetricsCollectorEntity, error) {
	const op = "services.MetricsService.collectSocialStats"
	const collectorPullSize = 100

	s.logger.Info("run_social_stats_collector", "op", op)

	var (
		queue             = make(chan struct{}, 10)
		startFromID       int
		socialStats       []*domain.MetricsEntity
		metricSheduleStat domain.MetricsCollectorEntity
	)
	defer close(queue)

	for {
		socialStats = make([]*domain.MetricsEntity, 0, collectorPullSize)

		profiles, err := s.repoProfile.GetList(ctx, startFromID, collectorPullSize)
		if err != nil {
			s.logger.Error(err.Error(), "op", op)
			return nil, err
		}
		if len(profiles) == 0 {
			break
		}
		startFromID = profiles[len(profiles)-1].ID

		wg := &sync.WaitGroup{}
		for _, p := range profiles {
			wg.Add(1)

			go func(p *domain.ProfileEnity) {
				defer func() {
					<-queue
					wg.Done()
				}()

				queue <- struct{}{}

				atomic.AddUint64(&metricSheduleStat.ProfileHandledTotal, 1)

				s.logger.Info("get_profile_social_stats", "op", op, "shedevrum_id", p.ShedevrumID, "collected", atomic.LoadUint64(&metricSheduleStat.ProfileHandledTotal))
				stat, err := s.shedAPI.Users.GetSocialStats(p.ShedevrumID)
				if err != nil {
					atomic.AddUint64(&metricSheduleStat.ProfileHandledBad, 1)
					s.logger.Error(err.Error(), "op", op)
				} else {
					atomic.AddUint64(&metricSheduleStat.ProfileHandledSuccess, 1)
					socialStats = append(socialStats, &domain.MetricsEntity{
						ProfileID:     p.ID,
						ShedevrumID:   p.ShedevrumID,
						Subscriptions: stat.Subscriptions,
						Subscribers:   stat.Subscribers,
						Likes:         stat.Likes,
					})
				}
			}(p)
		}
		wg.Wait()

		s.logger.Info("push_profile_social_stats", "op", op, "count", len(socialStats))
		if err := s.repoProfileMetrics.Create(ctx, socialStats); err != nil {
			s.logger.Error(err.Error(), "op", op)
			return nil, err
		}
	}

	return &metricSheduleStat, nil
}

func (s *ProfileService) profilesAssemblyScheduling() {
	for {
		profileCollectorFeedTopDay, err := s.repoProfileCollector.GetLastCollectedAt(domain.PROFILE_COLLECTOR_COLLECTOR_TYPE_FEED_TOP_DAY)
		if err != nil {
			s.logger.Error(err.Error())
			continue
		}
		if !profileCollectorFeedTopDay.Valid || profileCollectorFeedTopDay.Time.Add(time.Hour*24).Before(time.Now()) {
			s.topDayProfilesCollect()
			if err := s.repoProfileCollector.UpdateLastCollectedAt(domain.PROFILE_COLLECTOR_COLLECTOR_TYPE_FEED_TOP_DAY); err != nil {
				s.logger.Error(err.Error())
			}
		}

		time.Sleep(time.Minute * 1)
	}
}

func (s *ProfileService) topDayProfilesCollect() {
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
		profile.Name = p.Name
		profile.AvatarURL = p.AvatarURL
		profile.Link = p.Link
		profile.Subscriptions = p.Subscriptions
		profile.Subscribers = p.Subscribers
		profile.Likes = p.Likes
	} else {
		remoteProfile, err := s.shedAPI.Users.GetFeed(shedevrumID, 0, "")
		if err != nil {
			return nil, err
		}
		profile.Name = remoteProfile.User.DisplayName
		profile.AvatarURL = remoteProfile.User.AvatartURL
		profile.Link = remoteProfile.User.ShareLink

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
	return s.repoProfileMetrics.GetByShedevrumID(shedevrumID)
}

func (s *ProfileService) GetTop(ctx context.Context, filter domain.MetricsGetTopFilter, amount int) ([]*domain.ProfileEnity, error) {
	topList, err := s.repoProfileMetrics.GetTop(ctx, filter, amount)
	if err != nil {
		return nil, err
	}
	profiles := make([]*domain.ProfileEnity, 0, len(topList))
	for _, p := range topList {
		profile, err := s.GetByShedevrumID(ctx, p.ShedevrumID)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}
