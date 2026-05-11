FROM debian:bookworm

RUN apt-get update && \
    apt-get install -y mariadb-server && \
    rm -rf /var/lib/apt/lists/*

COPY --chmod=0755 /conf/db.conf /etc/mysql/mariadb.conf.d/50-server.cnf
COPY ./tools/mariadb.sh /mariadb.sh

ENTRYPOINT ["bash", "mariadb.sh"]