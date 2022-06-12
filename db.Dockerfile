FROM mysql:8.0.29

COPY ./database/*.sql /docker-entrypoint-initdb.d/