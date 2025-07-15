# Kontext

A tool to automatically visualize the event flows in your Kafka clusters and give you **kontext** on what consumers actually **do** with the data (in other words, your business logic). 

Nobody likes keeping diagrams or Markdown in sync anyway.

⚠️ *This project is in the early stages of development* ⚠️


## Getting Started

### Prerequisites
- Go >= 1.24
- Node 22
- Docker


### Setup datasources
```bash
cp .env.example .env

docker compose up -d --build
```

This will create...
- Kafka Clusters & Topics:

| Cluster       | Broker(s)     |  Admin API    | Console UI    | Topics |
| ------------- | ------------- | ------------- | ------------- | -------------|
| Rides  | kafka-rides-0:9092  | http://localhost:19644  | http://localhost:8080  | ride.requested, ride.fare.calculated, ride.matched, ride.started, ride.completed, ride.cancelled  | 
| Users  | kafka-users-0:9092  | http://localhost:29644  | http://localhost:8081  | user.created, user.updated|
| Drivers  | kafka-drivers-0:9092  | http://localhost:39644  | http://localhost:8082  |      driver.onboarded, driver.activated, driver.deactivated, driver.location.updated, driver.rating.submitted  | 
| Payments  | kafka-payments-0:9092  | http://localhost:49644  | http://localhost:8083  | payment.method.added, payment.method.removed, payment.initiated, payment.succeeded, payment.failed, refund.issued  | 
- A fake ride ridesharing application called **Glide** which runs in the background producing and consuming(WIP) messages from the topics above


For simplicity, the **web** app and the **api** are run outside of Docker

- To start the web app: `cd web` && `npm run dev` at http://localhost:3000

- To start the API: `cd api` && `go run .` at at http://localhost:3001


### Notes
- If running *inside* of Docker, make sure to update the URLs in your `.env` to point to `kafka-$CLUSTER-0:PORT` instead. See [docker-compose.yaml](docker-compose.yaml) for more info.
- The Admin API and console will be removed eventually as we're trying to recreate them *somewhat*.
- We are also using [Redpanda]("https://redpanda.com/") due to the smaller resource footprint. 
