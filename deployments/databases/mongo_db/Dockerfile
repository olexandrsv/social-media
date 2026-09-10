FROM mongo:7.0

ENV MONGO_INITDB_ROOT_USERNAME=mongo
CMD ["--replSet", "rs1", "--keyFile", "/run/secrets/mongodb-internal-secret","--bind_ip_all", "--port", "27017"]
