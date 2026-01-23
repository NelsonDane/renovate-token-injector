FROM golang:1.25.6 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o gh-app-token

FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=build /app/gh-app-token /app/gh-app-token
ENTRYPOINT ["/app/gh-app-token"]
