# Risotto Play

Similar to [Go playground](https://play.golang.org/), but for a language nobody will really use: [Risotto](https://github.com/raphaelvigee/risotto)

This is just the API, the front-end for this will be hosted in another repository...

## Run tests

```bash
go test pkg/**/*
```

## Run locally

Build the DockerFile:

```bash
docker build -t risotto-play-api -f Dockerfile .
```

Run the Dockerfile:

```bash
docker run -p 4000:4000 risotto-play-api
```

Or, just run the [hosted docker image](https://hub.docker.com/r/jjhaslanded/risotto-play):

```bash run -p 80:4000 jjhaslanded/risotto-play:latest
