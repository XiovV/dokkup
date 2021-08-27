package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type App struct {
	config *Config
}

func New(config *Config) *App {
	return &App{config: config}
}

func (a *App) HandleStatus(node, group string) {
	if node == "" && group == "" {
		fmt.Println("either node od group need to be set.")
		return
	}

	if node != "" && group != "" {
		fmt.Println("only node or group should be set, not both.")
		return
	}

	if node != "" {
		fmt.Printf("%s status\n", node)
		a.GetNodeStatus(node)
	}
}

func (a *App) GetNodeStatus(node string) {
	foundNode, ok := a.config.FindNodeByName(node)
	if !ok {
		fmt.Println("couldn't find node")
		os.Exit(1)
	}

	containers := a.findContainersInGroup(foundNode.Group)
	var request NodeStatusRequest
	for _, container := range containers {
		request.Container = container

		response, err := sendNodeStatusRequest(foundNode.Location, request)
		if err != nil {
			if errors.Is(err, ErrContainerNotFound) {
				fmt.Printf("container %s doesn't exist\n", container)
			}
		} else {
			fmt.Printf("name: %s\n", response.Name)
			fmt.Printf("id: %s\n", response.ID)
			fmt.Printf("image: %s\n", response.Image)
			fmt.Printf("status: %s\n", response.Status)
		}
	}

}

func (a *App) findContainersInGroup(groupName string) []string {
	group, _ := a.config.FindGroupByName(groupName)

	var containers []string
	for _, container := range group.Containers {
		containers = append(containers, container)
	}

	return containers
}

func (a *App) HandleUpdate(groupName string) {
	group, ok := a.config.FindGroupByName(groupName)
	if !ok {
		fmt.Printf("couldn't find a group with the name: %s\n", groupName)
		os.Exit(1)
	}

	err := a.pullImages(group.Nodes, group.Image)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = a.updateContainersInGroup(group)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (a *App) pingNodes(nodes []Node) {
	for _, node := range nodes {
		resCode, err := sendGetRequest(node.Location+"/api/health")
		if err != nil {
			fmt.Println("couldn't ping:", err)
		}

		if resCode != 200 {
			fmt.Printf("%s at %s is not online!\n", node.NodeName, node.Location)
		} else {
			fmt.Printf("%s is online\n", node.NodeName)
		}
	}

}

func sendGetRequest(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func (a *App) updateContainersInGroup(group Groups) error {
	var reqBody UpdateRequest
	reqBody.Image = group.Image

	for _, node := range group.Nodes {
		for	_, container := range group.Containers {
			reqBody.Container = container
			reqBodyMarshalled, err := json.Marshal(reqBody)
			if err != nil {
				return err
			}

			resCode, err := sendPostRequest(node.Location + "/api/containers/update", reqBodyMarshalled)
			if err != nil {
				fmt.Printf("error while updating %s on %s: %s\n", container, node.NodeName, err)
			}

			switch resCode {
			case http.StatusOK:
				fmt.Printf("successfully updated %s on %s\n", container, node.NodeName)
			case http.StatusInternalServerError:
				fmt.Printf("couldn't update %s on %s. Error: %d\n", container, node.NodeName, resCode)
			}
		}
	}

	return nil
}

func (a *App) pullImages(nodes []Node, image string) error {
	reqBody := PullImageRequest{Image: image}
	reqBodyMarshalled, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	
	for _, node := range nodes {
		resCode, err := sendPostRequest(node.Location + "/api/images/pull", reqBodyMarshalled)
		if err != nil {
			fmt.Printf("error while pulling %s on %s: %s\n", image, node.NodeName, err)
			continue
		}

		switch resCode {
		case http.StatusOK:
			fmt.Printf("successfully pulled %s on %s\n", image, node.NodeName)
		case http.StatusInternalServerError:
			fmt.Printf("couldn't pull %s on %s. Error: %d\n", image, node.NodeName, resCode)
		}
	}
	
	return nil
}

func sendPostRequest(url string, body []byte) (int, error){
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
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

func sendNodeStatusRequest(node string, body NodeStatusRequest) (NodeStatusResponse, error) {
	marshalBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", node+"/api/nodes/status", bytes.NewBuffer(marshalBody))
	if err != nil {
		return NodeStatusResponse{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return NodeStatusResponse{}, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 404:
		return NodeStatusResponse{}, ErrContainerNotFound
	}

	var response NodeStatusResponse
	responseBody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return NodeStatusResponse{}, err
	}

	return response, nil
}