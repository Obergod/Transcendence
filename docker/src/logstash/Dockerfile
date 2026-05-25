FROM docker.elastic.co/logstash/logstash:9.4.0

RUN rm -f /usr/share/logstash/pipeline/logstash.conf

COPY docker/src/logstash/logstash.conf \
     /usr/share/logstash/pipeline/logstash.conf