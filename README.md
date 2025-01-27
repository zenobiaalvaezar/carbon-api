# Carbon

[Feel free to try this API. Click to access the documentation.](https://carbon-api-70017640279.us-central1.run.app/swagger/index.html)

## Description: 

Carbon Management API goals is to calculate carbon (C02) emissions estimates for common C02 emitting activities.

## Background:

> Carbon API makes it easy to generate accurate emissions estimates from electricity consumption and fuel combustion.

## Highlights:

* Monolith Architecture
* Caching with Redis
* Serverless Deployment with Google Cloud Run
* Payment Gateway (Xendit)
* Email notifications

### Tech stacks:

* Go
* Echo
* Docker
* PostgreSQL
* MongoDB
* Redis
* GORM
* JWT-Authorization
* 3rd Party APIs (Xendit)
* SMTP
* REST
* Swagger
* Testify

## Application Flow

![Final Flow](./misc/flow.png)

## ERD

### Carbon Services (Postgres)

![Carbon services ERD](./misc/carbon_erd.png)

### Payment Method Service (MongoDB)

![ERD](./misc/payment_method_erd.png)

## Deployment

This app is containerized and deployed to Google Cloud Run as a monolith architecture. This means for each service (carbon_service, payment_method_service, and api_gateway) in a single instance. 

To deploy, go to the root folder for each service and type:

```bash
docker buildx build -t gcr.io/[PROJECT-ID]/[SERVICE-NAME]:latest --push .

gcloud run deploy carbon-api \
            --image gcr.io/[PROJECT-ID]/[SERVICE-NAME]:latest \
            --platform managed \
            --region [REGION] \
            --allow-unauthenticated \
            --timeout 10m \
```

- __PROJECT-ID__ refers to the project's ID on Google Cloud Console.
- __SERVICE-NAME__ refers to the service name that you want to deploy. Example: carbon-api.
- __REGION__ refers to the region name that you want to deploy. Example: us-central1.
