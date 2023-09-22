# Overview
dokkup is a simple container orchestration tool, designed to make managing containers across a cluster of servers as simple as possible.

It's ideal for use cases where you want to orchestrate containers accross one or multiple servers. Think of it as Nomad, Swarm or Kubernetes, but easier to use. 

# Setup
dokkup consists of two parts: the `agent` which runs on your servers, and the CLI tool which runs on your personal machine.

There is no need for a master node, your personal machine is the "master" node.

## Agent
The `agent` runs on your server and will do all of the container management once it receives a request from the CLI tool.

Running the `agent` is as simple as running a simple docker command:
```shell
docker run -d --name dokkup-agent -p 8080:8080 \ 
--restart always -v /path/to/config:/config \ 
-v /var/run/docker.sock:/var/run/docker.sock xiovv/dokkup:latest
```

Next, retreive the API key through the logs:
```shell
docker logs dokkup-agent
```

Output:
```
Your new API key is: jy9DbtDlfi5VJuAkbZYd4Kt0c2cQY8iQ
2023-09-22T10:31:42.175+0200    INFO    agent/main.go:36        server is listening...  {"port": "8080"}
```

Take note of that API key as it will not be printed out the next time the agent runs.

## CLI
The CLI is used to execute jobs and to tell the `agent` what to do.

### Install from source
Clone the repository:
```shell
git clone https://github.com/XiovV/dokkup.git
```

Install the binary:
```shell
cd dokkup/cmd/dokkup
go install
```

Run `dokkup` in your terminal to verify it's installed.

# Deploying jobs
In order to start deploying jobs, we need an inventory and a job file.

## Inventory
The inventory holds our agents, the concept is very similar to Ansible's [inventory](https://docs.ansible.com/ansible/latest/inventory_guide/intro_inventory.html) file. Let's define our `inventory.yaml`:
```yaml
nodes:
 - name: "lab1"
   location: "192.168.0.99:8080"
   key: "Z6wC4goD7V2EiL4XuecTuo8jVxfvwVxs"

 - name: "lab2"
   location: "192.168.0.100:8080"
   key: "EcwxaMO3kyBaKETesxInx7ga3Ti93gvI"

groups:
 - name: "labs"
   nodes: ["lab1", "lab2"]
```

In this example we defined two nodes. We assigned them a name, the IP:PORT where the agent is running and their API key. We also created a 
group called "labs" and assigned the two nodes to it so we can refer to them via the group name, instead of referring to them individually when we start managing our containers.
