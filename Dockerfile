FROM ubuntu:latest

COPY ./bin/ /suggest/
COPY ./configs/config.toml /suggest/

WORKDIR /suggest
