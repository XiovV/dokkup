# Getting started

## What is dokkup?
dokkup is a container orchestrator designed to make managing containers across a cluster of servers as simple as possible.

It's ideal for use cases where you want to orchestrate containers accross one or multiple servers. Think of it as Nomad, Swarm or Kubernetes.

## How does it work?
dokkup consists of two parts, the agent and the CLI tool:

- **agent** runs on the servers/nodes where you want to run your containers. It listens for incoming job requests and executes them.
- **CLI** is used to dispatch job requests to the servers/nodes.

The CLI requires two files, the inventory and the job files:

- **inventory** contains information about your nodes, such as their agent's API key, URL and name.
- **job** contains the configuration for the containers you want to deploy, it closely resembles a standard [docker-compose](https://docs.docker.com/compose) file.

## Where do I start?
You can start by reading [Agent setup](agent-setup.md) and [CLI setup](CLI-setup.md) then reading through the Deploying jobs section.
