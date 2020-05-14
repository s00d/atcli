FROM debian:stretch-slim

RUN DEBIAN_FRONTEND=noninteractive apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
  curl \
  wget \
  git \
  zip \
  jq


COPY *.sh /
ENTRYPOINT ["/entrypoint.sh"]

LABEL maintainer="s00d <virus191288@gmail.com>"