FROM golang:1.18 AS base
MAINTAINER Levin

ADD ./ /app/
WORKDIR /app
RUN cd /app && go build

