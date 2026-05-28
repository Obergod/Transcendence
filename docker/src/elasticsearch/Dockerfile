# # FROM docker.elastic.co/elasticsearch/elasticsearch:9.4.0

# # USER root

# # COPY entrypoint.sh /usr/local/bin/custom-entrypoint.sh
# # RUN chmod +x /usr/local/bin/custom-entrypoint.sh

# # USER elasticsearch

# # ENTRYPOINT ["/usr/local/bin/custom-entrypoint.sh"]

# FROM docker.elastic.co/elasticsearch/elasticsearch:9.4.0

# # Set environment variables
# ENV ELASTIC_PASSWORD=elastic123

# # Install dependencies
# USER root
# RUN apt-get update && apt-get install -y \
#     curl \
#     jq \
#     && rm -rf /var/lib/apt/lists/* \
#     && chown -R elasticsearch:elasticsearch /usr/share/elasticsearch

# # Copy your script
# COPY bootstrap-elasticsearch.sh /usr/share/elasticsearch/

# # Make it executable
# RUN chmod +x /usr/share/elasticsearch/bootstrap-elasticsearch.sh

# # Override the entrypoint correctly
# ENTRYPOINT ["sh", "-c", "/usr/share/elasticsearch/bootstrap-elasticsearch.sh"]