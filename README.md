# 432_final_project

## Overview

- Postgres lives on a single VM instance -- specifically using an e2-micro machine type and the base image for postgis [found here](https://gorm.io/index.html)
- Each table gets its own go file under the `cmd` subfolders separated by the frequency of extraction needed.
- Helper functions used throughout multiple programs live in the `pkg` subdir
    - custom functions for floats and datetimes are in `type_helpers.go`, Chicago Data Portal labels numeric data generically, we can too.
    - During extraction, `soda.go` filters the intake and funnels json data into specific structs, since we know the data coming in
    - `postgres.go` handles the loading with special help from the gorm package.
- Transformations shape the raw data into the format needed for the postgres tables
    - Making sense of the types from JSON and the types we need to combine things for postgres is handled through custom types
- Loading enforces type constraints & delivers records to postgres
- Geocoding is done through postGIS and queries after loading raw data in, matching on JOINS with geoJSON polygons for ZIP code and Community Area.
    - There is a bit of setup here but it only needs to happen once and it's very cost effective.
- The deployment uses GCP secrets management in one VM (or locally) to connect to the postgres database and continuously extract data on a batch schedule.
- The cron execution all happens in `app.go` and the program runs and sleeps continuously.


## Dependencies (necessary for developing locally):
- Google cloud console
- Postgres
- Go version 1.22.0
- ogr2ogr (installing GDAL)
- the following parameters as secrets

I created a .env file like this
```
POSTGRES_DB=YourDBName
POSTGRES_HOST=VMExternalIP
POSTGRES_USER=postgres
POSTGRES_PASSWORD=YOUR_PASSWORD
POSTGRES_PORT=5432
```
then added the json equivalent in GCP  

```
{
  "POSTGRES_DB": "mydatabase",
  "POSTGRES_HOST": "localhost",
  "POSTGRES_USER": "myuser",
  "POSTGRES_PASSWORD": "mypassword",
  "POSTGRES_PORT": "5432"
}
```

the .env file helps with tests locally.
The GCP secret ensures that the program can actually run on the production db

GCP needs a service account with access to secrets management for this to work.

## Set up -- postgres

- In a GCP project, create a VM instance under the Compute Engine
- Start the machine without a container.
- This [tutorial](https://joncloudgeek.com/blog/deploy-postgres-container-to-compute-engine/) helped set up the necessary configurations for the VM (up to the point where they talk about db migrations) -- the firewall rules are key here.
- Instead of using the pre-built Google container, we are submitting a custom container with PostGIS enabled. 
    - this is to allow for geocoding in the database when executing queries.

- Starting the container requires setting this up in the VM -- SSH into it and to the following
```
sudo apt-get update &&
sudo apt-get install -y docker.io &&
sudo docker pull postgis/postgis &&
sudo docker run --name <container_name> -e POSTGRES_PASSWORD=<your_password> -p 5432:5432 -d postgis/postgis
```

get ogr2ogr isntalled in the container (Not the VM) with
```
docker exec -it <container_id_or_name> bash
apt-get update && apt-get install -y gdal-bin && rm -rf /var/lib/apt/lists/*
```

confirm that's working with 

```
`ogr2ogr --version`
```

run "exit" to leave the container

then you should be able to upload the geojson files to the VM & copy them to the container with

```
sudo docker cp "Boundaries - Community Areas (current).geojson" <container_id_or_name>:/"Boundaries - Community Areas (current).geojson"
sudo docker exec -it <container_id_or_name> bash
ogr2ogr -f "PostgreSQL" PG:"dbname=<your_db> user=<your_username> password=<your_password>" /file.geojson
```

easiest to have the files represent what you want the table name to be -- i.e. boundaries_community_areas

Then get the new database set up -- still SSHd into the VM

```
docker exec -it <container_name/id> bash &&
psql -U postgres
```

then in psql:
```
CREATE DATABASE <db_name>;
CREATE DATABASE test_db;
\c <db_name>

CREATE EXTENSION postgis;
```

- confirm the database is set up 
- This process lets you define a db_name and password, the default user is postgres
- Confirm the connection outside the VM on your local machine with

```
psql -h <external_ip> -U postgres -d <db_name>
```

You can copy+paste the create table commands into the psql shell 
I tried to find a way to just run init.sql but got stumped. 

But you'll end up with the following

```
                   List of relations
 Schema |            Name            | Type  |  Owner
--------+----------------------------+-------+----------
 public | boundaries_community_areas | table | postgres
 public | boundaries_zip_codes       | table | postgres
 public | building_permits           | table | postgres
 public | covid_cases                | table | postgres
 public | geographies                | table | postgres
 public | spatial_ref_sys            | table | postgres
 public | taxi_rideshares            | table | postgres
 public | traffic_estimates          | table | postgres

```

- Copy geojson files into the container
`docker exec -it <container_id> bash`

- Running `app.go` will populate the database 
- then a query like this can obtain the zip codes and community areas through PostGIS

```
SELECT taxi_rideshares.*, boundaries_zip_codes.zip
FROM taxi_rideshares
JOIN boundaries_zip_codes 
ON ST_Contains(
    boundaries_zip_codes.wkb_geometry, 
    ST_SetSRID(ST_Point(taxi_rideshares.pickup_centroid_longitude, taxi_rideshares.pickup_centroid_latitude), 4326)
);
LIMIT 10
```
And with that, the geospatial data should be ready for analysis on the front end!

## Setup -- go code locally

get a google service account and connect it to your enviornment variables
export GOOGLE_APPLICATION_CREDENTIALS="service-account.json"

make the service account able to access secrets
I did this going to cloud IAM and adding a secrets manager accessor responsibility

And followed the steps [here](https://cloud.google.com/docs/authentication/provide-credentials-adc#how-to)
```

docker build -t your-image-name
docker run -v /path/to/application_default_credentials.json:/application_default_credentials.json -e GOOGLE_APPLICATION_CREDENTIALS=/application_default_credentials.json your-image-name

```
## Setup -- containerized deployment
- this assumes a private github repo is set up for the project

- creating a new VM instance (also e2 micro)
- setting up docker 

```
sudo apt-get update &&
sudo apt-get install -y docker.io
```

- providing a connection to github

```
ssh-keygen -t ed25519 -C "your_email"
eval "$(ssh-agent -s)"
ssh-add ~/.ssh/id_ed25519
cat ~/.ssh/id_ed25519.pub
```

Copy the output of the last command and add it to GitHub account

Then you should be able to clone
and be able to build with the same steps

```
sudo docker build -t 432-final . && sudo docker run 432-final
```

Assumes my private repository, but changing the go code to your solo reference should work too.

## Dependencies (handled by Go):
- [go-soda](https://pkg.go.dev/github.com/SebastiaanKlippert/go-soda@v1.0.1)
- [gorm](https://gorm.io/index.html)
- [secretmanager](https://pkg.go.dev/cloud.google.com/go/secretmanager)
- [cron](https://github.com/robfig/cron)


## Still to do:
- Frontend deployment