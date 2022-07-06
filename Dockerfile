FROM golang:1.18 AS base
ADD ./ /app/
WORKDIR /app
RUN cd /app && go build
#ENTRYPOINT ["./matching-story-robot-service", "-env", "dev"]
