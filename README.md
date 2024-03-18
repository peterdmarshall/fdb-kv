# fdb-kv
A simple HTTP key-value store layer on top of FoundationDB.

Read more about [FoundationDB](https://www.foundationdb.org/)

Also check out [Huma](https://huma.rocks/), a new HTTP API framework for Go.

## Usage

Run the server, port defaults to `8888` if not provided.
```go run main.go [--port <port>]```

To set a value for a key:

```GET /{key}```
Example:
```
curl --request GET \
  --url http://localhost:8888/foo \
  --header 'Accept: application/json'
```

To retrieve the value for a key:

```
PUT /{key}
{
    "value": <value>
}
```
Example:
```
curl --request PUT \
  --url http://localhost:8888/foo \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json' \
  --data '{
  "value": "bar"
}'
```
