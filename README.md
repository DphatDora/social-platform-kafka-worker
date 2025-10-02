# social-platform-microservice

This repository contains a Kafka worker service for processing tasks in a social platform application. The worker consumes messages from a Kafka topic, processes them, and produces results back to another Kafka topic. It also interacts with a PostgreSQL database to manage task states.

## Features

- Kafka consumer and producer for message handling
- PostgreSQL integration for task management
- Background processing of due tasks
- Email service for notifications

## Setup

1. Clone the repository:

   ```bash
   git clone `repository_url`
   cd social-platform-kafka-worker
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Configure environment variables:
   Create a `.env` file in the root directory using `.env.example` as a template.

4. Configure config.yaml:
   Create a `config.yaml` file in `/config` directory using `config.example.yaml` as a template.

5. Set up Kafka worker using docker-compose:

   ```bash
   docker-compose up -d
   ```

6. Run this service:

   ```bash
   go run cmd/main.go
   ```
