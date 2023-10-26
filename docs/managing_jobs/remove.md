# Remove/purge a job
If you wish to completely remove a job from your cluster of servers, you can do so with the `dokkup stop job --purge demo.yaml`:

!!! warning
    The rollback command cannot undo a purge, so be careful when running this command. The only way to return the containers is by running a deployment from scratch.

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
