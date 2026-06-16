FROM docker.elastic.co/logstash/logstash:7.15.0

RUN rm -f /usr/share/logstash/pipeline/logstash.conf

COPY docker/src/logstash/config/pipeline/logstash.conf \
     /usr/share/logstash/pipeline/logstash.conf
