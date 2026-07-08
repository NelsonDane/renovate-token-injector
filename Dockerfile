FROM golang:1.26.5@sha256:63f132d58c1f589f0dcda584933a9bb44bfda1150f1506377f5a902f34d86033 AS build
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
