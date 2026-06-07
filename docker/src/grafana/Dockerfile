# FROM debian:stable-slim as openssl-installer

# RUN apt-get update && \
#     apt-get install -y --no-install-recommends curl ca-certificates && \
#     apt-get clean && \
#     rm -rf /var/lib/apt/lists/*

# RUN apt-get update && \
#     apt-get install -y openssl && \
#     apt-get clean && \
#     rm -rf /var/lib/apt/lists/*

# FROM grafana/grafana:latest

# # Copy OpenSSL binary from the builder stage
# RUN mkdir -p /etc/grafana/ssl && \
#     chown -R grafana:grafana /etc/grafana/ssl && \
#     chmod 755 /etc/grafana/ssl

# # Copy the certificate generation script
# COPY generate-ssl.sh /usr/local/bin/generate-ssl.sh
# RUN chmod +x /usr/local/bin/generate-ssl.sh && \
#     chown grafana:grafana /usr/local/bin/generate-ssl.sh

# # Copy your custom grafana.ini
# COPY grafana.ini /etc/grafana/grafana.ini

# COPY --chown=472:472 grafana.key /etc/grafana/ssl/grafana.key || true
# COPY --chown=472:472 grafana.crt /etc/grafana/ssl/grafana.crt || true

# USER 472

FROM grafana/grafana:latest

USER root

RUN mkdir -p /etc/grafana/ssl \
    && chown -R 472:472 /etc/grafana/ssl

USER 472