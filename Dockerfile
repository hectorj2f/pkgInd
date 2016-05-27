FROM ubuntu:latest
MAINTAINER hectorj@gmail.com

COPY ./bin/pkgind /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/pkgind"]
