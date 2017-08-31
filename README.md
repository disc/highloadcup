# Highload Cup 2017

## Stack
Go

## Run
```
docker build -t golang-app .
docker run --rm -p 8080:80 -v $(pwd)/data.zip:/tmp/data/data.zip -t golang-app
```