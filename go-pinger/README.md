## Pinger

### Usage

```
cd api
PORT=3000 DB_HOST=mongo:27017 ./api
cd request
cp conf.toml.example conf.toml
./request
```

### Deploy

```
# Make sure you've built the binaries and committed
git clone ...
cd go-pinger
# create docker-compose.override.yml if needed
# create ./request/conf.toml
docker-compose up -d
```