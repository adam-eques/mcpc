# syntax=docker/dockerfile:1

# ---- build stage ----
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
ARG VERSION=docker
RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w -X github.com/adam-eques/mcpc/internal/version.Version=${VERSION}" \
    -o /out/mcpc ./cmd/mcpc

# ---- runtime stage ----
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/mcpc /usr/local/bin/mcpc
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/mcpc"]
