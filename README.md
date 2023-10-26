# Overview
dokkup is a container orchestrator designed to make managing containers across a cluster of servers as simple as possible.

It's ideal for use cases where you want to orchestrate containers accross one or multiple servers. Think of it as Nomad, Swarm or Kubernetes, but easier to use. 

# Features
- 🛠️ **Easy to use** - Build production-ready infrastructure from bare servers in a matter of minutes.
- 📈 **Zero downtime** - Dokkup ensures zero downtime during container updates, even if the update fails. 
- 💻 **No master node** - No need to set up a server just for a master node, your personal machine is all you need.
- 🚀 **Sky is the limit** - Got only one server? No problem. Got 10,000 servers? No problem.

# Table of contents
- [Documentation](#documentation)
- [Setup](#setup)
  - [Agent setup](#agent)
  - [CLI setup](#cli)
- [Deploying jobs](#deploying-jobs)
  - [Inventory file](#inventory)
  - [Job specificaton](#job)
  - [Deploy a job](#deploy-a-job)
- [Managing jobs](#managing-jobs)
  - [Update](#update-a-job)
  - [Canary update](#canary-update)
  - [Scaling](#upscale-or-downscale-a-job)
  - [Rollback](#rollback-a-job)
  - [Stop](#stop-a-job)
  - [Remove/purge](#removepurge-a-job)
- [Show job status](#show-job-status)

# Documentation
More detailed documentation can be found [here](https://xiovv.github.io/dokkup). However, this README will suffice if you're just getting started.

# Setup
dokkup consists of two parts: the agent and the CLI tool:

- **agent** runs on the servers/nodes where you want to run your containers. It listens for incoming job requests and executes them.
- **CLI** is used to dispatch job requests to the servers/nodes.

The CLI requires two files, the inventory and the job files:

- **inventory** contains information about your nodes, such as their agent's API key, URL and name.
- **job** contains the configuration for the containers you want to deploy, it closely resembles a standard [docker-compose](https://docs.docker.com/compose) file.

There is no need for a master node, your personal machine is the "master" node.

## Agent
The `agent` runs on your server and will do all of the container management once it receives a request from the CLI tool.

Running the `agent` is as simple as running a simple docker command:
```shell
$ docker run -d --name dokkup-agent -p 8080:8080 \ 
--restart always -v /path/to/config:/config \ 
-v /var/run/docker.sock:/var/run/docker.sock xiovv/dokkup:latest
```

Next, retreive the API key through the logs:
```shell
$ docker logs dokkup-agent
```
```
Your new API key is: jy9DbtDlfi5VJuAkbZYd4Kt0c2cQY8iQ
2023-09-22T10:31:42.175+0200    INFO    agent/main.go:36        server is listening...  {"port": "8080"}
```

Take note of that API key as it will not be printed out the next time the agent runs.

## CLI
The CLI is used to execute jobs and to tell the `agent` what to do.

### Install from Releases
Go to [Releases](https://github.com/XiovV/dokkup/releases) and download the binary for your OS and CPU architecture.

### Install from source
Clone the repository:
```shell
$ git clone https://github.com/XiovV/dokkup.git
```

Install the binary:
```shell
cd dokkup/cmd/dokkup
go install
```

Run `dokkup` in your terminal to verify it's installed:
```shell
$ dokkup version
```
```
Dokkup v0.1.0-beta, build fd13a57
```

# Deploying jobs
Before we can start deploying jobs, we need an inventory and a job file.

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

- **group** specifies which group we want this job to be deployed on. Earlier on we created a "labs" group, so that's what we're using here.
- **count** specifies how many containers we want to deploy on each node.
- **name** gives the job a name.
- **container** holds information about the container you want to deploy. It closely resembles a standard [docker-compose](https://docs.docker.com/compose) file.

## Deploy a job
Now that we have our inventory and job specification, we can deploy our job:

```shell
$ dokkup run job -i inventory.yaml demo.yaml
```
```
Deployment summary:

NAME     IMAGE                    RESTART     COUNT     GROUP     NETWORK
demo     crccheck/hello-world     always      2         labs      bridge

Node statuses:

NAME      STATUS     CONTAINERS     UPDATE     VERSION
lab1*     ONLINE     0 -> 2         true       55dab35
lab2*     ONLINE     0 -> 2         true       55dab35

Are you sure you want to proceed? (y/n) 
```

- Tip 1: If there's an `inventory.yaml` in your current directory, you can omit the -i flag, dokkup loads `inventory.yaml` files by default. 
- Tip 2: You can provide a -y or --yes flag to skip the confirmation prompt.

The CLI will show a deployment summary, showing some basic information about the job, such as how many containers it's going to run and the hashed version tag of the job, along with the nodes on which the job will be deployed. The asterisk next to the node name signifies that a job will be deployed from scratch.

Now we can run `docker ps` on our nodes and see our containers (this is for demonstration purposes, you don't have to do this):
```shell
$ docker ps
```
lab1 output:
```
CONTAINER ID   IMAGE                  COMMAND                   CREATED          STATUS                             PORTS                                       NAMES
adec3b3612ca   crccheck/hello-world   "/bin/sh -c 'echo \"h…"   14 seconds ago   Up 12 seconds (health: starting)   0.0.0.0:32805->8000/tcp                     demo-b905568e-942a-4ef4-b091-45f9fc2ddea9
e3a4b99761b1   crccheck/hello-world   "/bin/sh -c 'echo \"h…"   14 seconds ago   Up 12 seconds (health: starting)   0.0.0.0:32804->8000/tcp                     demo-78d37fdd-cbfc-468c-88b3-b8a1df17855b
23e0af35bdae   xiovv/dokkup:latest    "/agent"                  25 hours ago     Up 25 hours                        0.0.0.0:8080->8080/tcp, :::8080->8080/tcp   dokkup-agent
```

lab2 output:
```
CONTAINER ID   IMAGE                  COMMAND                   CREATED          STATUS                    PORTS                                       NAMES
cf00cd390db3   crccheck/hello-world   "/bin/sh -c 'echo \"h…"   36 seconds ago   Up 34 seconds (healthy)   0.0.0.0:32793->8000/tcp                     demo-d71fc907-6b29-4287-aab8-d29a5fe4e821
0e8e013adea2   crccheck/hello-world   "/bin/sh -c 'echo \"h…"   37 seconds ago   Up 35 seconds (healthy)   0.0.0.0:32792->8000/tcp                     demo-4d2c8f50-48e7-4817-b380-4b790196d34f
2386f86f788b   xiovv/dokkup:latest    "/agent"                  25 hours ago     Up 25 hours               0.0.0.0:8080->8080/tcp, :::8080->8080/tcp   dokkup-agent
```

In case you run the `dokkup run job` command without making any changes, you don't have to worry about dokkup wiping your existing containers and re-deploying them again, it will detect that nothing has changed and it won't do anything:
```
Deployment summary:

NAME     IMAGE                    RESTART     COUNT     GROUP     NETWORK
demo     crccheck/hello-world     always      2         labs      bridge

Node statuses:

NAME     STATUS     CONTAINERS     UPDATE     VERSION
lab1     ONLINE     2/2            false      55dab35
lab2     ONLINE     2/2            false      55dab35

Are you sure you want to proceed? (y/n) 

```
The CLI will show how many containers are running and the update status which will signify whether the job is going to be updated or not.

# Managing jobs

## Update a job
Updating a job is as simple as making a change in the job specification file and running the `dokkup run job` command again: \
`demo.yaml`
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
   # - "my.label.test=demo"
     - "my.label.test=somechange"
   environment:
     - MYENV=ENVEXAMPLE
   volumes:
     - myvolume:/home
   # networks:
   # - mynetwork

```
```shell
$ dokkup run job demo.yaml
```
```
Deployment summary:

NAME     IMAGE                    RESTART     COUNT     GROUP     NETWORK
demo     crccheck/hello-world     always      2         labs      bridge

Node statuses:

NAME     STATUS     CONTAINERS     UPDATE     VERSION
lab1     ONLINE     2/2            true       55dab35 -> 9470cdc
lab2     ONLINE     2/2            true       55dab35 -> 9470cdc

Are you sure you want to proceed? (y/n) 
```
Note: we omitted the -i inventory.yaml flag because dokkup automatically loads files called `inventory.yaml`.

The update status is now true, meaning that dokkup is going to take down the currently running containers and deploy new ones. If something goes wrong during the update, dokkup will abort and run the old containers, ensuring minimum downtime in case something goes wrong. It also shows the new version hash.

## Canary update
This is work in progress.

## Upscale or downscale a job
In case you want to change the number of running containers, you can do so by changing the `count` field in your job specification, then run the `dokkup run job` command: 
`demo.yaml`
```yaml
group: "labs"
# count: 2
count: 5
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
```shell
$ dokkup run job demo.yaml
```
```
Deployment summary:

NAME     IMAGE                    RESTART     COUNT     GROUP     NETWORK
test     crccheck/hello-world     always      5         local     bridge

Node statuses:

NAME          STATUS     CONTAINERS     UPDATE     VERSION
lab1          ONLINE     2 -> 5         true       55dab32
lab2          ONLINE     2 -> 5         true       55dab32

Are you sure you want to proceed? (y/n) 
```
This will upscale the job, meaning that the agent will keep the already running containers and start up 3 more containers. Conversely, if you were to reduce the `count` value, the agent will remove the extra containers, effectively downscaling the job.

## Rollback a job
In case you want to rollback an update (for example: you notice a serious issue with your new containers and want to return to the previous state as soon as possible), you can do so with the `dokkup rollback job` command:
```shell
$ dokkup rollback job demo.yaml
```
```
Deployment summary:

NAME     IMAGE                    RESTART     COUNT     GROUP     NETWORK
demo     crccheck/hello-world     always      2         labs      bridge

Node statuses:

NAME     STATUS     CONTAINERS     ROLLBACK     VERSION
lab1     ONLINE     2/2            true         9470cdc -> 55dab35
lab2     ONLINE     2/2            true         9470cdc -> 55dab35

Are you sure you want to proceed? (y/n) 
```

In case you've never done an update, the rollback field will be set to false, meaning that dokkup will not be able to do a rollback as there's no previous state to return to. 

## Stop a job
If you wish to stop a job, you can do so with the `dokkup stop job` command:
```shell
$ dokkup stop job demo.yaml
```
```
Stop job summary:

NAME     IMAGE                    GROUP
demo     crccheck/hello-world     labs

Node statuses:

NAME     STATUS     CONTAINERS     PURGE
lab1     ONLINE     2 -> 0         false
lab2     ONLINE     2 -> 0         false

Are you sure you want to proceed? (y/n) 
```

The CLI tool will show how many containers it's going to stop. Keep in mind that this will not remove the containers, it will only stop them (the PURGE flag is set to false). If you wish to remove the containers as well, take a look at [Remove/purge](#removepurge-a-job).

## Remove/purge a job
If you wish to completely remove a job from your cluster of servers, you can do so with the `dokkup stop job --purge demo.yaml`:
```shell
$ dokkup stop job --purge demo.yaml
```
```
Stop job summary:

NAME     IMAGE                    GROUP
demo     crccheck/hello-world     labs

Node statuses:

NAME     STATUS     CONTAINERS     PURGE
lab1     ONLINE     2 -> 0         true
lab2     ONLINE     2 -> 0         true

Are you sure you want to proceed? (y/n) 
```

The --purge flag will tell the `agent` to delete the containers after stopping them first.

Note: The rollback command will not undo a purge, so be careful when running this command. The only way to return the containers is by running a deployment from scratch.

# Show job status 
If you would like to get some information about a job, you can do so with the `dokkup show job` command:

```
$ dokkup show job demo.yaml
```
```
NODE     LOCATION              STATUS     JOB      IMAGE                    CONTAINERS     VERSION
lab1     192.168.0.99:8080     ONLINE     demo     crccheck/hello-world     2/2            9470cdc

CONTAINER ID     NAME                                          STATUS                     PORTS
99294b42a89d     demo-bbf84ca7-0a2b-4c65-9f39-f77fb8ce3624     Up 8 minutes (healthy)     0.0.0.0:32773->8000/tcp
28de75e833c0     demo-a9e8525a-103a-4749-b441-df676976b62d     Up 8 minutes (healthy)     0.0.0.0:32772->8000/tcp


NODE     LOCATION               STATUS    JOB      IMAGE                    CONTAINERS     VERSION
lab2     192.168.0.100:8080     ONLINE    demo     crccheck/hello-world     2/2            9470cdc

CONTAINER ID     NAME                                          STATUS                     PORTS
dc5fd3c210ed     demo-53fec61f-3eb6-4666-86d4-1d96b22d1291     Up 8 minutes (healthy)     0.0.0.0:32773->8000/tcp
9d51076ddd5a     demo-28949bc3-fc57-46b0-8c36-86b06d8a2228     Up 8 minutes (healthy)     0.0.0.0:32772->8000/tcp
```

This will show the nodes which the job is running on, along with the job name, image, how many containers are running and the version of the job. It will also display some basic information about each container, such as the ID, name, status and the exposed ports.
