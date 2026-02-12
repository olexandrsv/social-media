FROM postgres:16-alpine
ARG SERVICE

ENV POSTGRES_PASSWORD=database
ENV POSTGRES_USER=postgres
ENV POSTGRES_DB=media
COPY ${SERVICE}.sql /docker-entrypoint-initdb.d/