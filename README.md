# dockerlogs

A CLI for streaming logs from a host machine using a priviliged container daemon


```bash
dockerlogs -h
# Usage: dockerlogs [options] [command]

#Options:
#  -h, --help         output usage information

#Commands:
#  status             Check if the daemon is runnning
#  list               Print the available container log records
#  stream <streamId>  Streams logs for the specified streamId
```

## Getting Started
These instructions will get you a copy of the project up and running on your local machine for demonstration purposes.

## Prerequisites

What things you need to install the software and how to install them

* docker installed locally
* NodeJS (install easily with NVM on unix machines)

## Setup

### Building and running the Docker Daemon

```bash
# Build the docker-logs-daemon image
cd docker-logs-daemon
docker build -t docker-logs-daemon .

# Run the Daemon to start collecting the logs
docker run --rm -d \
--name docker-logs-daemon \
--privileged \
-e DOCKER_HOST="unix:///var/run/docker.sock" \
-v /var/run/docker.sock:/var/run/docker.sock \
-v $(pwd)/logs:/app/logs \
-p 8080:8080 \
docker-logs-daemon
```

### Building and running the example node app


```bash
# cd into the example-app directory
cd example-node-logs-app

# Build the image
docker build -t node_logs_app .

# Run the app with two different labels -
# indicating only one of them should have it's logs collected.
# We will use the "collect_logs" label for now.

docker run --rm -d \
-l collect_logs=false \
-e APP_NAME="Example app (logs NOT collected)" \
--name node_example_app_without_logs \
node_logs_app

docker run --rm -d \
-l collect_logs=true \
-e APP_NAME="Example app (logs collected)" \
--name node_example_app_with_logs \
node_logs_app

# Confirm only one of the containers is marked for log collection
docker ps -f "label=collect_logs=true"
```

## Using the CLI

```bash
# Install the CLI
cd docker-logs-cli
npm i -g docker-logs-cli/

# Now you can invoke the CLI directly
dockerlogs --help
```

## Cleanup


```bash
# Kill the example docker images
docker kill $(docker ps -q -f "name=node_example")
```

### Future Improvements

* Make the daemon docker image thin (use builder image for build and scratch for runtime)
* Find a way to not give the daemon root permissions (maybe use a docker unix sock permission proxy)

### Installing

A step by step series of examples that tell you how to get a development env running

Say what the step will be

```
Give the example
```

And repeat

```
until finished
```

End with an example of getting some data out of the system or using it for a little demo

## Known Issues

* When deleting a specific log from the storage folder - you have to re-run the daemon - it won't re-create the file on error



## Built With

* Golang
* NodeJS
* Docker Engine SDK

## Authors

* **Eli oldberg**

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc
