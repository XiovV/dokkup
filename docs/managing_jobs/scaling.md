# Scaling jobs 
With dokkup's scaling feature, you can easily change the number of running containers on your nodes.

!!! note
    Upscaling or downscaling a job will not have an effect on your uptime.

## Upscale a job
If you want to increase the number of running containers, you can do so by incrasing the `count` field in your job specification, then run the `dokkup run job` command: 

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

## Downscale a job
If you want to decrease the number of running containers, you can do so by decreasing the `count` field in your job specification, then run the `dokkup run job` command: 

```yaml title="demo.yaml"
group: "labs"
# count: 5
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
```
```shell
$ dokkup run job demo.yaml
```
```
Deployment summary:

NAME     IMAGE                    RESTART     COUNT     GROUP     NETWORK
test     crccheck/hello-world     always      2         local     bridge

Node statuses:

NAME          STATUS     CONTAINERS     UPDATE     VERSION
lab1          ONLINE     5 -> 2         true       55dab32
lab2          ONLINE     5 -> 2         true       55dab32

Are you sure you want to proceed? (y/n) 
```
