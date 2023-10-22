# Agent setup
The agent runs on your server and will do all of the container management once it receives a request from the CLI tool.

## Deploy the agent

Running the agent is as simple as running a simple docker command:

```shell
$ docker run -d --name dokkup-agent -p 8080:8080 \ 
--restart always -v /path/to/config:/config \ 
-v /var/run/docker.sock:/var/run/docker.sock xiovv/dokkup:latest
```

## Retrieve the API key
You can retrieve the API key through the logs. Please note that the API key will not be printed out the next time you run the agent.

```shell
$ docker logs dokkup-agent
```
```shell
Your new API key is: jy9DbtDlfi5VJuAkbZYd4Kt0c2cQY8iQ
2023-09-22T10:31:42.175+0200    INFO    agent/main.go:36        server is listening...  {"port": "8080"}
```

The agent is now ready to be used.
