# `dokkup`: Easy Container Updates

`dokkup` is a simple tool used for updating containers which can also handle rollbacks.

It's ideal for use cases when you want to update a container on demand. It communicates with 
the `dokkup-agent` via HTTP, so it doesn't require `docker` to be installed on your local machine, making it
a good fit for CI/CD.

# Setup

## Install

### Install from [Releases](https://github.com/XiovV/dokkup/releases)
To be added.

### Build and install from Source
With Go 1.16+, build and install the latest released version:

```
go install github.com/XiovV/dokkup@latest
```

# Update a Container
```shell
dokkup up -node http://localhost:8080 -container containerExample -tag latest -api-key [YOUR_API_KEY] -keep
```

### Command line flags

- `node` specifies where the `dokkup-agent` is located.
- `container` specifies the container you'd like to update.
- `tag` specifies the image tag (in our case, it will take whatever image `containerExample` is running on, and it will pull the `latest` image)
- `api-key` specifies your `dokkup-agent`'s API Key,
- `keep` (optional) will tell `dokkup-agent` to keep the old version of the container, this is useful if you want to use the `rollback` command.
- `image` specifies the full image (e.g. imageExample:latest). Note: you can either only use `tag` or `image`, using both is not allowed.

# Rollback a Container

In this example, `rollback` will look for a container named `containerExample-rollback` (which should exist if you had used the `keep` flag with the `up` command), if it manages to find it,
it will stop and remove `containerExample` and then it will run `containerExample-rollback`. 
```shell
dokkup rollback -node http://localhost:8080 -container containerExample -api-key [YOUR_API_KEY]
```

# Configuration
It's possible to use a config file instead of just using command line flags. You can configure anything using
a `dokkup.yaml` file. `dokkup` will look for this file and load it (the file needs to be named `dokkup.yaml`, at the moment it's 
still not possible to use a custom filename). Please note however that command line flags have a higher priority, meaning that `dokkup` will ignore the `dokkup.yaml`
file if you run it with command line flags. Example:
```yaml
api_key: '[YOUR_API_KEY]'
node_location: 'http://localhost:8080'
container: 'containerExample'
tag: 'latest'
# image: 'containerExample:latest' (we are using tag in this 
# example, so it's unnecessary to specify the full image.)
keep: true
```