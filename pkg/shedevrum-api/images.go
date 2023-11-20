package shedevrumapi

import (
	"fmt"
	"net/http"
)

type ImagesAPI struct {
	baseURL string
	client  *http.Client
	header  http.Header
}

func NewImagesAPI(baseURL string, client *http.Client, header http.Header) *ImagesAPI {
	return &ImagesAPI{
		baseURL: baseURL + "images/",
		client:  client,
		header:  header.Clone(),
	}
}

func (r *ImagesAPI) Like(id string) error {
	req, err := http.NewRequest("PUT", r.baseURL+id+"/like", nil)
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

func (r *ImagesAPI) Unlike(id string) error {
	req, err := http.NewRequest("PUT", r.baseURL+id+"/unlike", nil)
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
