DO
$do$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'fiber-starter') THEN
      CREATE DATABASE "fiber-starter";
   END IF;

   IF NOT EXISTS (SELECT FROM pg_user WHERE usename = 'user') THEN
      CREATE USER "user" WITH PASSWORD 'password';
   END IF;

   GRANT ALL PRIVILEGES ON DATABASE "fiber-starter" TO "user";
END
$do$;