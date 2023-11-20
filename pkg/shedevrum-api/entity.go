package shedevrumapi

type UserEntity struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	AvatartURL  string `json:"avatarURL"`
	AvatarID    string `json:"avatarIdentifier"`
	ShareLink   string `json:"shareLink"`
	Subscribed  bool   `json:"subscribed"`
	BlockedByMe bool   `json:"blockedByMe"`
	BlockedMe   bool   `json:"blockedMe"`
	Verified    bool   `json:"verified"`
}

type PostParamsEntity struct {
	Prompt string `json:"prompt"`
	Seed   int    `json:"seed"`
	IsExp  bool   `json:"is_experiment"`
}

type PostEntity struct {
	ID               string           `json:"id"`
	URL              string           `json:"url"`
	TranslatedPrompt string           `json:"translatedPrompt"`
	Kind             string           `json:"kind"`
	Likes            uint             `json:"likes"`
	User             UserEntity       `json:"user"`
	Liked            bool             `json:"liked"`
	Params           PostParamsEntity `json:"params"`
	PostURL          string           `json:"postURL"`
	CreatedAt        int64            `json:"createdAt"`
	Status           string           `json:"status"`
	ImageStatus      string           `json:"imageStatus"`
	CommentsBranchID string           `json:"commentsBranchID"`
}
