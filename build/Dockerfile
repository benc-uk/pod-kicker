# ================================================================================================
# === Stage 1
# ================================================================================================
FROM golang:1.17-alpine as go-build
WORKDIR /build
ARG GO_PACKAGE="github.com/benc-uk/pod-kicker/cmd"
ARG VERSION="0.0.0"
ARG BUILD_INFO="Not set"

# Install system dependencies
RUN apk update && apk add git gcc musl-dev

# Fetch and cache Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy in Go source files
COPY cmd/ ./cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o pod-kicker \
  -ldflags "-X main.version=$VERSION -X 'main.buildInfo=$BUILD_INFO'" \
  $GO_PACKAGE

# ================================================================================================
# === Stage 2
# ================================================================================================
FROM scratch
WORKDIR /app 
COPY --from=go-build /build/pod-kicker . 
CMD [ "./pod-kicker"]