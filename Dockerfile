# Dockerfile
FROM ubuntu:latest
MAINTAINER ish <ish@innogrid.com>

RUN mkdir -p /GraphQL_piano/
WORKDIR /GraphQL_piano/

ADD GraphQL_piano /GraphQL_piano/
RUN chmod 755 /GraphQL_piano/GraphQL_piano

EXPOSE 8001

CMD ["/GraphQL_piano/GraphQL_piano"]
