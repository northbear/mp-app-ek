# MP-APP-EK

Run minimal microservice application in docker-compose cluster and run ElasticSearch and Kibana stack to monitor
the microservice application.

The applicaiton  micro-services are written on Go, using web framework Gin2. It's not required any compilers or tools
preinstalled except one denoted in 'Prerequisites' section.

### Prerequisites

It requires Linux system that is capable to run Docker Engine version ~>18.x and docker-compose version ~>3.x.
It requires being installed essential build tools (that insluded `gmake`).

The host should have access image repositories:
- docker.io
- docker.elasticsearch.com

### Build images

The application consist of two services - `web`, and `auth`, that are located in respective directories.

To build images get into respective directories and run command `make build-image`. Ensure that respecitve
images are created successfully by command `docker images. The created images will be prefixed with application
name `mp-app-ek`. At this stage you are ready to run the application.

### Starting The Applcation

Run the application with command `docker-compose -f docker-compose-application.yml up`.
To run it in detached mode add options `-d` to the end of the command.


### Monitoring ES-Kibana Stack

The Kibana UI interface will be available on the host port 5601.

### Known Flaws

* Image builds can be integrated to docker-compose. So the images will be built on-fly
* Using gin2 framefork can be considered as overkill for this case. There is more simple frameworks and
libraries to handle microservice logics, session, auths
* Basic logs and logs gin framework are not unified (gin framework adds `[GIN]` prefix). It may make
log processing more comlicated
* There is requied manual definitions for indecies and filters in Kibana UI for getting useful data
and metrics
* README.md for web and auth services describing API are not provided. :(

## Reference

https://www.sarulabs.com/post/5/2019-08-12/sending-docker-logs-to-elasticsearch-and-kibana-with-filebeat.html
