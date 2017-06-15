# user service
Work in Progress!

## Dependencies
Go
docker
docker-compose

## Usage
### Build
Just enter the project directory and type

    ./docker-build.sh

### Run
Just enter the project directory and type

    docker-compose up

The user service is now running on port `2022`
ES is running on port `9200`

## User Ressource

### Get single user

GET /users/:username

### Sync single user with GitHub

GET /users-sync/:username

## Repository Ressource

### Get repositories for a single user

GET /repos/:username

### Sync respositories for a single user with GitHub

GET /repos-sync/:username







