package models

type UpdateContainerRequest struct {
	Image     string `json:"image"`
	Container string `json:"container"`
}

type PullImageRequest struct {
	Image string `json:"image"`
}

type NodeStatusRequest struct {
	Container string `json:"container"`
}

type NodeStatusResponse struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	Image  string `json:"image"`
	Status string `json:"status"`
}