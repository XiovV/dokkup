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

NAME      STATUS     CONTAINERS     UPDATE
lab1*     ONLINE     0/0            true
lab2*     ONLINE     0/0            true

Are you sure you want to proceed? (y/n) 
```
- Tip 1: If there's an `inventory.yaml` in your current directory, you can omit the -i flag, dokkup loads `inventory.yaml` files by default. 
- Tip 2: You can provide a -y or --yes flag to skip the confirmation prompt.

The CLI will show a deployment summary, showing some basic information about the job and the container you are about to deploy. And it will display the nodes on which the job will be deployed. The asterisk next to the node name signifies that a job will be deployed from scratch.

TODO: Show gif of the deployment process here.

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

NAME     STATUS     CONTAINERS     UPDATE
lab1     ONLINE     2/2            false
lab2     ONLINE     2/2            false

Are you sure you want to proceed? (y/n) 
```
The CLI will show how many containers are running and the update status which will signify if the job is going to be updated or not.

TODO: insert gif of "already up to date" clip

# Update a job
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

NAME     STATUS     CONTAINERS     UPDATE
lab1     ONLINE     2/2            true
lab2     ONLINE     2/2            true

Are you sure you want to proceed? (y/n) 
```
Note: we omitted the -i inventory.yaml flag because dokkup automatically loads files called `inventory.yaml`.

The update status is now true, meaning that dokkup is going to take down the currently running containers and deploy new ones. If something goes wrong during the update, dokkup will abort and run the old containers, ensuring minimum downtime in case something goes wrong.

TODO: insert gif of the update process

# Rollback a job
In case you want to rollback an update, you can do so with the `dokkup rollback job` command:
```shell
$ dokkup rollback job demo.yaml
```
```
Deployment summary:

NAME     IMAGE                    RESTART     COUNT     GROUP     NETWORK
demo     crccheck/hello-world     always      2         labs      bridge

Node statuses:

NAME     STATUS     CONTAINERS     ROLLBACK
lab1     ONLINE     2/2            true
lab2     ONLINE     2/2            true

Are you sure you want to proceed? (y/n)
```

TODO: insert gif of the rollback process

In case you've never done an update, the rollback field will be set to false, meaning that dokkup will not be able to do a rollback as there's no previous state to return to. 

# Stop a job
If you wish to stop a job, you can do so with the `dokkup stop job command`:
```shell
$ dokkup stop job demo.yaml
```
```
Stop job summary:

NAME     IMAGE                    GROUP
demo     crccheck/hello-world     labs

Node statuses:

NAME     STATUS     CONTAINERS
lab1     ONLINE     2/2
lab2     ONLINE     2/2

Are you sure you want to proceed? (y/n)
```

The CLI tool will show how many containers it's going to stop. Keep in mind that this will not remove the containers, it will only stop them. If you wish to remove the containers as well, take a look at the next chapter.

# Remove/purge a job
If you wish to completely remove a job from your cluster of servers, you can do so with the `dokkup stop job --purge demo.yaml`:
```shell
$ dokkup stop job --purge demo.yaml
```
```
Stop job summary:

NAME     IMAGE                    GROUP
demo     crccheck/hello-world     labs

Node statuses:

NAME     STATUS     CONTAINERS
lab1     ONLINE     2/2
lab2     ONLINE     2/2

Are you sure you want to proceed? (y/n)
```

The --purge flag will tell the `agent` to delete the containers after stopping them first.

Warning: The rollback command will not undo a purge, so be careful when running this command. The only way to return the containers is by running a deployment from scratch.
