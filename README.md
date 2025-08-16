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

# Create JMX auth files in ./jmx (required for metrics with auth enabled; matches defaults in .env)
mkdir jmx
echo "admin secret" > jmx/jmx_password
echo "admin readwrite" > jmx/jmx_access
chmod 400 jmx/jmx_password jmx/jmx_access

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
- A fake ride ridesharing application called **Glide** which runs in the background producing and consuming messages from the topics above


For simplicity, the **web** app and the **api** are run outside of Docker

- To start the web app: `cd web` && `npm run dev` at http://localhost:3000
- To start the API: `cd api` && `go run .` at http://localhost:3001


### Notes
- If running *inside* of Docker, make sure to update the URLs in your `.env` to point to the internal container names (e.g., kafka-rides:9092) instead. See [docker-compose.yaml](docker-compose.yaml) for more info.
- The clusters are running Apache Kafka (KRaft mode) with JMX enabled for metrics extraction. Topics are auto-created on first produce (via KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE=true) for demo simplicity; in prod, manage topics manually.
- To extract metrics, ensure your .env includes JMX URLs for each cluster, e.g.:
  ```
  KAFKA_RIDES_JMX_URL=service:jmx:rmi:///jndi/rmi://localhost:19644/jmxrmi
  # Repeat for other clusters. In prod, this would point to your broker's JMX endpoint (enable JMX on brokers if not already).
  # Optional: KAFKA_RIDES_JMX_USERNAME=admin, KAFKA_RIDES_JMX_PASSWORD=secret (if auth is enabled).
  ```
- Metrics are queried via JMX in the API code—see pkg/kafka/producers.go for details. This works across Apache Kafka, Confluent, AWS MSK, etc., as long as JMX is exposed.
- JMX auth files (jmx/jmx_password and jmx/jmx_access) are mounted for security. If you don't need auth, edit docker-compose.yaml to disable it in KAFKA_JMX_OPTS and remove the volumes.
