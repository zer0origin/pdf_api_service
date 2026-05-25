# API Backend Rework
Documentation can be created by using `swag init`.

# How to run locally?
When running the program, ensure the following environment variables have been set (I have included defaults):
* APP_HOST_NAME=localhost
* APP_PORT=8082
* DATA_SERVICE_URL=http://localhost:8083
* DATABASE_DB=postgres
* DATABASE_HOST=localhost
* DATABASE_PASSWORD=password
* DATABASE_PORT=5432
* DATABASE_USER=user

## Test Notes
- All tests must import flags.go inside testutil, a blank import is enough.