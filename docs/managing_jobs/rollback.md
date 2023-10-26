# Rollback a job
In case you want to rollback an update (for example: you notice a serious issue with your new containers and want to return to the previous state as soon as possible), you can do so with the `dokkup rollback job` command:

!!! note 
    If you've never done an update of a job before, the rollback field will be set to false, meaning that dokkup will not be able to do a rollback as there's no previous state to return to.
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
