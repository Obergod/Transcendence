FROM nginx

RUN apt update
RUN apt install -y nginx openssl

RUN mkdir -p /etc/nginx/ssl

COPY --chmod=0755 ./conf/nginx.conf /etc/nginx/nginx.conf

RUN echo "[req]" > /tmp/openssl.cnf && \
    echo "distinguished_name=req" >> /tmp/openssl.cnf && \
    echo "[san]" >> /tmp/openssl.cnf && \
    echo "subjectAltName=DNS:localhost,DNS:gbakulin.42.fr" >> /tmp/openssl.cnf && \
    openssl req -x509 -nodes -newkey rsa:4096 \
      -days 365 \
      -keyout /etc/ssl/private/inception.key \
      -out /etc/ssl/certs/inception.crt \
      -subj "/C=FR/ST=IDF/L=Paris/O=42/OU=Student/CN=gbakulin.42.fr" \
      -extensions san -config /tmp/openssl.cnf

RUN mkdir -p /var/run/nginx

RUN chmod 755 /var/www/html
RUN chown -R www-data:www-data /var/www/html


COPY ./tools/nginx_script.sh nginx_script.sh

ENTRYPOINT ["bash", "nginx_script.sh"]