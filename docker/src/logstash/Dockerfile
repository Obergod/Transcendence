FROM docker.elastic.co/logstash/logstash:7.15.0

USER root

COPY docker/src/logstash/config/logstash.yml \
     /usr/share/logstash/config/logstash.yml

RUN chown logstash:root /usr/share/logstash/config/logstash.yml && \
    chmod 640 /usr/share/logstash/config/logstash.yml

USER logstash