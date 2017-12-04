# Jenkins Exporter

[![Build Status](http://github.dronehippie.de/api/badges/exporters/jenkins/status.svg)](http://github.dronehippie.de/exporters/jenkins)
[![Go Doc](https://godoc.org/github.com/exporters/jenkins?status.svg)](http://godoc.org/github.com/exporters/jenkins)
[![Go Report](http://goreportcard.com/badge/github.com/exporters/jenkins)](http://goreportcard.com/report/github.com/exporters/jenkins)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b18de3d4c8034064947ed1daaf223c62)](https://www.codacy.com/app/exporters/jenkins?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=exporters/jenkins&amp;utm_campaign=Badge_Grade)
[![](https://images.microbadger.com/badges/image/exporters/jenkins.svg)](http://microbadger.com/images/exporters/jenkins "Get your own image badge on microbadger.com")
[![Join the Matrix chat at https://matrix.to/#/#prometheus-exporters:matrix.org](https://matrix.to/img/matrix-badge.svg)](https://matrix.to/#/#prometheus-exporters:matrix.org)

[Prometheus](https://prometheus.io/) exporter that collects Jenkins metrics.


## Installation

If you are missing something just write us on our nice [Matrix](https://matrix.to/#/#prometheus-exporters:matrix.org) or [IRC](https://webchat.freenode.net/?channels=prometheus-exporters) channel. If you want to use our pre-built binaries just head over to the [releases](https://github.com/exporters/jenkins/releases).


### Docker

```
docker pull exporters/jenkins:0.1
docker run --rm -p 9212:9212 -e DEBUG=true -e JENKINS_ADDRESS=http://jenkins.example.com -e JENKINS_USERNAME=username -e JENKINS_PASSWORD=p455w0rd exporters/jenkins:0.1
```

If you are using docker-compose to orchestrate your Prometheus setup you can use this quite simple snippet:

```
jenkins_exporter:
    image: exporters/jenkins:0.1
    restart: always
    environment:
    - DEBUG=true
    - JENKINS_ADDRESS=http://jenkins.example.com
    - JENKINS_USERNAME=username
    - JENKINS_PASSWORD=p455w0rd
    ports:
    - "127.0.0.1:9212:9212"
```


## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| DEBUG | `false` | If set to true also debug information will be logged, otherwise only info |
| JENKINS_ADDRESS |  | Address where we can find the Jenkins server |
| JENKINS_USERNAME |  | Username for the server authentication |
| JENKINS_PASSWORD |  | Password for the server authentication |
| WEB_ADDR | `:9212` | Address for this exporter to run |
| WEB_PATH | `/metrics` | Path for metrics  |


## Metrics


| Name | Type | Cardinality | Help |
|------|------|-------------|------|
| jenkins_ | gauge | 1 | Description |


### Alerts & Recording Rules

Example alerts and recording rules you can find at [example.rules](example.rules)


## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). It is also possible to just simply execute the `go get github.com/exporters/jenkins` command, but we are mostly using our `Makefile`. Make sure you copy the `.env.example` to `.env` and change the variables matching your credentials.

```
go get -d github.com/exporters/jenkins
cd $GOPATH/src/github.com/exporters/jenkins
make build && ./jenkins_exporter
```


## Contributing

Fork -> Patch -> Push -> Pull Request


## Authors

* [Thomas Boerger](https://github.com/tboerger)


## License

Apache-2.0


## Copyright

```
Copyright (c) 2017 Thomas Boerger <thomas@webhippie.de>
```
