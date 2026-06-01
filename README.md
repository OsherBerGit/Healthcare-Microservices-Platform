# 🏥 HealthSaga: Distributed Healthcare Architecture

![Spring Boot](https://img.shields.io/badge/Spring_Boot-6DB33F?logo=spring&logoColor=white)
![FastAPI](https://img.shields.io/badge/FastAPI-009688?logo=fastapi&logoColor=white)
![Apache Kafka](https://img.shields.io/badge/Kafka-231F20?logo=apachekafka&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-DC382D?logo=redis&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?logo=postgresql&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-4EA94B?logo=mongodb&logoColor=white)
![Keycloak](https://img.shields.io/badge/Keycloak-JBoss-blue?logo=keycloak)
![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=white)

## 📖 About
**HealthSaga** is a highly scalable, event-driven microservices ecosystem designed for modern healthcare management. The system seamlessly integrates a high-performance API Gateway, multiple domain-specific microservices, and an AI-powered medical triage engine. 

Built strictly on **Domain-Driven Design (DDD)** principles, the project showcases enterprise-grade patterns such as the **Saga Pattern** for distributed transactions, real-time event streaming, and advanced security mechanisms including Zero-Trust IAM and Data Masking.

## 🛠 Tech Stack

### API Gateway (The Entry Point)
* **Core:** Golang
* **Features:** WAF Middleware, Data Masking, Dynamic Load Balancing, Intelligent Cost-Based Router.

### Core Microservices (Business Logic)
* **Core:** Java 17, Spring Boot 3
* **Architecture:** Hexagonal Architecture, Domain-Driven Design (DDD).
* **Services:** Patient, Appointment, Billing, Audit, Notification.

### Machine Learning Engine (Triage)
* **Core:** Python 3.10, FastAPI
* **Model:** Tiny ML Prediction Model (based on [Kaggle Patient Dataset](https://www.kaggle.com/datasets/mitishaagarwal/patient)).

### Infrastructure & Data
* **Event Streaming:** Apache Kafka (Outbox Pattern & Event Choreography).
* **Databases:** PostgreSQL (Relational), MongoDB (NoSQL).
* **Caching:** Redis (High-Performance Caching).
* **IAM:** Keycloak (Identity & Access Management).
* **DevOps:** Docker, Docker Compose.

## ✨ Technical Highlights & Features

### 🧩 Enterprise Microservices Architecture
* **Domain-Driven Design (DDD):** Deep separation of concerns. The Patient Service utilizes Hexagonal Architecture (Ports and Adapters) to isolate domain logic from infrastructure.
* **Distributed Transactions:** Full implementation of the **Saga Pattern**. Orchestrates complex admission processes (Patient Creation ➔ Appointment Scheduling ➔ Triage Prediction) with automated **Compensating Transactions (Rollbacks)** if any step fails.

### 🛡️ Advanced API Gateway (Golang)
* **Custom WAF & Security:** Built-in Web Application Firewall middleware and on-the-fly Data Masking for sensitive PII/PHI.
* **Smart Routing:** Intelligent Cost-Based Router and custom Load Balancing algorithm with active Health Checks.

### ⚡ Event-Driven & High Performance
* **Kafka Cluster:** Decouples core services from side-effects. Appointment creations emit events consumed asynchronously by Audit, Billing, and Notification services.
* **Redis Integration:** High-performance caching layers to reduce database hits on frequently accessed medical records.

### 🧠 AI Medical Triage
* **FastAPI ML Endpoint:** An isolated microservice evaluating patient vitals (Age, Heart Rate, Blood Pressure, Temperature) to predict urgency, deeply integrated into the Saga admission flow.

## 🚀 Quick Start

### 1. Prerequisites
* Docker & Docker Compose
* Go 1.21+
* Java 17 (JDK)
* Python 3.10+

### 2. Spin Up Infrastructure
Start the required databases, Kafka cluster, Keycloak, and Redis using Docker Compose:
```bash
docker-compose up -d
```

### 3. Build & Run Services
You can run the entire ecosystem via Docker, or individually for development:

**API Gateway (Golang):**
```bash
cd api-gateway
go run main.go loadbalancer.go
```

**Patient Service (Java/Spring Boot):**
```bash
cd patient-service
./mvnw spring-boot:run
```

**Triage Service (Python/FastAPI):**
```bash
cd triage-service
pip install -r requirements.txt
uvicorn api.main:app --port 8000
```

### 4. Trigger the Saga Admission Process
Send a POST request to the Gateway to initiate the Distributed Transaction:

```http
POST http://localhost:8080/api/admission/full-process
Content-Type: application/json
Authorization: Bearer <Keycloak_Token>

{
  "name": "Israel Israeli",
  "age": 45,
  "date": "2026-06-05",
  "heart_rate": 85,
  "blood_pressure_systolic": 120,
  "temperature": 37.2
}
```

### 📁 Project Structure
```text
📦 HealthSaga
 ├── 📂 api-gateway/          # Golang Gateway, LB, WAF, Saga Orchestrator
 ├── 📂 patient-service/      # Java Spring Boot, Hexagonal Arch, MongoDB
 ├── 📂 appointment-service/  # Java Spring Boot, PostgreSQL, Kafka Producer
 ├── 📂 triage-service/       # Python FastAPI, Tiny ML Model
 ├── 📂 billing-service/      # Java Spring Boot, Kafka Consumer
 ├── 📂 audit-service/        # Java Spring Boot, Kafka Consumer
 ├── 📂 notification-service/ # Java Spring Boot, Kafka Consumer
 └── 📜 docker-compose.yml    # Infra: Postgres, Mongo, Redis, Kafka, Keycloak
```

---

Developed as an advanced showcase of Distributed Systems, Zero-Trust Security, and Event-Driven Architecture.