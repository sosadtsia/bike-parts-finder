# Bike Parts Finder

A cloud-native application for finding bike parts by brand, model, and year.

## Project Overview

This application helps cyclists find compatible parts for their bicycles by searching across multiple online retailers. It uses web scraping to gather real-time parts data and presents structured results with images.

## Architecture

```
User
|
React.js Frontend <--> Go API Backend (mTLS)
|
+-----------+------------+
|            |
PostgreSQL    Redis Cache
^            |
|            |
Kafka Consumer <--- Kafka Topic: scrape_results
^
|
Kafka on EKS
^
|
Scraper (Go) <--- Kafka Topic: scrape_requests
|
External websitesß
```

## Components

- **Frontend**: Go WebAssembly for searching and viewing bike parts
- **Backend API**: Go service handling API logic and database/cache access
- **Security**: mTLS for secure client-server communication with certificate management
- **Scraper**: Go service that scrapes websites for bike parts data
- **Database**: PostgreSQL for storing structured parts data
- **Cache**: Redis for caching frequent queries
- **Message Broker**: Kafka for handling asynchronous scraping jobs

## Infrastructure

- Hosted on Amazon EKS
- Infrastructure provisioned with OpenTofu
- Deployed via GitOps using ArgoCD
- Certificate management using cert-manager for mTLS
- Monitoring with Prometheus and Grafana
- Backups and disaster recovery with Velero

## Development

See [DEVELOPMENT.md](./docs/DEVELOPMENT.md) for local development instructions.

## Deployment

See [DEPLOYMENT.md](./docs/DEPLOYMENT.md) for deployment instructions.
