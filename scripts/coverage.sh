#!/bin/bash


go test -count=1 -p 1 -cover -coverprofile=coverage.txt ./...
go tool cover -func=coverage.txt
