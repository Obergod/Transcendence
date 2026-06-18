# FROM debian:latest

# RUN apt-get update && apt-get install -y cron && rm -rf /var/lib/apt/lists/*

# COPY docker/src/backup/backup_script.sh /backup_script.sh

# RUN chmod +x /backup_script.sh

# RUN cron

# RUN echo "*/5 * * * * /backup_script.sh >>" > /etc/cron.d/backup_script

# CMD ["tail", "-f", "/dev/null"]

FROM debian:latest

RUN apt-get update && \
    apt-get install -y cron && \
    rm -rf /var/lib/apt/lists/*

COPY docker/src/backup/backup_script.sh /backup_script.sh
RUN chmod +x /backup_script.sh

RUN echo "*/5 * * * * root /backup_script.sh >> /var/log/backup.log 2>&1" \
    > /etc/cron.d/backup_script && \
    chmod 0644 /etc/cron.d/backup_script

CMD ["cron", "-f"]