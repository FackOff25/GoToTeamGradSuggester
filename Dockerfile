FROM ubuntu:latest

COPY ./bin/ /suggest/
COPY ./configs/config.yaml /suggest/

WORKDIR /suggest
