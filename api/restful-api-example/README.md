# About

Example of a RESTful web service API deployed on AWS ECR and ECS using
Terraform. For learning and later copy/paste purposes.

* ECR - fully managed Docker container registry
* ECS - fully managed container orchestration service
* Terraform - open source infrastructure as code tool

Based on: https://medium.com/avmconsulting-blog/how-to-deploy-a-dockerised-node-js-application-on-aws-ecs-with-terraform-3e6bceb48785

# Steps

* [x] Write the app
* [x] Containerize the app
---
* [x] Create ECR, ECS, load balancer and related resources
* [x] Push the app image to ECR
---
* [ ] Add test driven security (CodeQL)
* [ ] Apply AWS WAF
* [ ] Collect logs and create security alerts

## Write the app

https://go.dev/doc/tutorial/web-service-gin

The app

```
go run .

curl -XPOST localhost:8080/albums --data @album4.json
```

## Containerize the app 

https://docs.docker.com/language/golang/build-images/

```
docker build --tag restful-api-example .
docker image ls
docker run --publish 8080:8080 restful-api-example
```

## Create ECR, ECS, load balancer and related resources

```
terraform init
terraform apply
```

## Push the app image to ECR

Amazon ECR > Repositories > restful-api-example: View push commands

EC2 > Load Balancers: DNS name

# Cleanup

```
terraform destroy
```