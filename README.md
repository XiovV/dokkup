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

- **nodes** is an array of nodes (or servers, both terms are used interchangeably) where the `agents` are running.
- **groups** is an array of groups, where you can group together multiple nodes to avoid referring to them individually.

## Job
The job specification holds information which tells the `agent` how to deploy a job. Let's define our `demo.yaml`:
```yaml
group: "labs"
count: 2
name: "demo"

container:
 - image: "crccheck/hello-world"
   ports:
    - in: 8000
   restart: always 
   labels:
     - "my.label.test=demo"
   environment:
     - MYENV=ENVEXAMPLE
   volumes:
     - myvolume:/home
   # networks:
   # - mynetwork
```

- **group** specifies which group we want this job to be deployed on. Earlier on we created a "labs" group, so the job will be deployed on our lab1 and lab2 nodes.
- **count** specifies how many containers we want to deploy on each node.
- **name** gives the job a name.
- **container** holds information about the container you want to deploy. It closely resembles a standard [docker-compose](https://docs.docker.com/compose) file.

## Deploying the job
Now that we have our inventory and job specification, we can deploy our job:

```shell
$ dokkup run job -i inventory.yaml demo.yaml
```
```
Deployment summary:

NAME     IMAGE                    RESTART     COUNT     GROUP     NETWORK
demo     crccheck/hello-world     always      2         labs      bridge

Node statuses:

NAME      STATUS     CONTAINERS     UPDATE
lab1*     ONLINE     0/0            true
lab2*     ONLINE     0/0            true

Are you sure you want to proceed? (y/n) 
```
Tip #1: If there's an `inventory.yaml` in your current directory, you can omit the -i flag, dokkup loads `inventory.yaml` files by default. 

Tip #2: You can provide a -y or --yes flag to skip the confirmation prompt.

The CLI will show a deployment summary, showing some basic information about the job and the container you are about to deploy. And it will display the nodes on which the job will be deployed. The asterisk next to the node name signifies that a job will be deployed from scratch.

TODO: Show gif of the deployment process here.
