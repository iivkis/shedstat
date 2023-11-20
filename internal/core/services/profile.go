package services

import (
	"context"
	"shedstat/internal/adapters/repository"
	"shedstat/internal/core/domain"
	shedevrumapi "shedstat/pkg/shedevrum-api"
	"time"

	"golang.org/x/exp/slog"
)

type ProfileService struct {
	logger               *slog.Logger
	repoProfile          *repository.ProfilePostgresRepository
	repoProfileCollector *repository.ProfileCollectorPostgresRepository
	shedAPI              *shedevrumapi.ShedevrumAPI
}

func NewProfileService(
	logger *slog.Logger,
	repoProfile *repository.ProfilePostgresRepository,
	repoProfileCollector *repository.ProfileCollectorPostgresRepository,
	shedAPI *shedevrumapi.ShedevrumAPI,
) *ProfileService {
	svc := &ProfileService{
		logger:               logger,
		repoProfile:          repoProfile,
		repoProfileCollector: repoProfileCollector,
		shedAPI:              shedAPI,
	}
	go svc.shedulerFeedTopDay()
	return svc
}

func (s *ProfileService) shedulerFeedTopDay() {
	for range time.Tick(time.Minute) {
		profileCollectorFeedTopDay, err := s.repoProfileCollector.GetLastCollectedAt(domain.PROFILE_COLLECTOR_COLLECTOR_TYPE_FEED_TOP_DAY)
		if err != nil {
			s.logger.Error(err.Error())
			continue
		}
		if !profileCollectorFeedTopDay.Valid || profileCollectorFeedTopDay.Time.Add(time.Hour*24).Before(time.Now()) {
			s.collectProfileFromFeedTopDay()
		}
	}
}

func (s *ProfileService) collectProfileFromFeedTopDay() {
	const op = "services.ProfileService.collectProfileFromTopDay"
	for startFrom := ""; ; {
		feed, err := s.shedAPI.Feed.GetTop(shedevrumapi.FEED_TOP_PERIOD_DAY, 100, startFrom)
		if err != nil {
			s.logger.Error(err.Error(), "op", op)
			break
		}
		for _, post := range feed.Posts {
			if exists, _ := s.repoProfile.ExistsByShedevrumID(context.Background(), post.User.ID); !exists {
				if err := s.repoProfile.Create(context.Background(), &domain.ProfileEnity{
					ShedevrumID: post.User.ID,
				}); err != nil {
					s.logger.Error(err.Error(), "op", op)
				}
			}
		}
		if startFrom = feed.Next; startFrom == "" {
			break
		}
	}
}

func (s *ProfileService) Create(ctx context.Context, profile *domain.ProfileEnity) error {
	return s.repoProfile.Create(ctx, profile)
}

func (s *ProfileService) Get(ctx context.Context, id int) (*domain.ProfileEnity, error) {
	return s.repoProfile.GetByID(ctx, id)
}

func (s *ProfileService) GetList(ctx context.Context, startFromID int, amount int) ([]*domain.ProfileEnity, error) {
	return s.repoProfile.GetList(ctx, startFromID, amount)
}
