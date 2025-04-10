# Kontext

**Automated Kafka event flow visualization and business logic mapping**

> ⚠️ **This project is in the very early stages of development** ⚠️


## Overview
Understanding Kafka event flows, their downstream impacts, and evolving schemas often requires significant effort. Traditional documentation methods like markdown or static diagrams quickly become outdated, and are tedious to maintain. Existing observability tools might offer metrics on throughput or lineage, but they typically lack insight into *what* your services *do* with the data.
**Kontext provides an always up-to-date, visual understanding of your event-driven architecture with a low operational cost, freeing your team from the chore of manual documentation.**


## Planned Features


- **Continuous Kafka Discovery**:  Automatically discovers and continuously updates your Kafka topology (producers, consumers, topics, schemas) across clusters.
- **Code-Aware Lineage**: Visualizes event flows, automatically linking Kafka topics to the specific code functions that handle them via code analysis.
- **Self-Hosted & Secure**: Deploy entirely within your infrastructure, maintaining full control over access configuration (Kafka, code repos, LLM keys) for security and privacy.
- **Live Message Sampling**: View live message examples from topics for concrete data context and debugging.



## Getting Started

### Prerequisites

- [Rust](https://www.rust-lang.org/tools/install)
- [Node.js](https://nodejs.org/)
- [Docker](https://docs.docker.com/get-docker/)

### Setup and Usage

```bash
# Get help for any script with --help!

# Builds shared library, API, installs web dependencies, starts Kafka and MySQL
./scripts/setup.sh

# Starts specified services (API, web) in development mode
./scripts/run.sh                # all services
./scripts/run.sh api web        # specific services

# Runs tests for the specified services
./scripts/test.sh               # all tests
./scripts/test.sh api           # only API tests

```

## Contributing

Ideas, feedback, and contributions are welcome! Feel free to open an issue to discuss anything and everything. (Please keep in mind the early stage of the project).


## License

Kontext is licensed under the [Apache 2.0 License](LICENSE)

## Troubleshooting

See [TROUBLESHOOTING.md](docs/TROUBLESHOOTING.md)
