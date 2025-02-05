FROM golang:1.23.3-alpine3.20 AS build

RUN apk --no-cache add make git

COPY . /go/src
WORKDIR /go/src

RUN make build

FROM alpine:3.20

COPY --from=build /go/src/bin/finapi-server /usr/local/bin/finapi-server

EXPOSE 8080

ENTRYPOINT ["finapi-server"]