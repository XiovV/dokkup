# Stop a job
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
