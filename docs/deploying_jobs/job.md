# Job

The job specification contains the configuration for the containers you want to deploy. It closely resembles a standard [docker-compose](https://docs.docker.com/compose) file.

```yaml title="job.yaml"
# The group where you want to deploy the job.
#
# Cannot be used if node OR nodes have been specified.
group: "labs"

# The node on which you want to deploy the job.
#
# Cannot be used if group OR nodes have been specified.
# node: "lab1"

# List of nodes on which you want to deploy the job.
#
# Cannot be used if group OR node have been specified.
# nodes: ["lab1", "lab2"]

# How many containers you want to run per node.
count: 2

# Name of the job.
name: "demo"

container:
   # Docker image
 - image: "traefik/whoami:latest"

   # Array of ports
   ports:
      # Port inside the container
    - in: 8000
      # Exposed port
      #
      # If omitted docker will dynamically assign the exposed port.
      # Recommended to omit if running more than one instance of a container (e.g. count is greater than 1)
    - out: 8080

   # Restart policy
   restart: always 

   # Array of labels
   labels:
     - "my.label.test=demo"
   
   # Array of environment variables
   environment:
     - MYENV=ENVEXAMPLE
    
   # Array of volumes
   volumes:
     - myvolume:/home
   
   # Array of networks
   #
   # Default: bridge
   networks:
    - mynetwork

   # Array of command flags
   commands:
    - "--my-custom-flag"
```
