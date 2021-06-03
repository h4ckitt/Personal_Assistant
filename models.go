package main

type data struct {
	User             user             `json:"user"`
	ExtendedEntities extendedentities `json:"extended_entities"`
	RetweetCount     int              `json:"retweet_count"`
	FavoriteCount    int              `json:"favorite_count"`
}

type user struct {
	Name     string `json:"name"`
	Username string `json:"screen_name"`
}

type extendedentities struct {
	Media []media `json:"media"`
}

type media struct {
	Type      string    `json:"type"`
	VideoInfo videoinfo `json:"video_info"`
}

type videoinfo struct {
	Variants []variants `json:"variants"`
}

type variants struct {
	Bitrate     int    `json:"bitrate,omitempty"`
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
}

type replyData struct {
	Data      data
	Qualities []string
}
