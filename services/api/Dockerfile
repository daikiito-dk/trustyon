FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /trustyon-api ./cmd/server

FROM alpine:3.22
RUN adduser -D app
USER app
COPY --from=build /trustyon-api /trustyon-api
EXPOSE 8080
ENTRYPOINT ["/trustyon-api"]
