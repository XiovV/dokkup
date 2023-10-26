# Stop a job
!!! warning
    Stopping a job will NOT remove the containers, it will simply stop them. If you'd like to remove the containers, read [Remove/purge](remove.md).

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

!!! note
    After stopping a job, you can easily start it back up again with the `dokkup run job` command.
