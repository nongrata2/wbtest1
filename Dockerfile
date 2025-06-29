FROM golang:1.23 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build -o /app ./cmd/app/main.go

FROM alpine:3.20

COPY --from=build /app /app

ENTRYPOINT ["/app"]