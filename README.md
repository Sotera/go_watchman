# go_watchman

[![Build Status](https://travis-ci.org/Sotera/go_watchman.svg?branch=master)](https://travis-ci.org/Sotera/go_watchman)

Watchman processes for which go is specifically well suited

## Docker build

```
cd deploy
# use go exec as entrypoint
./build.sh annotations 99 standalone [push]
# use supervisord as entrypoint
./build.sh annotations 99 supervisord [push]


# on remote host
# update docker-compose.yml
docker-compose up -d
```

## Deployment TODO

* Use env vars in docker-compose.yml to configure your app.
* If using standalone mode, add cli options in docker run or compose file.
* Mount supervisord.conf in docker-compose for custom supervisor conf.