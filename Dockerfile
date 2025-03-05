FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN \
  apt-get -y update && \
  apt-get -y --no-install-recommends install \
  ca-certificates \
  curl \
  python3 \
  python3-pip \
  python3-venv \
  golang && \
  rm -rf /var/lib/apt/lists/*

RUN <<EOF
  mkdir -p /usr/share/ca-certificates/dn42
  echo "dn42/ca-xuu.crt" >> /etc/ca-certificates.conf
EOF

COPY ca-xuu.crt /usr/share/ca-certificates/dn42

RUN update-ca-certificates
