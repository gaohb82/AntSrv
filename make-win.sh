#! /bin/bash 
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o ants.exe -ldflags -H=windowsgui main.go