# Upscale or downscale a job
In case you want to change the number of running containers, you can do so by changing the `count` field in your job specification, then run the `dokkup run job` command: 

```yaml title="demo.yaml"
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
