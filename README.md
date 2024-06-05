1. Run postgres container
   docker run -d --name todo \
    -e POSTGRES_USER=todo \
    -e POSTGRES_PASSWORD=password \
    -e POSTGRES_DB=todo \
    -p 5432:5432 \
    postgres:latest
