FROM debian

ENV REDIS_PASS=avito
ENV REDIS_TLS_PORT=6379


RUN apt-get update
RUN apt-get install 
RUN apt install -y lsb-release curl gpg

RUN curl -fsSL https://packages.redis.io/gpg | gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg

RUN apt install -y redis

COPY ./tls /etc/redis/tls

# команда для запуска редиса через tls
CMD redis-server --requirepass ${REDIS_PASS} --tls-port ${REDIS_TLS_PORT} --port 0 \
	--tls-cert-file /etc/redis/tls/redis.crt \
	--tls-key-file /etc/redis/tls/redis.key \
	--tls-ca-cert-file /etc/redis/tls/ca.crt




