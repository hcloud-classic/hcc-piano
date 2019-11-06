# Dockerfile
FROM ubuntu:latest
MAINTAINER ish <ish@innogrid.com>

RUN mkdir -p /piano/
WORKDIR /piano/

ADD flute /piano/
RUN chmod 755 /piano/piano

EXPOSE 7300

CMD ["/piano/piano"]
