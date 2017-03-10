# go_watchman
Watchman processes for which go is specifically well suited

## Docker build

```
cd deploy
./build.sh annotations mytag [push]

# on remote host
# update docker-compose.yml
docker-compose up -d
```

## Deployment TODO

* Use env vars in docker-compose to config your app.
* Supervisor keeps ur app running. Do you want that?
* Mount supervisord.conf in docker-compose for custom supervisor conf.