FROM jenkins/slave:alpine

LABEL maintainer="Stephen Lane-Walsh <sdl.slane@gmail.com>"

USER root
ARG APP_PATH=github.com/WhoBrokeTheBuild

ENV GOPATH=/opt/go
ENV AGENT_WORKDIR=${GOPATH}/src/${APP_PATH}

RUN apk update && \
    apk add --update --no-cache go make mesa-dev glu-dev libc-dev libxcursor-dev libxrandr-dev libxinerama-dev libxi-dev && \
    mkdir -p /opt/go/src/github.com/whoBrokeTheBuild
