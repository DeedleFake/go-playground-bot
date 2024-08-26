FROM golang:alpine AS build

WORKDIR /build
COPY . .

ENV CGO_ENABLED=0
RUN go build ./cmd/playbot

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /build/playbot /

ENTRYPOINT ["/playbot"]
