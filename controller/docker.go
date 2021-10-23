package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type DockerController struct {
	node   string
	apiKey string
}

func NewDockerController(node, apiKey string) DockerController {
	return DockerController{node: node, apiKey: apiKey}
}

func (dc DockerController) handleAgentResponse(r *http.Response) error {
	if r.StatusCode != 200 {
		var errResponse ErrorResponse
		if err := json.NewDecoder(r.Body).Decode(&errResponse); err != nil {
			return fmt.Errorf("received unexpected response from the agent")
		}
		return fmt.Errorf(errResponse.Error)
	}

	return nil
}

func (dc DockerController) GetContainerImage(containerName string) (string, error) {
	URL := fmt.Sprintf("%s%s/%s", dc.node, GetContainerImageURL, containerName)

	r, err := dc.getRequest(URL)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()

	err = dc.handleAgentResponse(r)

	if err != nil {
		return "", err
	}

	var body struct {
		Image string `json:"image"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return "", err
	}

	return body.Image, nil
}

func (dc DockerController) RollbackContainer(containerName string) error {
	URL := fmt.Sprintf("%s%s?container=%s", dc.node, RollbackContainerURL, containerName)

	r, err := dc.putRequest(URL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return dc.handleAgentResponse(r)
}

func (dc DockerController) PullImage(image string) error {
	URL := fmt.Sprintf("%s%s?image=%s", dc.node, PullImageURL, image)

	r, err := dc.putRequest(URL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return dc.handleAgentResponse(r)
}

func (dc DockerController) UpdateContainer(containerName, image string, keepContainer bool) error {
	URL := fmt.Sprintf("%s%s?container=%s&image=%s&keep=%s", dc.node, UpdateContainerURL, containerName, image, strconv.FormatBool(keepContainer))

	r, err := dc.putRequest(URL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return dc.handleAgentResponse(r)
}

func (dc DockerController) getRequest(location string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, location, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("key", dc.apiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (dc DockerController) putRequest(location string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPut, location, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("key", dc.apiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
