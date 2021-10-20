package controller

import (
	"fmt"
	"net/http"
)

type DockerController struct {
	node string
	apiKey string
}

func NewDockerController(node, apiKey string) DockerController {
	return DockerController{node: node, apiKey: apiKey}
}

func (dc DockerController) PullImage(image string) error {
	URL := fmt.Sprintf("%s%s?image=%s", dc.node, PullImageURL, image)

	statusCode, err := dc.putRequest(URL)
	if err != nil {
		return err
	}

	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusForbidden:
		return fmt.Errorf("invalid api key")
	case http.StatusInternalServerError:
		return fmt.Errorf("image couldn't be pulled")
	default:
		return fmt.Errorf("an unknown error has occured while pulling the image")
	}
}

func (dc DockerController) UpdateContainer(containerName, image string) error {
	URL := fmt.Sprintf("%s%s?container=%s&image=%s", dc.node, UpdateContainerURL, containerName, image)

	statusCode, err := dc.putRequest(URL)
	if err != nil {
		return err
	}

	switch statusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("container %s doesn't exist", containerName)
	case http.StatusForbidden:
		return fmt.Errorf("invalid api key")
	case http.StatusInternalServerError:
		return fmt.Errorf("container couldn't be updated")
	default:
		return fmt.Errorf("an unknown error has occured while pulling the image")
	}
}

func (dc DockerController) putRequest(location string) (int, error){
	request, err := http.NewRequest(http.MethodPut, location, nil)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	request.Header.Set("key", dc.apiKey)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	defer response.Body.Close()

	return response.StatusCode, nil
}