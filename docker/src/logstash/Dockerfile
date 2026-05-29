FROM docker.elastic.co/logstash/logstash:9.4.0

RUN rm -f /usr/share/logstash/pipeline/logstash.conf

COPY docker/src/logstash/config/pipeline/logstash.conf \
     /usr/share/logstash/pipeline/logstash.conf


# FROM docker.elastic.co/logstash/logstash:9.4.0

# COPY docker/src/logstash/setup.sh /setup.sh
# COPY docker/src/logstash/logstash.yml /usr/share/logstash/config/logstash.yml

# RUN /setup.sh

# ENTRYPOINT ["/setup.sh"]