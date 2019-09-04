FROM golang:alpine as build

# Copy in the local repository to build from.
COPY / /go/src/gitlab.com/braneproject/branehub/

# Force Go to use the cgo based DNS resolver. This is required to ensure DNS
# queries required to connect to linked containers succeed.
ENV GODEBUG netdns=cgo

# Install dependencies and install/build branehub.
RUN apk add --no-cache --update git ca-certificates \
&&  go get github.com/gorilla/mux \
&&  go install /go/src/gitlab.com/braneproject/branehub

# Start a new, final image to reduce size.
FROM alpine as final

# Expose http port
EXPOSE 80

# Copy the binaries and TLS certificates from the build image.
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/branehub /bin/

ENTRYPOINT ["branehub"]
