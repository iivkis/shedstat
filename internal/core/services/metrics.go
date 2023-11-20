package services

import (
	"context"
	"database/sql"
	"errors"
	"shedstat/internal/adapters/repository"
	"shedstat/internal/core/domain"
	shedevrumapi "shedstat/pkg/shedevrum-api"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/exp/slog"
)

type MetricsService struct {
	logger               *slog.Logger
	repoMetrics          *repository.MetricsClickHouseRepository
	repoMetricsCollector *repository.MetricsCollectorPostgresRepository
	svcProfile           *ProfileService
	shedAPI              *shedevrumapi.ShedevrumAPI
}

func NewMetricsService(
	logger *slog.Logger,
	repoMetrics *repository.MetricsClickHouseRepository,
	repoMetricsCollector *repository.MetricsCollectorPostgresRepository,
	svcProfile *ProfileService,
	shedAPI *shedevrumapi.ShedevrumAPI,
) *MetricsService {
	svc := &MetricsService{
		logger:               logger,
		repoMetrics:          repoMetrics,
		repoMetricsCollector: repoMetricsCollector,
		svcProfile:           svcProfile,
		shedAPI:              shedAPI,
	}
	go svc.sheduler()
	return svc
}

func (s *MetricsService) sheduler() {
	const op = "services.MetricsService.sheduler"
	for range time.Tick(time.Second) {
		lastMetricShedule, err := s.repoMetricsCollector.GetLast(context.Background())
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				s.logger.Error(err.Error(), "op", op)
				continue
			}
		}
		if errors.Is(err, sql.ErrNoRows) || lastMetricShedule.CreatedAt.Add(time.Hour*8).Before(time.Now()) {
			metricShedule, err := s.collectSocialStats(context.Background())
			if err != nil {
				s.logger.Error(err.Error(), "op", op)
			} else {
				if err := s.repoMetricsCollector.Create(context.Background(), metricShedule); err != nil {
					s.logger.Error(err.Error(), "op", op)
				}
			}
		}
	}
}

func (s *MetricsService) collectSocialStats(ctx context.Context) (*domain.MetricsCollectorEntity, error) {
	const op = "services.MetricsService.collectSocialStats"
	s.logger.Info("run_social_stats_collector", "op", op)

	collectorPullSize := 50

	var (
		queue             = make(chan struct{}, 30)
		startFromID       int
		socialStats       []*domain.MetricsEntity
		metricSheduleStat domain.MetricsCollectorEntity
	)
	defer close(queue)

	for {
		socialStats = make([]*domain.MetricsEntity, 0, collectorPullSize)

		profiles, err := s.svcProfile.GetList(ctx, startFromID, collectorPullSize)
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
			metricSheduleStat.ProfileHandledTotal += 1

			go func(p *domain.ProfileEnity) {
				defer func() {
					<-queue
					wg.Done()
				}()

				queue <- struct{}{}

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
		if err := s.repoMetrics.Create(ctx, socialStats); err != nil {
			s.logger.Error(err.Error(), "op", op)
			return nil, err
		}
	}

	return &metricSheduleStat, nil
}
