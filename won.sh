#!/bin/bash

GOOS=windows 
GOARCH=amd64 
go build -o bin/my-desktop-agent-windows-amd64.exe ./...
