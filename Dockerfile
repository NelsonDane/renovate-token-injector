FROM golang:1.27.1@sha256:137ca8442e368f5f5bb4d4f84d8cc1f6c2d898edf8a380459eb407f89d149b0d AS build
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o gh-app-token

FROM gcr.io/distroless/static@sha256:cd64bec9cec257044ce3a8dd3620cf83b387920100332f2b041f19c4d2febf93

WORKDIR /app
COPY --from=build /app/gh-app-token /app/gh-app-token

ENTRYPOINT ["/app/gh-app-token"]
