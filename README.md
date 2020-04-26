# Risotto Play 🍲

![License](https://img.shields.io/github/license/risotto/play)
![Go version](https://img.shields.io/github/go-mod/go-version/risotto/play)
![Uptime](https://img.shields.io/uptimerobot/ratio/m784827227-5b49008ca1a00d6181932718)
![Publish image and Deploy](https://github.com/risotto/play/workflows/Publish%20image%20and%20Deploy/badge.svg)
[![codecov](https://codecov.io/gh/risotto/play/branch/master/graph/badge.svg)](https://codecov.io/gh/risotto/play)
[![Docker image size](https://img.shields.io/docker/image-size/jjhaslanded/risotto-play?sort=date)](https://hub.docker.com/r/jjhaslanded/risotto-play)
[![Total alerts](https://img.shields.io/lgtm/alerts/g/risotto/play.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/risotto/play/alerts/)

[api.play.risotto.dev](https://api.play.risotto.dev)

Similar to [Go playground](https://play.golang.org/), but for a language nobody will really use: [Risotto](https://github.com/risotto/risotto).

Play about with it on [play.risotto.dev](https://play.risotto.dev)

This is just the API, the front-end is hosted in [risotto-play](https://github.com/risotto/play-ui), but you can just whack this api into any old front end if you so wish.

## Run tests

```bash
go test pkg/**/*
```

Or with docker:

```bash
docker build --target tester -t risotto-play-tester .
docker run --rm -iv${PWD}:/host-volume risotto-play-tester
```

## Run locally

Build the DockerFile:

```bash
docker build -t risotto-play-api -f Dockerfile .
```

Run the Dockerfile:

```bash
docker run -p 80:4000 risotto-play-api
```

Or, just run the [hosted docker image](https://hub.docker.com/r/jjhaslanded/risotto-play):

```bash
docker run -p 80:4000 jjhaslanded/risotto-play:latest
```
