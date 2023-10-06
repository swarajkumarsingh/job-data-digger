FROM redis:5.0.7-alpine3.11
COPY redis.conf /usr/local/etc/redis/redis.conf
RUN mkdir -p /data
CMD [ "redis-server", "/usr/local/etc/redis/redis.conf" ]