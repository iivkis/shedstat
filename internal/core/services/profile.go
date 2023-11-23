package services

import (
	"context"
	"database/sql"
	"errors"
	"shedstat/internal/adapters/repository"
	"shedstat/internal/core/domain"
	shedevrumapi "shedstat/pkg/shedevrum-api"
	"strings"
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
	repoProfileMetrics          *repository.ProfileMetricsClickHouseRepository
	repoProfileMetricsCollector *repository.ProfileMetricsCollectorPostgresRepository
	shedAPI                     *shedevrumapi.ShedevrumAPI
	cache                       *cache.Cache
}

func NewProfileService(
	logger *slog.Logger,
	repoProfile *repository.ProfilePostgresRepository,
	repoProfileCollector *repository.ProfileCollectorPostgresRepository,
	repoProfileMetrics *repository.ProfileMetricsClickHouseRepository,
	repoProfileMetricsCollector *repository.ProfileMetricsCollectorPostgresRepository,
	shedAPI *shedevrumapi.ShedevrumAPI,
) *ProfileService {
	svc := &ProfileService{
		logger:                      logger,
		repoProfile:                 repoProfile,
		repoProfileCollector:        repoProfileCollector,
		repoProfileMetrics:          repoProfileMetrics,
		repoProfileMetricsCollector: repoProfileMetricsCollector,
		shedAPI:                     shedAPI,
		cache:                       cache.New(time.Hour, time.Minute),
	}
	svc.runSchedulers()
	return svc
}

func (s *ProfileService) runSchedulers() {
	go s.profilesAssemblyScheduling()
	go s.profileMetricsAssemblyScheduling()
}

