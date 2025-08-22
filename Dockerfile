# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.25 AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
COPY cmd cmd
COPY pkg pkg
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o openapi2go

FROM scratch
COPY --from=build /app/openapi2go /usr/bin/

ENTRYPOINT ["/usr/bin/openapi2go"]
