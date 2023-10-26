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
