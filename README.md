# go_watchman

[![Build Status](https://travis-ci.org/Sotera/go_watchman.svg?branch=master)](https://travis-ci.org/Sotera/go_watchman)

Watchman processes for which go is specifically well suited

## Dependencies

* Go 1.6+
* Glide 0.12

## Develop

```
cd <sub-project>
glide install
```

## Docker build

```
cd deploy
# use go exec as entrypoint
./build.sh annotations_mock_server 99 annotations/mockery/cmd standalone [push]
# use supervisord as entrypoint
./build.sh annotations_mock_server 99 annotations/mockery/cmd supervisord [push]


# on remote host
# update docker-compose.yml
docker-compose up -d
```

## Deployment TODO

* Use env vars in docker-compose.yml to configure your app.
* If using standalone mode, add cli options in docker run or compose file.
* Mount supervisord.conf in docker-compose for custom supervisor conf.