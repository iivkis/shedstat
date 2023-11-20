package shedevrumapi

import (
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
)

type UserAPI struct {
	baseURL string
	client  *http.Client
	header  http.Header
}

func NewUserAPI(baseURL string, client *http.Client, header http.Header) *UserAPI {
	return &UserAPI{
		baseURL: baseURL + "user/",
		client:  client,
		header:  header.Clone(),
	}
}

type UserGet struct {
	UserEntity
}

func (r *UserAPI) GetMe() (*UserGet, error) {
	req, err := http.NewRequest("GET", r.baseURL, nil)
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
	var body UserGet
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return &body, nil
}
