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

    curl -X GET localhost:2022/users/:username

### Sync single user with GitHub

    curl -H "Authorization: Bearer OAUTH-TOKEN" -X GET localhost:2022/users-sync/:username

## Repository Ressource

### Get repositories for a single user

    curl -X GET localhost:2022/repos/:username

### Sync respositories for a single user with GitHub

    curl -H "Authorization: Bearer OAUTH-TOKEN" localhost:2022/repos-sync/:username







