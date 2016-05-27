[![Build Status](https://api.travis-ci.org/hectorj2f/pkgInd.svg)](https://travis-ci.org/hectorj2f/pkgInd)
[![](https://godoc.org/github.com/hectorj2f/pkgInd?status.svg)](http://godoc.org/github.com/hectorj2f/pkgInd)
[![](https://img.shields.io/docker/pulls/hectorj2f/pkgind.svg)](http://hub.docker.com/hectorj2f/pkgind
[![Go Report Card](https://goreportcard.com/badge/github.com/hectorj2f/pkgInd)](https://goreportcard.com/report/github.com/hectorj2f/pkgInd)

# Package Indexer

Clients will connect to this package Indexer and inform which packages should be indexed,
and which dependencies they might have on other packages.

Messages from clients follow this pattern:

`<command>|<package>|<dependencies>\n`

Where:

* `<command>` is mandatory, and is either `INDEX`, `REMOVE`, or `QUERY`

* `<package>` is mandatory, the name of the package referred to by the command, e.g. `mysql`, `openssl`, `pkg-config`, `postgresql`, etc.

* `<dependencies>` is optional, and if present it will be a comma-delimited list of packages that need to be present before `<package>` is installed. e.g. `cmake,sphinx-doc,xz`

* The message always ends with the character `\n`


## Run using pkgIndctl

INDEX: `pkgIndctl index --package=cloog --dependencies=a,b,c`

REMOVE: `pkgIndctl index --package=cloog`

QUERY: `pkgIndctl query --package=cloog`

Here are some sample messages:

```
INDEX|cloog|gmp,isl,pkg-config\n
INDEX|ceylon|\n
REMOVE|cloog|\n
QUERY|cloog|\n
```


## Run within a container

DOCKER: `docker run --rm -p 8080:8080 hectorj2f/pkgind start`

RKT: `rkt run --port=host:8080 docker://hectorj2f/pkgind --exec=/usr/local/bin/pkgInd -- start`
