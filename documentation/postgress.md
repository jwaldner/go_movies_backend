- sudo service postgresql status
- sudo service postgresql start

CREATE USER joe WITH PASSWORD '';

GRANT ALL PRIVILEGES ON DATABASE "go_movies" to joe;

psql -d go_movies -f go_movies.sql
