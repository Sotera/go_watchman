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