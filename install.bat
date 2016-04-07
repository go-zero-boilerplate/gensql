@echo off
go-build-git -out "$GOPATH\bin\gensql.exe" -injectvar "main.GitHash"