FROM debian:stable-slim as openssl-installer

RUN apt-get update && \
    apt-get install -y --no-install-recommends curl ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN apt-get update && \
    apt-get install -y openssl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

FROM docker.elastic.co/kibana/kibana:9.4.0

COPY --from=openssl-installer /usr/bin/openssl /usr/bin/openssl

# COPY --from=builder /app/ssl/localhost+2-key.pem /etc/kibana/ssl/localhost-key.pem
# COPY --from=builder /app/ssl/localhost+2.pem /etc/kibana/ssl/localhost.crt

RUN chown -R 1000:1000 /etc/kibana/ssl && \
    chmod 644 /etc/kibana/ssl/localhost-key.pem && \
    chmod 644 /etc/kibana/ssl/localhost.crt

USER 1000:1000 