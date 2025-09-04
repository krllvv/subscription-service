# Subscription REST Service

### Project setup

1. #### Clone repository
```bash
  git clone https://github.com/krllvv/subscription-service
  cd subscription-service
```

2. ####  Copy values from file `.env.example` to `.env` and update if necessary.

3. #### Run application using Docker
```bash
  docker-compose up --build
```

The application will run on `http://localhost:8080`

Swagger documentation will be available at `http://localhost:8080/swagger/index.html`

### Endpoints

| Method | Path                    | Description                 |
|:------:|:------------------------|-----------------------------|
|  POST  | `/subscriptions`        | Create subscription         |
|  GET   | `/subscriptions`        | List of subscriptions       |
|  GET   | `/subscription/{subID}` | Get subscription by ID      |
|  PUT   | `/subscription/{subID}` | Update subscription         |
| DELETE | `/subscription/{subID}` | Delete subscription         |
|  GET   | `/subscriptions/total`  | Sum total cost for a period |


Create `curl` example:
```bash
curl -L -X POST 'http://localhost:8080/subscriptions' \
-H 'Content-Type: application/json' \
-d '{
    "service_name":"Yandex Plus",
    "price":400,
    "user_id":"60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date":"07-2025"
  }'
```