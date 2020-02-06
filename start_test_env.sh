#!/usr/bin/env bash

if [ "" == "$(docker ps | grep dynamodb-playground)" ]; then
  echo "Starting dynamodb docker stuff"
  docker network create wojtek-learning-go
  docker run --net wojtek-learning-go -d --rm --name dynamodb-plaground -p 6000:6000 amazon/dynamodb-local --port 6000
  while ! nc -z localhost 6000; do
    echo "Could not connect to docker, is it running?"
    sleep 2;
  done
  echo "Dynamodb started!"
else
  echo "Docker already running"
fi
