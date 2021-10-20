package app

import (
	"github.com/XiovV/docker_control_cli/services"
)

type App struct {
	config *Config
	inventory *Inventory
	dockerService services.DockerService
}

func New(config *Config, inventory *Inventory, dockerService services.DockerService) *App {
	return &App{config: config, inventory: inventory, dockerService: dockerService}
}

//func (a *App) HandleStatus(node, group string) {
//	if node == "" && group == "" {
//		fmt.Println("either node od group need to be set.")
//		return
//	}
//
//	if node != "" && group != "" {
//		fmt.Println("only node or group should be set, not both.")
//		return
//	}
//
//	if node != "" {
//		fmt.Printf("%s status\n", node)
//		a.GetNodeStatus(node)
//	}
//}

//func (a *App) GetNodeStatus(node string) {
//	foundNode, ok := a.config.FindNodeByName(node)
//	if !ok {
//		fmt.Println("couldn't find node")
//		os.Exit(1)
//	}
//
//	containers := a.config.FindContainersInGroup(foundNode.Group)
//	for _, container := range containers {
//		containerStatus, err := a.dockerService.GetContainerStatus(foundNode.Location, container)
//
//		if err != nil {
//			if errors.Is(err, models.ErrContainerNotFound) {
//				fmt.Printf("container %s doesn't exist\n", container)
//			}
//		} else {
//			fmt.Printf("name: %s\n", containerStatus.Name)
//			fmt.Printf("id: %s\n", containerStatus.ID)
//			fmt.Printf("image: %s\n", containerStatus.Image)
//			fmt.Printf("status: %s\n", containerStatus.Status)
//		}
//	}
//
//}

//func (a *App) HandleUpdate(groupName string) {
//	group, ok := a.config.FindGroupByName(groupName)
//	if !ok {
//		fmt.Printf("couldn't find a group with the name: %s\n", groupName)
//		os.Exit(1)
//	}
//
//	start := time.Now()
//	err := a.updateContainersInGroup(group)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println("duration:", time.Since(start))
//}
//
//func (a *App) updateContainersInGroup(group Groups) error {
//	var reqBody models.UpdateContainerRequest
//	reqBody.Image = group.Image
//
//	for _, node := range group.Nodes {
//		for	_, container := range group.Containers {
//			fmt.Printf("attempting to update %s on %s\n", container, node.NodeName)
//
//			reqBody.Container = container
//
//			resCode, err := a.dockerService.UpdateContainer(node.Location, reqBody)
//
//			if err != nil {
//				fmt.Printf("error while updating %s on %s: %s\n", container, node.NodeName, err)
//			}
//
//			switch resCode {
//			case http.StatusOK:
//				fmt.Printf("successfully updated %s on %s\n", container, node.NodeName)
//			case http.StatusNotFound:
//				fmt.Printf("%s doesn't exist on %s\n", container, node.NodeName)
//			}
//		}
//	}
//
//	return nil
//}
//
//func (a *App) pullImages(nodes []Group, image string) error {
//	reqBody := models.PullImageRequest{Image: image}
//
//	for _, node := range nodes {
//		fmt.Printf("attempting to pull %s on %s\n", image, node.NodeName)
//		resCode, err := a.dockerService.PullImage(node.Location, reqBody)
//		if err != nil {
//			fmt.Printf("error while pulling %s on %s: %s\n", image, node.NodeName, err)
//			continue
//		}
//
//		switch resCode {
//		case http.StatusOK:
//			fmt.Printf("successfully pulled %s on %s\n", image, node.NodeName)
//		case http.StatusInternalServerError:
//			fmt.Printf("couldn't pull %s on %s. Error: %d\n", image, node.NodeName, resCode)
//		}
//	}
//
//	return nil
//}