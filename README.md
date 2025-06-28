# Kontext

I want to build a tool that visualizes the event flows in my Kafka clusters and gives me **Kontext** on what the consumers *DO* with the data (business logic). 

This project is in the early stages of development, contributions are welcome!


## Getting Started

```bash
cp .env.example .env

docker compose up -d --build
```

This will start the following services:
- Web app on port 3000
- API on port 4000
- Kafka Production Cluster 
    - Broker URL: 19092
    - Schema Registry URL: 18081
    - Admin API URL: 19644
    - Console URL: 8080
- Kafka Analytics Cluster 
    - Broker URL: 29092
    - Schema Registry URL: 28081
    - Admin API URL: 29644
    - Console URL: 8081

See [docker-compose.yaml](./docker-compose.yaml) for more details. We are also using [Redpanda]("https://redpanda.com/") due to the smaller resource footprint.
