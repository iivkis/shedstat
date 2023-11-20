package shedevrumapi

import (
	"net/url"
	"strings"
)

// https://shedevrum.ai/profile/PrinceBelka/
// https://shedevrum.ai/@prince/
func ParseProfileIDFromURL(profileURL string) (string, error) {
	u, err := url.Parse(profileURL)
	if err != nil {
		return "", err
	}

	path := strings.Split(u.Path, "/")
	if path[1] == "profile" {
		return path[2], nil
	}

	return path[1], nil
}

func CreateProfileURL(profileID string) string {
	return "https://shedevrum.ai/profile/" + profileID
}
