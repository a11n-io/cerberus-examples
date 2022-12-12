# backend build stage
FROM golang:alpine as backend-builder
RUN apk --no-cache add build-base
ENV GO111MODULE=on
WORKDIR /app
#COPY go.mod .
#COPY go.sum .
COPY ./ /app
RUN go mod download
RUN CGO_ENABLED=1 go build cmd/api/api.go

# final stage
FROM alpine
COPY --from=backend-builder /app/api /app/
COPY --from=backend-builder /app/.env.dev /app/
COPY --from=backend-builder /app/migrations /app/migrations
RUN apk add --no-cache bash busybox-extras
ENV GIN_MODE=release
WORKDIR /app

ENTRYPOINT ["./api"]