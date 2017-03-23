## Docker Deploy

```
# cp docker-compose.yml docker-compose.override.yml
# modify override file as needed, then:
docker-compose up -d
curl -H "Content-Type: application/json" "http://localhost:8080/annotations/refid/1?from_date=2017-03-14&to_date=2017-03-17&type=label"
```