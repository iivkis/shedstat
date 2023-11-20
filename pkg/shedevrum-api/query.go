package shedevrumapi

import "net/url"

func AddQuery(query url.Values) url.Values {
	query.Set("content", "combined")
	query.Set("from", "web")
	query.Set("appPlatform", "web")
	return query
}
