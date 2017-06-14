#!/bin/sh
go build -o px86 main.go
./px86 ../tolset_p86/exec-helloworld/helloworld.bin
