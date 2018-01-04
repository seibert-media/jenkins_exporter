FROM quay.io/prometheus/busybox:latest
MAINTAINER Kevin Wiesmueller <kwiesmueller@seibert-media.net>

COPY jenkins_exporter /bin/jenkins_exporter

EXPOSE 9103
ENTRYPOINT ["/bin/jenkins_exporter"]
