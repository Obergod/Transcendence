FROM postgres-bookworm

RUN apt-get update && \
    apt-get install -y postgres:14.22 && \
    rm -rf /var/lib/apt/lists/*

# COPY --chmod=0755 /conf/db.conf /etc/mysql/mariadb.conf.d/50-server.cnf
COPY ./tools/script.sh /script.sh

ENTRYPOINT ["bash", "script.sh"]