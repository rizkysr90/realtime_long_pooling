Prerequisites :

- PostgreSQL
- Run migration database :

- Install dbmate
- dbmate --url "postgresql://postgres:yourpassword@localhost:4444/postgres?sslmode=disable" -d "resources/pgsql/migrations" up
