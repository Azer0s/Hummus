FROM golang:alpine as builder
RUN apk update && apk add --no-cache git && apk add --update alpine-sdk
WORKDIR $GOPATH/src/github.com/Azer0s/Hummus/
COPY . .
RUN rm go.mod
RUN chmod +x scripts/prepare_release.sh
RUN ./scripts/prepare_release.sh
RUN go build -o bin/hummus
RUN cp -r bin/* /go/bin/

FROM alpine:latest
COPY --from=builder /go/bin/ .
ENTRYPOINT ["/hummus"]