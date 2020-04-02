# Risotto Play 🍲

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
