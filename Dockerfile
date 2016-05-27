FROM ubuntu:latest
MAINTAINER hectorj@gmail.com

COPY ./bin/pkgInd /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/pkgInd"]
