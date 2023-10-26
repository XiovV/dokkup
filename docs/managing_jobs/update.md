# Update a job

Updating a job is as simple as making a change in the job specification file and running the `dokkup run job` command again: 

```yaml title="demo.yaml"
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
```

!!! note
    Dokkup will execute the update with zero downtime.

    In the event that something goes wrong during the update, dokkup will abort and return the previous containers, ensuring minimum downtime in the event of an error.

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

The update status is now true, meaning that dokkup is going to take down the currently running containers and deploy new ones. The new version hash is also shown. 
