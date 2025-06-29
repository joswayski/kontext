# Kontext

I want to build a tool that visualizes the event flows in my Kafka clusters and gives me **Kontext** on what the consumers *DO* with the data (business logic). 

This project is in the early stages of development, contributions are welcome!


## Getting Started

```bash
cp .env.example .env

docker compose up -d --build
```

This will start the following services:
- Web app at http://localhost:3000
- API at http://localhost:4000
- Kafka Production Cluster 
    - Broker URL at kafka-production-0:9092 and localhost:19092
    - Admin API URL at http://localhost:19644 
    - Console URL at http://localhost:8080
- Kafka Analytics Cluster 
    - Broker URL at kafka-analytics-0:9092 and localhost:29092
    - Admin API URL at http://localhost:29644
    - Console URL at http://localhost:8081

The Admin API and console will be removed eventually, as we're trying to recreate this somewhat. We are also using [Redpanda]("https://redpanda.com/") due to the smaller resource footprint. See [docker-compose.yaml](./docker-compose.yaml) for more details. 
