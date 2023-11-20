package shedevrumapi

import (
	"net/http"
)

type Config struct {
	OAuth     string
	UserAgent string
}

type ShedevrumAPI struct {
	Feed   *FeedAPI
	User   *UserAPI
	Users  *UsersAPI
	Images *ImagesAPI
}

func NewShedevrumAPI(config Config) *ShedevrumAPI {
	// set default headers
	header := make(http.Header)
	if config.OAuth != "" {
		header.Add("Authorization", "OAuth "+config.OAuth)
	}
	header.Add("User-Agent", config.UserAgent)

	// init client
	client := &http.Client{}

	return &ShedevrumAPI{
		Feed:   NewFeedAPI(baseURL, client, header),
		User:   NewUserAPI(baseURL, client, header),
		Users:  NewUsersAPI(baseURL, client, header),
		Images: NewImagesAPI(baseURL, client, header),
	}
}
