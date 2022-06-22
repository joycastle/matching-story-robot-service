#!/bin/bash
protoc --proto_path=./ --go_out=:. --go_opt=paths=source_relative ./*.proto
protoc --proto_path=./ --go_out=plugins=grpc:. --go_opt=paths=source_relative ./pigeon.proto

