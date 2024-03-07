# advanced-container-counter
An upgrade version of container counter.

## 1. Go to demo branch
Please use [demo](https://github.com/Trip1eLift/container-counter/tree/demo) for demo run where it removes some development tool from docker.
```text
git checkout demo
```

## 2. Run services with docker, do:
```text
make start
```

## 3. Client Mock
Mimic client traffics to continue to send traffic to cluster.

1. Interval is default to 3
2. mode can go from 0 to 10 (0-4 for now)

### 2 active containers with client traffic interval of 3
```text
curl --header "Content-Type: application/json" --request POST --data '{"mode":"2","interval":"3"}' http://localhost:7001/toggle
```

### 4 active containers with client traffic interval of 2
```text
curl --header "Content-Type: application/json" --request POST --data '{"mode":"4","interval":"2"}' http://localhost:7001/toggle
```

### Visit [toggle.rest](toggle.rest) for all toggle mode suggestion

The current setting only spawn 4 cluster containers due to computation resource congestion.

Use rest client vscode plugin to run http requests in [toggle.rest](toggle.rest).

## TODO:
1. Nothing currently...