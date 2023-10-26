# Reverse Proxy & Load Balancing with TLS using Traefik

This example is going to show you step-by-step how to deploy multiple containers on a node, and how to reverse proxy and load balance them (TLS included) via Traefik.

Make sure you read [Getting started](../index.md) before continuing.


## Inventory file
Learn more about inventory files [here](../deploying_jobs/inventory.md).

The first step is to create an inventory file, so dokkup can know about our node:

!!! Warning
    Make sure to use your own IP:PORT and agent API key!

```yaml title="inventory.yaml"
nodes:
 - name: "lab1"
   location: "192.168.0.99:8080"
   key: "EcIGwgiWhA4K7dDM8mhTSAEVU74PS3CI"

groups:
 - name: "labs"
   nodes: ["lab1"]
```

## Traefik setup 
### Traefik job specification
Learn more about job specifications [here](http://localhost:8000/deploying_jobs/job).

Let's write the job specification for Traefik:

```yaml title="traefik.yaml"
group: "labs"
count: 1
name: "traefik"

container:
 - image: "traefik:latest"
   ports:
    - in: 80
      out: 80
    - in: 8080
      out: 8081
   restart: always 
   volumes:
     - "/var/run/docker.sock:/var/run/docker.sock:ro"
   command:
       # Do not use in production!
     - "--api.insecure=true"
     - "--providers.docker=true"

     - "--entrypoints.web.address=:80"
```

We only want one instance of Traefik, so we set the `count` value to 1. We have also exposed the ports 80 (for http) and 8081 (for the Traefik dashboard), and set some basic Traefik flags.

### Deploy traefik
Now that we have our Traefik job specification ready, it's time to deploy it:

```
$ dokkup run job traefik.yaml
```

```
Deployment summary:

NAME        IMAGE              RESTART     COUNT     GROUP     NETWORK
traefik     traefik:latest     always      1         labs      [bridge]

Node statuses:

NAME      STATUS     CONTAINERS     UPDATE     VERSION
lab1*     ONLINE     0 -> 1         true       72a16a1

Are you sure you want to proceed? (y/n) 
```

Now you can go to 192.168.0.99:8081 to make sure Traefik is running.

## Container setup
In this example we will be deploying traefik's [whoami](https://hub.docker.com/r/traefik/whoami) image, but you can use any image you'd like. 

### Whoami job specification
Let's write the job specification for whoami:
```yaml title="whoami.yaml"
group: "labs"
count: 3
name: "demo"

container:
 - image: "traefik/whoami"
   restart: always 
   ports:
    - in: 80
   labels:
     - "traefik.enable=true"
     - "traefik.http.routers.whoami.rule=Host(`whoami.example.com`)"

     - "traefik.http.routers.whoami.entrypoints=web"
     - "traefik.http.routers.whoami.service=whoami"
     - "traefik.http.services.whoami.loadbalancer.server.port=80"
```

Here we are telling dokkup to deploy 3 instances of this container, and we are setting labels for Traefik.

Don't worry! We will set up TLS soon!

### Deploy whoami
!!! warning
    Make sure you've set up your DNS records before deploying!

Let's deploy our containers:
```
dokkup run job whoami.yaml
```

```
Deployment summary:

NAME     IMAGE              RESTART     COUNT     GROUP     NETWORK
demo     traefik/whoami     always      3         labs      [bridge]

Node statuses:

NAME      STATUS     CONTAINERS     UPDATE     VERSION
lab1*     ONLINE     0 -> 3         true       6ac3f7b

Are you sure you want to proceed? (y/n) 
```

Now run `curl` to make sure it's working:
```
$ curl http://whoami.example.com
```

```
Hostname: 3d23e4572caf
IP: 127.0.0.1
IP: 172.17.0.4
RemoteAddr: 172.17.0.2:53158
GET / HTTP/1.1
Host: whoami.example.com
...
```

Excellent! It seems to be working! If you keep running `curl`, you will see the Hostname changing on each request, that means our containers are being load balanced properly. It's time to set up TLS.


## Setting up TLS
In order to set up TLS, we first have to update our Traefik and whoami job specifications.

### Traefik configuration
Read Traefik's official [documentation](https://doc.traefik.io/traefik/https/acme) to learn more.

```yaml title="traefik.yaml"
# ...
container:
   ports:
   # ...
    - in: 443
      out: 443
   volumes:
   #  ...
     - "letsencrypt:/letsencrypt"
   command:
   #  ...
     - "--entrypoints.web-secure.address=:443"

     - "--certificatesresolvers.le.acme.httpchallenge=true"
     - "--certificatesresolvers.le.acme.httpchallenge.entrypoint=web"

     - "--certificatesresolvers.le.acme.email=email@example.com"
     - "--certificatesresolvers.le.acme.storage=/letsencrypt/acme.json"
```

Here we've exposed the port 443 (for https), added a volume for Let's Encrypt, and added command flags for TLS.

### Update Traefik
Now let's update Traefik with our new settings:

```
$ dokkup run job traefik.yaml
```

```
Deployment summary:

NAME        IMAGE              RESTART     COUNT     GROUP     NETWORK
traefik     traefik:latest     always      1         labs      [bridge]

Node statuses:

NAME     STATUS     CONTAINERS     UPDATE     VERSION
lab1     ONLINE     1/1            true       72a16a1 -> 65f5458

Are you sure you want to proceed? (y/n) 
```

Traefik should be set now, let's update our whoami job specification!

### Whoami configuration
```yaml title="whoami.yaml"
# ...
container:
   # ...
   labels:
     # ...
     - "traefik.http.middlewares.myredirect.redirectscheme.scheme=https"
     - "traefik.http.routers.whoami.middlewares=myredirect"

     - "traefik.http.routers.whoami-secure.rule=Host(`whoami.example.com`)"
     - "traefik.http.routers.whoami-secure.entrypoints=web-secure"
     - "traefik.http.routers.whoami-secure.tls.certresolver=le"
     - "traefik.http.routers.whoami-secure.tls=true"
```

Now it's time to update our whoami containers!

### Update whoami
```
$ dokkup run job whoami.yaml
```

```
Deployment summary:

NAME     IMAGE              RESTART     COUNT     GROUP     NETWORK
demo     traefik/whoami     always      3         labs      [bridge]

Node statuses:

NAME     STATUS     CONTAINERS     UPDATE     VERSION
lab1     ONLINE     3/3            true       6ac3f7b -> e909b4b

Are you sure you want to proceed? (y/n)
```

### Testing TLS
Now that we have made the necessary changes for TLS to work, let's test it out. We will run `curl` multiple times to make sure load balancing is working.

```
$ curl https://whoami.example.com
```

```
Hostname: 828f6cad787f
IP: 127.0.0.1
IP: 172.17.0.6
RemoteAddr: 172.17.0.2:54178
GET / HTTP/1.1
```

```
$ curl https://whoami.example.com
```

```
Hostname: c638a325519e
IP: 127.0.0.1
IP: 172.17.0.5
RemoteAddr: 172.17.0.2:38422
GET / HTTP/1.1
```

```
$ curl https://whoami.example.com
```

```
Hostname: 3d23e4572caf
IP: 127.0.0.1
IP: 172.17.0.4
RemoteAddr: 172.17.0.2:50274
GET / HTTP/1.1
```

The Hostname value is different each time we make a request, which shows that load balancing is working properly! 

And that is it, you have successfully used dokkup to deploy multiple containers and put them behind a load balancer, and all of that with free TLS! 

Now you can use [scaling](../managing_jobs/scaling.md) to add or remove containers, and they will be automatically load balanced by Traefik!
