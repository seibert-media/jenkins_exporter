FROM webhippie/alpine:latest

LABEL maintainer="Thomas Boerger <thomas@webhippie.de>" \
  org.label-schema.name="Jenkins Exporter" \
  org.label-schema.vendor="Thomas Boerger" \
  org.label-schema.schema-version="1.0"

EXPOSE 9106
ENTRYPOINT ["/bin/jenkins_exporter"]

RUN apk add --no-cache ca-certificates
COPY dist/jenkins_exporter-*-linux-amd64 /bin/jenkins_exporter
