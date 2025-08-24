# Kontext

A tool to automatically visualize the event flows in your Kafka clusters and give you **Kontext** on what consumers actually **do** with the data (in other words, your business logic). 



Nobody likes keeping diagrams or Markdown in sync anyway.

⚠️ *This project is in the early stages of development* ⚠️


## Getting Started

#### Prerequisites
- Rust + Cargo
- Node 22
- Docker


#### Setup datasources
```bash
cp .env.example .env

docker compose up -d --build
```

This will create...
1. Kafka Clusters & Topics:

| Cluster       | Broker(s)     | Host Connection | Topics |
| ------------- | ------------- | ------------- | -------------|
| Rides  | rides:9092  | localhost:39092  | ride.requested, ride.fare.calculated, ride.matched, ride.started, ride.completed, ride.cancelled  | 
| Users  | users:9093  | localhost:39093  | user.created, user.updated|
| Drivers  | drivers:9094  | localhost:39094  |      driver.onboarded, driver.activated, driver.deactivated, driver.location.updated, driver.rating.submitted  | 
| Payments  | payments:9095  | localhost:39095  | payment.method.added, payment.method.removed, payment.initiated, payment.succeeded, payment.failed, refund.issued  | 
2. A fake ride ridesharing application called **Glide** which runs in the background producing and consuming messages from the topics above

#### Web & API
For simplicity, the **web** app and the **api** are run outside of Docker. 
To start both: `npm run dev` or individually `npm run web|api` on ports 3000/4000




