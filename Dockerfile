# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /cata
COPY . .
COPY gitconfig /etc/gitconfig

RUN rm /bin/sh && ln -s /bin/bash /bin/sh

RUN apt-get update
RUN apt-get install -y protobuf-compiler
RUN go get -u google.golang.org/protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

ENV NODE_VERSION=20.11.1
ENV NVM_DIR="/root/.nvm"

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.38.0/install.sh | bash \
    && . $NVM_DIR/nvm.sh \
    && nvm install $NODE_VERSION \
    && nvm alias default $NODE_VERSION \
    && nvm use default

#RUN . "$NVM_DIR/nvm.sh" && nvm install ${NODE_VERSION}
#RUN . "$NVM_DIR/nvm.sh" && nvm use v${NODE_VERSION}
#RUN . "$NVM_DIR/nvm.sh" && nvm alias default v${NODE_VERSION}

ENV PATH="/root/.nvm/versions/node/v${NODE_VERSION}/bin/:${PATH}"

EXPOSE 8080/tcp
