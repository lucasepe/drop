#!/bin/bash


go test -count=1 -p 1 -cover -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
