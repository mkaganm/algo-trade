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

Projeyi başlatmak için docker-compose komutlarını veya 
tanımlanan make komutlarını kullanabilirsiniz.

---
## DESCRIPTION

![FLOWCHART](https://raw.githubusercontent.com/mkaganm/algo-trade/refs/heads/master/documents/1.png)

Algo trade projesini 3 ana modülden oluşmaktadır. Bu modüller sırasıyla;
- collector
- proccessor
- trader

### Collector
Collector binance üzerinden anlık olarak web socket ile btc/usdt verileri toplar.
Topladığı bu verileri mongo üzerine yazar.

### Processor
Processor mongo db üzerindeki verileri alır ve işler.
Bu verilerin sma değerlerini hesaplar. Buna göre BUY veya SELL kararı verir.
Verdiği kararı mongo db üzerinde loglar.
Verdiği kararı redis ile trader modülüne streams ile sinyal olarak gönderir.

### Trader
Trader modülü redis üzerinden streams ile sinyalleri alır.
İşlediği sinyaller doğrultusunda redis üzerindeki stream mesajlarını acknowledge eder
Aldınan sinyaller doğrultusunda buraya alma ve satma komutları entegre edilebilir.


Bütün servisler pyroscope üzerinden metriklerinin prometheusa gönderir.
Prometheus üzerindeki metrikler grafana ile görselleştirilebilir.

---
### HEALTH CHECKER

3 servis içinde health check endpointi bulunmaktadır.
- collector: http://localhost:8080/healthcheck
- processor: http://localhost:8081/healthcheck
- trader: http://localhost:8082/healthcheck

Collector health check endpointi üzerinden mongo db bağlantısını kontrol eder.
Processor health check endpointi üzerinden mongo db ve redis bağlantısını kontrol eder.
Trader health check endpointi üzerinden redis bağlantısını kontrol eder.

HEALCHECK PHOTO

---

### METRICS

Metrikleri toplamak için pyroscope kullanıldı.
Pyroscope toplanan metrikleri prometheus ile entegre eder.
Grafana ile görselleştirilebilir.

GRAFANA PHOTO

---

## TODO List
- [ ] Add log system
- [ ] Add recover for all Go routines
- [ ] Add test ()
- [ ] Add more more comments