FROM golang:1.12-alpine as builder

# Copy in the local repository to build from.
COPY . /go/src/gitlab.com/braneproject/branehub

# Force Go to use the cgo based DNS resolver. This is required to ensure DNS
# queries required to connect to linked containers succeed.
ENV GODEBUG netdns=cgo

# Install dependencies and install/build branehub.
RUN apk add --no-cache --update git \
&& cd /go/src/gitlab.com/braneproject/branehub \
&&  go get github.com/gorilla/mux \
&&  go install gitlab.com/braneproject/branehub

# Start a new, final image to reduce size.
FROM alpine as final

# Expose lnd ports (server, rpc).
EXPOSE 80

# Copy the binaries and entrypoint from the builder image.
COPY --from=builder /go/bin/branehub /bin/

# Add bash.
RUN apk add --no-cache \
    bash

# Copy the entrypoint script.
COPY "start-branehub.sh" .
RUN chmod +x start-branehub.sh
