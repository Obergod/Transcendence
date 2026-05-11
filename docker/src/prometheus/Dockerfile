FROM prometheus

RUN apt-get update && \
    apt-get install -y \
        ca-certificates \
        curl \
        git \
        golang-go \
        build-essential \
        wget \
        tar \
        && apt-get clean && rm -rf /var/lib/apt/lists/*

WORKDIR /opt/monitoring
RUN mkdir node_exporter prometheus

RUN wget https://github.com/prometheus/node_exporter/releases/download/v1.10.2/node_exporter-1.10.2.linux-amd64.tar.gz && \
    tar xvf node_exporter-1.10.2.linux-amd64.tar.gz && \
    cp node_exporter-1.10.2.linux-amd64/node_exporter /usr/local/bin/node_exporter && \
    chmod +x /usr/local/bin/node_exporter

RUN wget https://github.com/prometheus/prometheus/releases/download/v3.7.3/prometheus-3.7.3.linux-amd64.tar.gz && \
    tar xvf prometheus-3.7.3.linux-amd64.tar.gz && \
    cp prometheus-3.7.3.linux-amd64/prometheus /usr/local/bin/prometheus && \
    chmod +x /usr/local/bin/prometheus && \
    rm -rf prometheus-3.7.3.linux-amd64*

COPY ./config/prometheus.yml /etc/prometheus.yml

EXPOSE 9090

CMD /usr/local/bin/node_exporter & /usr/local/bin/prometheus \
    --config.file=/etc/prometheus.yml \
    --web.listen-address=:9090
