env "local" {
  url     = "postgres://todo:password@172.17.0.2:5432/todo?sslmode=disable"
  src     = "./schema.sql"
  dev     = "postgres://postgres:post@172.17.0.2:5432/postgres?sslmode=disable"
}
