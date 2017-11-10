FROM alpine:latest
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN mkdir -p /opt/vice/
WORKDIR /opt/vice/
ADD vice-store /opt/vice/
RUN chmod +x /opt/vice/vice-store
ENV RETHINKDB_LOCATION=localhost \
    RETHINKDB_DATABASE=vice \
    RABBITMQ_LOCATION=localhost \
    RABBITMQ_USER=admin \
    RABBITMQ_PASS=admin \
    STORAGE_BASEPATH=/tmp/
CMD /opt/vice/vice-store \
    --rethinkdb-location $RETHINKDB_LOCATION \
    --rethinkdb-database $RETHINKDB_DATABASE \
    --rabbitmq-location $RABBITMQ_LOCATION \
    --rabbitmq-user $RABBITMQ_USER \
    --rabbitmq-pass $RABBITMQ_PASS \
    --storage-basepath $STORAGE_BASEPATH