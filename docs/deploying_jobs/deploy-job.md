# Deploy a job
Before we can start deploying jobs, we need an inventory and a job file. Check out [Inventory](inventory.md) and [Job specification](job.md) to learn more about them.

## Inventory
The inventory contains information about our agents. The concept is very similar to Ansible's [inventory](https://docs.ansible.com/ansible/latest/inventory_guide/intro_inventory.html) file.

```yaml title="inventory.yaml"
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

Here we've defined two nodes and assigned them to a group. This way we can refer to them with a group name instead of individually. 

Note that you do not have to use 2 or more nodes, you can use dokkup with a single node if you wish to do so.


## Job specification
The job specification contains the configuration for the containers you want to deploy. It closely resembles a standard [docker-compose](https://docs.docker.com/compose) file.

Check out [Job specification](job.md) for a full example.

```yaml title="demo.yaml"
group: "labs"
count: 2
name: "demo"

container:
 - image: "crccheck/hello-world"
   ports:
    - in: 8000
   restart: always 
```

Here we are telling dokkup:

- on which nodes we want the job to be deployed to.
- how many instances of the container we want to run per node (in our case this will total to 4 containers, as we are deploying 2 per node).
- to assign the job the name "demo". 
- information about the container, such as the image, ports and restart policy. 

In this example we have ommited assigning an exposed port. It's recommended to do so when running more than one container (e.g. count is greater than 1) because each container needs to have its own unique exposed port. Docker will automatically assign exposed ports for each of our containers.

## Deploy the job
Now that we have our inventory and job specification, we can deploy our job:


```
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

!!! tip
    - If there's an `inventory.yaml` in your current directory, you can omit the -i flag, dokkup loads `inventory.yaml` files by default.
    - You can provide a -y or --yes flag to skip the confirmation prompt.

The CLI will show a deployment summary, showing some basic information about the job, such as how many containers it's going to run and the hashed version tag of the job, along with the nodes on which the job will be deployed. The asterisk next to the node name signifies that a job will be deployed from scratch.

!!! note
    In case you run the `dokkup run job` command without making any changes, you don't have to worry about dokkup wiping your existing containers and re-deploying them again, it will detect that nothing has changed and it won't do anything.

## Show the containers
Now we can run `docker ps` on our nodes and see our containers (this is for demonstration purposes, you don't have to do this):
```
$ docker ps
```

```title="lab1"
CONTAINER ID   IMAGE                  COMMAND                   CREATED          STATUS                             PORTS                                       NAMES
adec3b3612ca   crccheck/hello-world   "/bin/sh -c 'echo \"h…"   14 seconds ago   Up 12 seconds (health: starting)   0.0.0.0:32805->8000/tcp                     demo-b905568e-942a-4ef4-b091-45f9fc2ddea9
e3a4b99761b1   crccheck/hello-world   "/bin/sh -c 'echo \"h…"   14 seconds ago   Up 12 seconds (health: starting)   0.0.0.0:32804->8000/tcp                     demo-78d37fdd-cbfc-468c-88b3-b8a1df17855b
23e0af35bdae   xiovv/dokkup:latest    "/agent"                  25 hours ago     Up 25 hours                        0.0.0.0:8080->8080/tcp, :::8080->8080/tcp   dokkup-agent
```
```title="lab2"
CONTAINER ID   IMAGE                  COMMAND                   CREATED          STATUS                    PORTS                                       NAMES
cf00cd390db3   crccheck/hello-world   "/bin/sh -c 'echo \"h…"   36 seconds ago   Up 34 seconds (healthy)   0.0.0.0:32793->8000/tcp                     demo-d71fc907-6b29-4287-aab8-d29a5fe4e821
0e8e013adea2   crccheck/hello-world   "/bin/sh -c 'echo \"h…"   37 seconds ago   Up 35 seconds (healthy)   0.0.0.0:32792->8000/tcp                     demo-4d2c8f50-48e7-4817-b380-4b790196d34f
2386f86f788b   xiovv/dokkup:latest    "/agent"                  25 hours ago     Up 25 hours               0.0.0.0:8080->8080/tcp, :::8080->8080/tcp   dokkup-agent
```

Alternatively, you can use the `dokkup show job` command to see your containers:
```
$ dokkup show job demo.yaml
```
```
NODE     LOCATION              STATUS     JOB      IMAGE                    CONTAINERS     VERSION
lab1     192.168.0.99:8080     ONLINE     demo     crccheck/hello-world     2/2            55dab35

CONTAINER ID     NAME                                          STATUS                      PORTS
adec3b3612ca     demo-b905568e-942a-4ef4-b091-45f9fc2ddea9     Up 12 seconds (healthy)     0.0.0.0:32805->8000/tcp
e3a4b99761b1     demo-78d37fdd-cbfc-468c-88b3-b8a1df17855b     Up 12 seconds (healthy)     0.0.0.0:32804->8000/tcp


NODE     LOCATION               STATUS    JOB      IMAGE                    CONTAINERS     VERSION
lab2     192.168.0.100:8080     ONLINE    demo     crccheck/hello-world     2/2            55dab35

CONTAINER ID     NAME                                          STATUS                      PORTS
cf00cd390db3     demo-d71fc907-6b29-4287-aab8-d29a5fe4e821     Up 34 seconds (healthy)     0.0.0.0:32793->8000/tcp
0e8e013adea2     demo-4d2c8f50-48e7-4817-b380-4b790196d34f     Up 35 seconds (healthy)     0.0.0.0:32792->8000/tcp
```

And that's it! You have successfully deployed 4 instances of a container spread accross 2 nodes. 

Read the Managing jobs section and [Reverse Proxy & Load Balancing with TLS using Traefik](/examples/traefik) to learn more. 
