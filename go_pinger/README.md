## Pinger

Collect system stats by 'ping'ing the REST API, and provide simple monitoring UI.

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
git clone https://github.com/Sotera/go_watchman
# create docker-compose.yml
# create docker-compose.override.yml if needed
# create conf.toml
docker-compose up -d
```