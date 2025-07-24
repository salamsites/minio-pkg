package util

type Media struct {
	Sizes   []Size        `json:"sizes"`
	Content []interface{} `json:"content"`
}

type Size struct {
	Width   int `json:"width"`
	Height  int `json:"height"`
	Quality int `json:"quality"`
}

type Err struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type FeedResultTypeImage struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Mime string `json:"mime"`
}

type FeedResultTypeVideo struct {
	Path     string `json:"path"`
	Type     string `json:"type"`
	Duration int64  `json:"duration"`
	Mime     string `json:"mime"`
}

type FeedResultTypeAudio struct {
	Path     string `json:"path"`
	Type     string `json:"type"`
	Duration int64  `json:"duration"`
	Mime     string `json:"mime"`
}

type FeedResultTypeFile struct {
	Path     string `json:"path"`
	Type     string `json:"type"`
	FileSize int64  `json:"file_size"`
	Mime     string `json:"mime"`
}
