FROM quay.io/prometheus/busybox:latest
MAINTAINER //SEIBERT/MEDIA GmbH <docker@seibert-media.net>

COPY jenkins_exporter /bin/jenkins_exporter

EXPOSE 9103
ENTRYPOINT ["/bin/jenkins_exporter"]
