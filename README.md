## What is Docker Control (temporary name)?
Docker Control is a simple CLI tool with an ability to mass-update docker containers.
If you've got hundreds of servers running containers, you can update the containers running on
all of them with this tool using a single command. Simply create your own config.json, and you're good to go.

## Usage
In order to use this tool, you first have to define your own config.json.

### config.json example
```json
{
  "groups": [
    {
      "group": "rest-api",
      "containers": ["api-instance-1", "api-instance-2", "api-instance-3"],
      "image": "example/example-image:1.0.0",
      "nodes": [
        {
          "location": "https://node1.example.com:5006",
          "node_name": "api-node1"
        },
        {
          "location": "https://node2.example.com:5006",
          "node_name": "api-node2"
        }
      ]
    },
    {
      "group": "database-servers",
      "containers": ["postgres"],
      "image": "postgres:latest",
      "nodes": [
        {
          "location": "https://node1.example.com:5006",
          "node_name": "database-node1"
        },
        {
          "location": "https://node2.example.com:5006",
          "node_name": "database-node2"
        }
      ]
    }
  ]
}
```

#### config.json explanation:
- `group` is a server group
- `containers` an array defining which containers you'd like to update
- `image` an image you'd like the containers to update to 
- `nodes` an array where you define all your nodes (a server that's running docker-control-agent)
- `location` where the docker-control-agent is running
- `node_name` custom name for that specific node

## Command Line Usage
```
$ docker_control update --group rest-api
```
This command will attempt to update the containers with the names 
"api-instance-1", "api-instance-2" and "api-instance-3" which are running on
nodes "api-node1" and "api-node2" to the image "example/example-image:1.0.0".