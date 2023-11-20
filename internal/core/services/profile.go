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
	logger  *slog.Logger
	repo    *repository.ProfilePostgresRepository
	shedAPI *shedevrumapi.ShedevrumAPI
}

func NewProfileService(
	logger *slog.Logger,
	repo *repository.ProfilePostgresRepository,
	shedAPI *shedevrumapi.ShedevrumAPI,
) *ProfileService {
	svc := &ProfileService{
		logger:  logger,
		repo:    repo,
		shedAPI: shedAPI,
	}
	go svc.sheduler()
	return svc
}

func (s *ProfileService) sheduler() {
	for {
		s.collectProfileFromTopDay()
		time.Sleep(time.Hour * 24)
	}
}

func (s *ProfileService) collectProfileFromTopDay() {
	const op = "services.ProfileService.collectProfileFromTopDay"
	for startFrom := ""; ; {
		feed, err := s.shedAPI.Feed.GetTop(shedevrumapi.FEED_TOP_PERIOD_DAY, 100, startFrom)
		if err != nil {
			s.logger.Error(err.Error(), "op", op)
			break
		}
		for _, post := range feed.Posts {
			if exists, _ := s.repo.ExistsByShedevrumID(context.Background(), post.User.ID); !exists {
				if err := s.repo.Create(context.Background(), &domain.ProfileEnity{
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
	return s.repo.Create(ctx, profile)
}

func (s *ProfileService) Get(ctx context.Context, id int) (*domain.ProfileEnity, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProfileService) GetList(ctx context.Context, startFromID int, amount int) ([]*domain.ProfileEnity, error) {
	return s.repo.GetList(ctx, startFromID, amount)
}
