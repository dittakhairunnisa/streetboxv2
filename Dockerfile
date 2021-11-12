FROM golang:alpine as builder

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR /opt/streetbox/

#use modules
COPY go.mod .

ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

COPY . .
COPY config.yml.example config.yml

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
      -ldflags='-w -s -extldflags "-static"' -a \
      -o /go/bin/streetbox .

############################
# STEP 2 build a small image
############################
FROM alpine:latest

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

RUN mkdir -p /opt/streetbox
RUN mkdir -p /opt/streetbox/image
RUN mkdir -p /opt/streetbox/doc
# Copy our static executable
COPY --from=builder /go/bin/streetbox .
COPY --from=builder /opt/streetbox/config.yml .
COPY --from=builder /opt/streetbox/credentials/streetbox-private-key.json .
COPY --from=builder /opt/streetbox/assets/menutemplate.csv .
# Use an unprivileged user.
# Cannot use this command because this project has ssl certs from streetbox certs folder (permission denied)
#USER appuser:appuser 
ENV GIN_MODE=release
ENV TZ=Asia/Jakarta
# Run the api-gateway binary.
ENTRYPOINT ["/streetbox"]