func (s *ProfileService) profileMetricsAssemblyScheduling() {
	const op = "services.ProfileService.profileMetricsAssemblyScheduling"
	for {
		lastCollection, err := s.repoProfileMetricsCollector.GetLast(context.Background())
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				s.logger.Error(err.Error(), "op", op)
				continue
			}
		}
		if errors.Is(err, sql.ErrNoRows) || lastCollection.CreatedAt.Add(time.Hour*12).Before(time.Now()) {
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

func (s *ProfileService) profileMetricsCollect(ctx context.Context) (*domain.ProfileMetricsCollectorEntity, error) {
	const op = "services.ProfileService.profileMetricsCollect"
	const pullSize = 100

	s.logger.Info("run_social_stats_collector", "op", op)

	var (
		queue          = make(chan struct{}, 15)
		startFromID    uint64
		socialStats    []*domain.ProfileMetricsEntity
		collectorStats domain.ProfileMetricsCollectorEntity
	)
	defer close(queue)

	for {
		socialStats = make([]*domain.ProfileMetricsEntity, 0, pullSize)

		profiles, err := s.repoProfile.GetList(ctx, startFromID, pullSize)
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

				atomic.AddUint64(&collectorStats.ProfileHandledTotal, 1)

				s.logger.Info("get_profile_social_stats", "op", op, "shedevrum_id", p.ShedevrumID, "collected", atomic.LoadUint64(&collectorStats.ProfileHandledTotal))
				stat, err := s.shedAPI.Users.GetSocialStats(p.ShedevrumID)
				if err != nil {
					atomic.AddUint64(&collectorStats.ProfileHandledBad, 1)
					s.logger.Error(err.Error(), "op", op)
				} else {
					atomic.AddUint64(&collectorStats.ProfileHandledSuccess, 1)
					socialStats = append(socialStats, &domain.ProfileMetricsEntity{
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

	return &collectorStats, nil
}

func (s *ProfileService) profilesAssemblyScheduling() {
	const op = "services.ProfileService.profilesAssemblyScheduling"
	for {
		profileCollectorFeedTopDay, err := s.repoProfileCollector.GetLastCollectedAt(domain.PROFILE_COLLECTOR_COLLECTOR_TYPE_FEED_TOP_DAY)
		if err != nil {
			s.logger.Error(err.Error(), "op", op)
			continue
		}
		if !profileCollectorFeedTopDay.Valid || profileCollectorFeedTopDay.Time.Add(time.Hour*24).Before(time.Now()) {
			s.topDayProfilesCollect()
			if err := s.repoProfileCollector.UpdateLastCollectedAt(domain.PROFILE_COLLECTOR_COLLECTOR_TYPE_FEED_TOP_DAY); err != nil {
				s.logger.Error(err.Error(), "op", op)
			}
		}
		time.Sleep(time.Minute * 1)
	}
}

func (s *ProfileService) topDayProfilesCollect() {
	const op = "services.ProfileService.topDayProfilesCollect"
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

func (s *ProfileService) cacheSetRemoteUser(user *shedevrumapi.UserEntity) {
	s.cache.Set("remote_user_"+user.ID, user, cache.DefaultExpiration)
}

func (s *ProfileService) cacheGetRemoteUser(userID string) (*shedevrumapi.UserEntity, bool) {
	if u, ok := s.cache.Get("remote_user_" + userID); ok {
		return u.(*shedevrumapi.UserEntity), true
	}
	return nil, false
}

func (s *ProfileService) cacheSetRemoteSocialStats(userID string, stats *shedevrumapi.UsersGetSocialStats) {
	s.cache.Set("remote_social_stats_"+userID, stats, cache.DefaultExpiration)
}

func (s *ProfileService) cacheGetRemoteSocialStats(userID string) (*shedevrumapi.UsersGetSocialStats, bool) {
	if stats, ok := s.cache.Get("remote_social_stats_" + userID); ok {
		return stats.(*shedevrumapi.UsersGetSocialStats), true
	}
	return nil, false
}

func (s *ProfileService) cacheSetShedevrumIDByUsertag(usertag string, shedevrumID string) {
	s.cache.Set("usertag_"+usertag, shedevrumID, cache.DefaultExpiration)
}

func (s *ProfileService) cacheGetShedevrumIDByUsertag(usertag string) (string, bool) {
	if shedevrumID, ok := s.cache.Get("usertag_" + usertag); ok {
		return shedevrumID.(string), true
	}
	return "", false
}

func (s *ProfileService) cacheSetTopProfiles(profiles []*domain.ProfileEnity, sort domain.ProfileMetrics_GetTopSort) {
	s.cache.Set("top_profiles_"+string(sort), profiles, cache.DefaultExpiration)
}

func (s *ProfileService) cacheGetTopProfiles(sort domain.ProfileMetrics_GetTopSort) ([]*domain.ProfileEnity, bool) {
	if top, ok := s.cache.Get("top_profiles_" + string(sort)); ok {
		return top.([]*domain.ProfileEnity), true
	}
	return nil, false
}

func (s *ProfileService) shedevrumIDNormalize(shedevrumID string) (string, error) {
	const op = "services.ProfileService.shedevrumIDNormalize"

	if strings.HasPrefix(shedevrumID, "@") {
		if cachedShedevrumID, ok := s.cacheGetShedevrumIDByUsertag(shedevrumID); ok {
			shedevrumID = cachedShedevrumID
		} else {
			remoteFeed, err := s.shedAPI.Users.GetFeed(shedevrumID, 0, "")
			if err != nil {
				s.logger.Error(err.Error(), "op", op)
				return "", err
			}
			s.cacheSetRemoteUser(&remoteFeed.User)
			s.cacheSetShedevrumIDByUsertag(shedevrumID, remoteFeed.User.ID)
			shedevrumID = remoteFeed.User.ID
		}
	}
	return shedevrumID, nil
}

func (s *ProfileService) Create(ctx context.Context, shedevrumID string) error {
	const op = "services.ProfileService.Create"
	if err := s.repoProfile.Create(ctx, shedevrumID); err != nil {
		s.logger.Error(err.Error(), "op", op)
		return nil
	}
	return nil
}

func (s *ProfileService) Get(ctx context.Context, id int) (*domain.ProfileEnity, error) {
	const op = "services.ProfileService.Get"
	profile, err := s.repoProfile.GetByID(ctx, id)
	if err != nil {
		s.logger.Error(err.Error(), "op", op)
		return nil, err
	}
	return profile, nil
}

func (s *ProfileService) GetByShedevrumID(ctx context.Context, shedevrumID string) (*domain.ProfileEnity, error) {
	const op = "services.ProfileService.GetByShedevrumID"

	shedevrumID, err := s.shedevrumIDNormalize(shedevrumID)
	if err != nil {
		s.logger.Error(err.Error(), "op", op)
		return nil, err
	}

	profile, err := s.repoProfile.GetByShedevrumID(ctx, shedevrumID)
	if err != nil {
		s.logger.Error(err.Error(), "op", op)
		return nil, err
	}

	var remoteUser *shedevrumapi.UserEntity
	if chachedRemoteUser, ok := s.cacheGetRemoteUser(shedevrumID); ok {
		remoteUser = chachedRemoteUser
	} else {
		remoteFeed, err := s.shedAPI.Users.GetFeed(shedevrumID, 0, "")
		if err != nil {
			s.logger.Error(err.Error(), "op", op)
			return nil, err
		}
		s.cacheSetRemoteUser(&remoteFeed.User)
		remoteUser = &remoteFeed.User
	}

	var remoteSocialStats *shedevrumapi.UsersGetSocialStats
	if chachedRemoteSocialStats, ok := s.cacheGetRemoteSocialStats(shedevrumID); ok {
		remoteSocialStats = chachedRemoteSocialStats
	} else {
		stats, err := s.shedAPI.Users.GetSocialStats(shedevrumID)
		if err != nil {
			s.logger.Error(err.Error(), "op", op)
			return nil, err
		}
		s.cacheSetRemoteSocialStats(shedevrumID, stats)
		remoteSocialStats = stats
	}

	profile.Name = remoteUser.DisplayName
	profile.AvatarURL = remoteUser.AvatartURL
	profile.Link = remoteUser.ShareLink
	profile.Subscriptions = remoteSocialStats.Subscriptions
	profile.Subscribers = remoteSocialStats.Subscribers
	profile.Likes = remoteSocialStats.Likes

	return profile, nil
}

func (s *ProfileService) GetList(ctx context.Context, startFromID uint64, amount int) ([]*domain.ProfileEnity, error) {
	const op = "services.ProfileService.GetList"
	profiles, err := s.repoProfile.GetList(ctx, startFromID, amount)
	if err != nil {
		s.logger.Error(err.Error(), "op", op)
		return nil, err
	}
	return profiles, nil
}

func (s *ProfileService) GetMetrics(ctx context.Context, shedevrumID string) ([]*domain.ProfileMetricsEntity, error) {
	const op = "services.ProfileService.GetMetrics"
	shedevrumID, err := s.shedevrumIDNormalize(shedevrumID)
	if err != nil {
		s.logger.Error(err.Error(), "op", op)
		return nil, err
	}
	metrics, err := s.repoProfileMetrics.GetByShedevrumID(shedevrumID)
	if err != nil {
		s.logger.Error(err.Error(), "op", op)
		return nil, err
	}
	return metrics, nil
}

func (s *ProfileService) GetTop(ctx context.Context, sort domain.ProfileMetrics_GetTopSort, amount int) ([]*domain.ProfileEnity, error) {
	const op = "services.ProfileService.GetTop"

	var profiles []*domain.ProfileEnity

	if cachedProfiles, ok := s.cacheGetTopProfiles(sort); ok {
		profiles = cachedProfiles
	} else {
		list, err := s.repoProfileMetrics.GetTop(ctx, sort, amount)
		if err != nil {
			s.logger.Error(err.Error(), "op", op)
			return nil, err
		}

		profiles = make([]*domain.ProfileEnity, len(list))
		wg, queue := &sync.WaitGroup{}, make(chan struct{}, 15)

		wg.Add(len(list))
		for i, p := range list {
			go func(i int, p *domain.ProfileMetricsEntity) {
				defer func() {
					<-queue
					wg.Done()
				}()
				queue <- struct{}{}
				profile, err := s.GetByShedevrumID(ctx, p.ShedevrumID)
				if err != nil {
					s.logger.Error(err.Error(), "op", op)
					profile = &domain.ProfileEnity{}
				}
				profiles[i] = profile
			}(i, p)
		}
		wg.Wait()

		s.cacheSetTopProfiles(profiles, sort)
	}

	return profiles, nil
}
