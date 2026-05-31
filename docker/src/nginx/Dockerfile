FROM nginx:alpine

RUN apk add --no-cache openssl bash

# COPY docker/src/nginx/conf/generate_certs.sh /usr/local/bin/generate_certs.sh
# RUN chmod +x /usr/local/bin/generate_certs.sh

COPY docker/src/nginx/tools/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN mkdir -p /etc/nginx/ssl && \
    chmod 700 /etc/nginx/ssl

ENTRYPOINT ["/entrypoint.sh"]