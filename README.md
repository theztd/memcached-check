# Memcached check 

This utility generates prometheus metrics after each run. It work's well with a promethesu node_exporter with configured metrics path.


```bash
Usage of memcached-check:
  -addr string
        Memcached server address. (default "127.0.0.1:11211")
  -metrics string
        Define path where the metrics have to be saved (default "./metrics_test.prom")
  -service string
        Service tag variable (usefull to identify which memcached we measure). (default "example")
  -write
        Do a write and read test
```