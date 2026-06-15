FROM grafana/grafana:latest

USER root

RUN mkdir -p /etc/grafana/ssl \
    && chown -R 472:472 /etc/grafana/ssl

USER 472