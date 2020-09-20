FROM golang:alpine as builder
RUN apk update && apk upgrade && apk add --no-cache \
    git \
    gcc \
    libc-dev
WORKDIR /hummus
COPY . .
RUN go get ./...
RUN chmod +x scripts/prepare_release.sh
RUN ./scripts/prepare_release.sh
RUN go build -o bin/hummus
WORKDIR /
ENV PATH="/hummus/bin:${PATH}"
ENTRYPOINT ["/hummus/bin/hummus"]