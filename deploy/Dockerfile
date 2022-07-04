FROM golang:1.18 AS base
MAINTAINER Levin


ADD ./ /app/
WORKDIR /app
RUN cd /app && go build

ENTRYPOINT ./matching-story-robot-service -env=$RUNENV
