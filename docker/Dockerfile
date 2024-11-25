# Base image for certificates
FROM alpine:3.6 AS alpine
RUN apk add -U --no-cache ca-certificates && update-ca-certificates

# Builder stage for Flyway installation
FROM alpine:3.6 AS builder

RUN apk add --no-cache bash

# Install necessary tools for downloading and extracting Flyway
RUN apk add --no-cache wget tar

# Define Flyway version and installation directory
ENV FLYWAY_VERSION=10.21.0
ENV FLYWAY_DIR=/opt/flyway-${FLYWAY_VERSION}

# Download and extract Flyway
RUN wget --no-check-certificate https://repo1.maven.org/maven2/org/flywaydb/flyway-commandline/${FLYWAY_VERSION}/flyway-commandline-${FLYWAY_VERSION}-linux-x64.tar.gz && \
    tar -zxf flyway-commandline-${FLYWAY_VERSION}-linux-x64.tar.gz && \
    mkdir -p /opt && \
    mv flyway-${FLYWAY_VERSION} ${FLYWAY_DIR}

# Final image stage
FROM alpine:3.6 AS final
ENV GODEBUG netdns=go

# Copy certificates from the alpine stage
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the Flyway installation from the builder stage
COPY --from=builder ${FLYWAY_DIR} ${FLYWAY_DIR}

# Add the Flyway binary to the PATH
ENV PATH="${FLYWAY_DIR}:${PATH}"

# Copy the plugin binary
ADD release/linux/amd64/plugin /bin/

# Set the entrypoint
ENTRYPOINT ["/bin/plugin"]
