# Kontext

A tool to automatically visualize the event flows in your Kafka clusters and give you **Kontext** on what consumers actually **do** with the data (in other words, your business logic). 

Nobody likes keeping diagrams or Markdown in sync anyway.

⚠️ *This project is in the early stages of development* ⚠️


## Getting Started

```bash
cp .env.example .env

docker compose up -d --build
```

This will start the following services to simulate a ridesharing application:
- Web app at http://localhost:3000
- API at http://localhost:4000
- Kafka Clusters:
  - **rides**: Broker at kafka-rides-0:9092, Admin API at http://localhost:19644, Console at http://localhost:8080
  - **users**: Broker at kafka-users-0:9092, Admin API at http://localhost:29644, Console at http://localhost:8081
  - **drivers**: Broker at kafka-drivers-0:9092, Admin API at http://localhost:39644, Console at http://localhost:8082


### Notes
- If running outside of Docker, make sure to update the URLs in your `.env` to point to `localhost:PORT` instead. See [docker-compose.yaml](docker-compose.yaml) for more info.
- The Admin API and console will be removed eventually as we're trying to recreate them *somewhat*.
- We are also using [Redpanda]("https://redpanda.com/") due to the smaller resource footprint. 
