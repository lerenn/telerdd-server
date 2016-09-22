FROM debian:latest
MAINTAINER Louis Fradin <louis.fradin@gmail.com>

# Ports
EXPOSE 8080

# Volumes
VOLUME /var/log
VOLUME /telerdd-server

# Workdir
WORKDIR /telerdd-server

# Command
CMD ./bin/telerdd-server

# Data
COPY ./ ./
