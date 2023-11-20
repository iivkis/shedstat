package shedevrumapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/goccy/go-json"
)

type FeedAPI struct {
	baseURL string
	client  *http.Client
	header  http.Header
}

func NewFeedAPI(baseURL string, client *http.Client, header http.Header) *FeedAPI {
	return &FeedAPI{
		baseURL: baseURL + "feed/",
		client:  client,
		header:  header.Clone(),
	}
}

type FeedTopPeriod string

const (
	FEED_TOP_PERIOD_DAY  FeedTopPeriod = "day"
	FEED_TOP_PERIOD_WEEK FeedTopPeriod = "week"
	FEED_TOP_PERIOD_ALL  FeedTopPeriod = "all"
)

type FeedGetTopResponse struct {
	Next  string        `json:"next"`
	Posts []*PostEntity `json:"posts"`
}

func (r *FeedAPI) GetTop(period FeedTopPeriod, amount uint64, startFrom string) (*FeedGetTopResponse, error) {
	req, err := http.NewRequest("GET", r.baseURL+"top", nil)
	if err != nil {
		return nil, err
	}
	req.Header = r.header
	q := req.URL.Query()
	q = AddQuery(q)
	q.Add("period", string(period))
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
	var body FeedGetTopResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return &body, nil
}

type FeedGetSubscriptionsResponse struct {
	Next  string        `json:"next"`
	Posts []*PostEntity `json:"posts"`
}

func (r *FeedAPI) GetSubscriptions(amount uint64, startFrom string) (*FeedGetSubscriptionsResponse, error) {
	req, err := http.NewRequest("GET", r.baseURL+"subscriptions", nil)
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
	var body FeedGetSubscriptionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return &body, nil
}

type FeedGetRecentResponse struct {
	Next  string        `json:"next"`
	Posts []*PostEntity `json:"posts"`
}

func (r *FeedAPI) GetRecent(amount uint64, startFrom string) (*FeedGetRecentResponse, error) {
	req, err := http.NewRequest("GET", r.baseURL+"recent", nil)
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
	var body FeedGetRecentResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return &body, nil
}
