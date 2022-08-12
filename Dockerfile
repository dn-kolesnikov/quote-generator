# STEP 1 build executable binary
FROM golang:alpine as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates && rm -rf /var/cache/apk/*

#ENV \
#    APP_USER=app \
#    APP_UID=1001

# See https://stackoverflow.com/a/55757473/12429735
#RUN \
#    adduser \
#    --disabled-password \
#    --gecos "" \
#    --home "/nonexistent" \
#    --shell "/sbin/nologin" \
#    --no-create-home \
#    --uid "$APP_UID" \
#    "$APP_USER"

WORKDIR ${GOPATH}/src/app/
COPY . .

# Fetch dependencies.
RUN \
    go mod download && \
    go mod verify && \
    go mod vendor
#    && \
#    go vet -v && \
#    go test -v

# Build the binary
RUN \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/app \
    ./cmd/quote-generator/main.go

# STEP 2 build a small image
FROM scratch

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy our static executable
COPY --from=builder /go/bin/app /app/app

# Use an unprivileged user.
USER nobody:nogroup

# Run the hello binary.
ENTRYPOINT ["/app/app"]
