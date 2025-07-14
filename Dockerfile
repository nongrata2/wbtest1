FROM golang:1.24 AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY migrations ./migrations
COPY static ./static

RUN CGO_ENABLED=0 go build -o /app ./cmd/app/main.go

FROM alpine:3.20

COPY --from=build /app /app
COPY --from=build /src/static ./static

ENTRYPOINT ["/app"]