FROM debian:latest
MAINTAINER Louis Fradin <louis.fradin@gmail.com>

# Ports
EXPOSE 8080

# Volumes
VOLUME /var/log
VOLUME /telerdd

# Workdir
WORKDIR /telerdd

# Command
CMD ./bin/telerdd

# Data
COPY ./ ./
