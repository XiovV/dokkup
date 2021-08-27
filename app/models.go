package app

type UpdateRequest struct {
	Image     string `json:"image"`
	Container string `json:"container"`
}

type PullImageRequest struct {
	Image string `json:"image"`
}

