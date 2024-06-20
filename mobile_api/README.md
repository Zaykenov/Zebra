сервис генерации pdf

https://localhost:7147/swagger/index.html


#run before migration:

``` sql 
  CREATE SCHEMA mobile_api;
  CREATE ROLE mobile_api_user NOSUPERUSER NOCREATEDB NOCREATEROLE NOINHERIT LOGIN NOREPLICATION NOBYPASSRLS PASSWORD '123456';
  GRANT ALL ON SCHEMA mobile_api TO mobile_api_user;
```