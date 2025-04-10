# ALGO TRADE

---

### SETUP

```
make build
```
```
make run
```

or

```
docker-compose build
```
```
docker-compose up -d
```

You can start the project using the docker-compose commands or the defined make commands.

---
## DESCRIPTION

![FLOWCHART](https://raw.githubusercontent.com/mkaganm/algo-trade/refs/heads/master/documents/flowchart.png)

The Algo Trade project consists of 3 main modules. These are, respectively:
- collector
- processor
- trader

### Collector
The Collector collects real-time BTC/USDT data from Binance via WebSocket.
It writes the collected data to MongoDB.

### Processor
The Processor retrieves and processes data from MongoDB. 
It calculates the SMA values of this data and makes BUY or SELL decisions accordingly. 
The decision is logged in MongoDB. 
The decision is sent as a signal to the Trader module via Redis streams.

### Trader
The Trader module receives signals via Redis streams. 
It acknowledges the stream messages in Redis based on the processed signals.
Buy and sell commands can be integrated here based on the received signals.

---
All services have health check endpoints.

All services send their metrics to Prometheus via Pyroscope. 
The metrics in Prometheus can be visualized using Grafana.

---
### HEALTH CHECKER

There are health check endpoints for 3 services.
- collector:
```
curl -X GET http://localhost:8080/healthcheck
```
- processor:
```
curl -X GET http://localhost:8082/healthcheck
```
- trader:
```
curl -X GET http://localhost:8083/healthcheck
```

The Collector health check endpoint checks the MongoDB connection. 

The Processor health check endpoint checks the MongoDB and Redis connections. 

The Trader health check endpoint checks the Redis connection.

![](https://raw.githubusercontent.com/mkaganm/algo-trade/refs/heads/master/documents/healthcheck.png)

---

### METRICS

Pyroscope is used to collect metrics.
Pyroscope integrates the collected metrics with Prometheus.
They can be visualized using Grafana.

Example metric queries are provided below.
Many more metrics are already being collected through Pyroscope and Prometheus,
and they can also be visualized via Grafana.

Collector cpu 
```
sum(rate(node_cpu_seconds_total{job="collector-metrics", mode!="idle"}[1m]))
```
Collector memory 
```
sum(node_memory_MemTotal_bytes{job="processor-metrics"} - node_memory_MemAvailable_bytes{job="processor-metrics"})
```
Processor cpu 
```
sum(rate(node_cpu_seconds_total{job="processor-metrics", mode!="idle"}[1m]))
```
Processor memory 
```
sum(node_memory_MemTotal_bytes{job="processor-metrics"} - node_memory_MemAvailable_bytes{job="processor-metrics"})
```
Trader cpu 
```
sum(rate(node_cpu_seconds_total{job="processor-metrics", mode!="idle"}[1m]))
```
Trader memory
```
sum(node_memory_MemTotal_bytes{job="trader-metrics"} - node_memory_MemAvailable_bytes{job="trader-metrics"})
```

![](https://raw.githubusercontent.com/mkaganm/algo-trade/refs/heads/master/documents/grafana.png)
---

### MongoDB

```
http://127.0.0.1:8081/db/btc_data/
```

There are 2 collections in MongoDB.
- trade_signals: Logs of processed signals
  ![](https://raw.githubusercontent.com/mkaganm/algo-trade/refs/heads/master/documents/processlogs.png)
---
- depth: Collected price data
  ![](https://raw.githubusercontent.com/mkaganm/algo-trade/refs/heads/master/documents/btcdatadb.png)
---

### Redis

Signals are being sent and received via trade_signals_stream in Redis.

```
http://127.0.0.1:8001/redis-stack/browser
```


![](https://raw.githubusercontent.com/mkaganm/algo-trade/refs/heads/master/documents/redis.png)

---
  

### Code Quality
golangci-lint was used for code quality, 
and the project's code was written according to these rules.

---
# TECHNOLOGIES

 - Golang
 - MongoDB
 - Redis
 - Pyroscope
 - Prometheus
 - Grafana
 - Docker
 - golinter

---
### TEST 

Example tests have been written for SMA calculations and signal processing.

![](https://raw.githubusercontent.com/mkaganm/algo-trade/refs/heads/master/documents/unittest.png)
---


### TODO
2. A short document that explains:
   ○ How you would ensure scalability, fault tolerance, and security.
   ○ Challenges faced and how you solved them.
