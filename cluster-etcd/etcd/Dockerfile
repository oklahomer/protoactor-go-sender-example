FROM alpine:latest

ENV ETCD_VERSION v3.2.32

# Install etcd
RUN wget https://github.com/etcd-io/etcd/releases/download/${ETCD_VERSION}/etcd-${ETCD_VERSION}-linux-amd64.tar.gz && \
    tar -xzvf etcd-${ETCD_VERSION}-linux-amd64.tar.gz -C /usr/local/bin && \
    ln -s /usr/local/bin/etcd-${ETCD_VERSION}-linux-amd64/etcd /usr/local/bin/

EXPOSE 2379

CMD ["/usr/local/bin/etcd", \
    "--name", "node1", \
    "--data-dir", "/etcd-data", \
    "--listen-client-urls", "http://0.0.0.0:2379", \
    "--listen-peer-urls", "http://0.0.0.0:2380", \
    "--advertise-client-urls", "http://127.0.0.1:2379", \
    "--initial-advertise-peer-urls", "http://127.0.0.1:2380", \
    "--initial-cluster", "node1=http://127.0.0.1:2380", \
    "--auto-compaction-retention", "1"]