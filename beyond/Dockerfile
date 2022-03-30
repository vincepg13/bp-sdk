# Simple usage with a mounted data directory:
# > docker build -t beyond .
# > docker run -it -p 46657:46657 -p 46656:46656 -p 26656:26656 -p 26657:26657 -v ~/.beyondd:/root/.beyondd -v ~/.beyondcli:/root/.beyondcli beyond beyondd init
# > docker run -it -p 46657:46657 -p 46656:46656 -p 26656:26656 -p 26657:26657 -v ~/.beyondd:/root/.beyondd -v ~/.beyondcli:/root/.beyondcli beyond beyondd start
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES make git libc-dev bash gcc linux-headers eudev-dev

# Set working directory for the build
WORKDIR /go/src/github.com/vincepg13/bp-sdk/beyond

# Add source files
COPY . .

# Install minimum necessary dependencies, build Beyond SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
    make get_tools && \
    make get_vendor_deps && \
    make build

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/vincepg13/bp-sdk/beyond/bin/beyondd /usr/bin/beyondd
COPY --from=build-env /go/src/github.com/vincepg13/bp-sdk/beyond/bin/beyondcli /usr/bin/beyondcli

# Run beyondd by default, omit entrypoint to ease using container with beyondcli
CMD ["beyondd"]
