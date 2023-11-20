package shedevrumapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/goccy/go-json"
)

type UsersAPI struct {
	baseURL string
	client  *http.Client
	header  http.Header
}

func NewUsersAPI(baseURL string, client *http.Client, header http.Header) *UsersAPI {
	return &UsersAPI{
		baseURL: baseURL + "users/",
		client:  client,
		header:  header.Clone(),
	}
}

type UsersGetFeed struct {
	UserID string        `json:"userID"`
	User   UserEntity    `json:"user"`
	Posts  []*PostEntity `json:"posts"`
	Next   string        `json:"next"`
}

func (r *UsersAPI) GetFeed(userID string, amount uint64, startFrom string) (*UsersGetFeed, error) {
	req, err := http.NewRequest("GET", r.baseURL+userID+"/feed", nil)
	if err != nil {
		return nil, err
	}
	req.Header = r.header
	q := req.URL.Query()
	q = AddQuery(q)
	q.Add("amount", strconv.FormatUint(amount, 10))
	if startFrom != "" {
		q.Add("startFrom", startFrom)
	}
	req.URL.RawQuery = q.Encode()
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrBadStatusCode, resp.StatusCode)
	}
	var body UsersGetFeed
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return &body, nil
}

type UsersGetSubscriptions struct {
	Next          string        `json:"next"`
	Subscriptions []*UserEntity `json:"subscriptions"`
}

func (r *UsersAPI) GetSubscriptions(userID string, amount uint64, startFrom string) (*UsersGetSubscriptions, error) {
	req, err := http.NewRequest("GET", r.baseURL+userID+"/subscriptions", nil)
	if err != nil {
		return nil, err
	}
	req.Header = r.header
	q := req.URL.Query()
	q = AddQuery(q)
	q.Add("amount", strconv.FormatUint(amount, 10))
	if startFrom != "" {
		q.Add("startFrom", startFrom)
	}
	req.URL.RawQuery = q.Encode()
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrBadStatusCode, resp.StatusCode)
	}
	var body UsersGetSubscriptions
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return &body, nil
}

type UsersGetSubscribers struct {
	Next        string        `json:"next"`
	Subscribers []*UserEntity `json:"subscribers"`
}

func (r *UsersAPI) GetSubscribers(userID string, amount uint64, startFrom string) (*UsersGetSubscribers, error) {
	req, err := http.NewRequest("GET", r.baseURL+userID+"/subscribers", nil)
	if err != nil {
		return nil, err
	}
	req.Header = r.header
	q := req.URL.Query()
	q = AddQuery(q)
	q.Add("amount", strconv.FormatUint(amount, 10))
	if startFrom != "" {
		q.Add("startFrom", startFrom)
	}
	req.URL.RawQuery = q.Encode()
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrBadStatusCode, resp.StatusCode)
	}
	var body UsersGetSubscribers
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return &body, nil
}

func (r *UsersAPI) Subscribe(userID string) error {
	req, err := http.NewRequest("PUT", r.baseURL+userID+"/subscribe", nil)
	if err != nil {
		return err
	}
	req.Header = r.header
	q := req.URL.Query()
	q = AddQuery(q)
	req.URL.RawQuery = q.Encode()
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %d", ErrBadStatusCode, resp.StatusCode)
	}
	return nil
}

func (r *UsersAPI) Unsubscribe(userID string) error {
	req, err := http.NewRequest("PUT", r.baseURL+userID+"/unsubscribe", nil)
	if err != nil {
		return err
	}
	req.Header = r.header
	q := req.URL.Query()
	q = AddQuery(q)
	req.URL.RawQuery = q.Encode()
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %d", ErrBadStatusCode, resp.StatusCode)
	}
	return nil
}

type UsersGetSocialStats struct {
	Likes         uint64 `json:"likes"`
	Subscribers   uint64 `json:"subscribers"`
	Subscriptions uint64 `json:"subscriptions"`
}

func (r *UsersAPI) GetSocialStats(userID string) (*UsersGetSocialStats, error) {
	req, err := http.NewRequest("GET", r.baseURL+userID+"/social_stats", nil)
	if err != nil {
		return nil, err
	}
	req.Header = r.header
	q := req.URL.Query()
	q = AddQuery(q)
	req.URL.RawQuery = q.Encode()
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrBadStatusCode, resp.StatusCode)
	}
	var body UsersGetSocialStats
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return &body, nil
}
