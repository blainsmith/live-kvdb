# Live Key/Value Database

## Database Requirements

1. Get(key) pulls a value from the database
2. Set(key, value) creates or updates a value from the database
3. Del(key) removes a value from the database
4. Keys() returns all the keys from the database
5. Atomic updates and avoid race conditions to ensure data integrity

## Network Connections

### Persistant TCP Socket

Should be able to connect with `telnet` to interact with the database

Commands should follow `COMMAND KEY VALUE` space separated pattern:

1. `get key` returns value to the connected client
2. `set key value`
3. `del key`
4. `keys` returns a list of keys
5. `quit`: closes the tcp connection

### HTTP REST API

Should be able to use Postman, `curl`, or any other HTTP-based program

Commands should map to HTTP methods:

1. `GET /key` returns value in the body response
2. `POST /key` with value in the request body
3. `DELETE /key`
4. `GET /` returns all keys in the body response