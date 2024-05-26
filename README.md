# 432_final_project

## Overview

- Postgres lives on a single VM instance -- specifically using an e2-micro machine type and the base image for postgres version 15 [found here](https://console.cloud.google.com/marketplace/product/google/postgresql15?hl=en&project=final-project-424101)


## Set up -- postgres

- In a GCP project, create a VM instance under the Compute Engine
- Start the machine with a container using the image `marketplace.gcr.io/google/postgresql15:latest`
- This [tutorial](https://joncloudgeek.com/blog/deploy-postgres-container-to-compute-engine/) helped set up the necessary configurations for the VM & container (up to the point where they talk about db migrations)
- confirm the database is set up 
- This process let's you define a db_name and password, the default user is postgres
- Confirm the connection with

```
psql -h <external_ip> -U postgres -d <db_name>
```

You can copy+paste the create table commands into the psql shell 
I tried to find a way to just run init.sql but got stumped. 

But you'll end up with the following

```
               List of relations
 Schema |       Name        | Type  |  Owner
--------+-------------------+-------+----------
 public | building_permits  | table | postgres
 public | covid             | table | postgres
 public | geographies       | table | postgres
 public | taxi_rideshares   | table | postgres
 public | traffic_estimates | table | postgres

```

## How the Go ETL works

Each table gets its own go file under the `cmd` subfolders separated by the frequency of extraction needed. 

The pipeline consists of the following steps

- Extraction:
    - 
    - Socrata API [docs for SoQL](https://dev.socrata.com/docs/queries/) (used for the extraction queries)


## Dependencies:
- [go-soda](https://pkg.go.dev/github.com/SebastiaanKlippert/go-soda@v1.0.1)
- [gota](https://pkg.go.dev/github.com/go-gota/gota/dataframe)
- [pq](https://pkg.go.dev/github.com/lib/pq#section-readme)
- [godotenv](https://pkg.go.dev/github.com/joho/Godotenv)
- Postgres version 15



## Still to do:
- Apply extraction/loading utils to the other tables
- Dockerize & Deploy
- Frontend deployment