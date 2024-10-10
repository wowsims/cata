# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /cata

SHELL ["/bin/bash", "-c"]

COPY . .
COPY gitconfig /etc/gitconfig

# Install all Go dependencies
RUN apt-get update \
	&& apt-get install -y protobuf-compiler \
	&& go get -u google.golang.org/protobuf \
	&& go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
	&& curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin

ENV NODE_VERSION=20
ENV NVM_DIR="/root/.nvm"

# Install all Frontend dependencies
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash \
	&& . $NVM_DIR/nvm.sh \
	&& nvm install $NODE_VERSION \
	&& nvm alias default $NODE_VERSION \
	&& nvm use default

ENV PATH="/root/.nvm/versions/node/v${NODE_VERSION}/bin/:${PATH}"

EXPOSE 8080 3333 5173
