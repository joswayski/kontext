# Kontext

A tool to automatically visualize the event flows in your Kafka clusters and give you **Kontext** on what consumers actually **do** with the data (in other words, your business logic). 

Nobody likes keeping diagrams or Markdown in sync anyway.

⚠️ *This project is in the early stages of development* ⚠️


## Getting Started

```bash
cp .env.example .env

docker compose up -d --build
```

This will start the following services:
- Kontext Web app at http://localhost:3000
- Kontext API at http://localhost:4000
- Kafka Clusters & Topics to simulate a ridesharing application called **Glide**:

| Cluster       | Broker(s)     |  Admin API    | Console UI    | Topics |
| ------------- | ------------- | ------------- | ------------- | -------------|
| Users  | kafka-users-0:9092  | http://localhost:29644  | http://localhost:8081  | user.created, user.updated|
| Drivers  | kafka-drivers-0:9092  | http://localhost:39644  | http://localhost:8082  |      driver.onboarded, driver.activated, driver.deactivated, driver.location.updated, driver.rating.submitted  | 
| Rides  | kafka-rides-0:9092  | http://localhost:19644  | http://localhost:8080  | ride.requested, ride.fare.calculated, ride.matched, ride.started, ride.completed, ride.cancelled  | 
| Payments  | kafka-payments-0:9092  | http://localhost:49644  | http://localhost:8083  | payment.method.added, payment.method.removed, payment.initiated, payment.succeeded, payment.failed, refund.issued  | 


### Notes
- If running outside of Docker, make sure to update the URLs in your `.env` to point to `localhost:PORT` instead. See [docker-compose.yaml](docker-compose.yaml) for more info.
- The Admin API and console will be removed eventually as we're trying to recreate them *somewhat*.
- We are also using [Redpanda]("https://redpanda.com/") due to the smaller resource footprint. 
