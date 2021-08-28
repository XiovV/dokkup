package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/XiovV/docker_control_cli/models"
	"io/ioutil"
	"net/http"
)

type DockerService interface {
	GetContainerStatus(string, string) (models.NodeStatusResponse, error)
	UpdateContainer(string, models.UpdateContainerRequest) (int, error)
	PullImage(string, models.PullImageRequest) (int, error)
}

type DockerController struct {}

func NewDockerController() DockerController {
	return DockerController{}
}

func (dc DockerController) PullImage(location string, request models.PullImageRequest) (int, error) {
	body, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", location + "/api/images/pull", bytes.NewBuffer(body))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func (dc DockerController) UpdateContainer(location string, request models.UpdateContainerRequest) (int, error) {
	body, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", location + "/api/containers/update", bytes.NewBuffer(body))
	if err != nil {
		return http.StatusInternalServerError, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func (dc DockerController) GetContainerStatus(url, containerName string) (models.NodeStatusResponse, error) {
	fmt.Println("getting container status")
	body := models.NodeStatusRequest{Container: containerName}
	marshalBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url+"/api/nodes/status", bytes.NewBuffer(marshalBody))
	if err != nil {
		return models.NodeStatusResponse{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.NodeStatusResponse{}, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 404:
		return models.NodeStatusResponse{}, models.ErrContainerNotFound
	}

	var response models.NodeStatusResponse
	responseBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return models.NodeStatusResponse{}, err
	}

	return response, nil
}