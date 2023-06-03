# Hotel backend

This repo contains a demo hotel backend written in go and is powered by
[fiber](https://gofiber.io/) and uses [MongoDB](https://www.mongodb.com/) as
its database. Authentication is done via JWT.

The backend is designed with portability in mind, should the user ever need to
switch to another type of database such as PostgreSQL, SQLite, MariaDB or MySQL
(it could be literally any way of storing actually) all that would be needed
for the switch is to implement the functions needed for the interfaces in the
`db` folder.

Head on over to the [API Documentation](DOCUMENTATION.md) for a complete map
of all available routes.

## Environment variables

Configure the `.env` file:
```
HTTP_LISTEN_ADDRESS=:3000
JWT_SECRET=abcdefghijklmnopqrstuvwxyz
MONGODB_DB_NAME=hotel-reservation
MONGODB_DB_URI=mongodb://1.2.3.4:27017
MONGODB_TEST_DB_URI=mongodb://1.2.3.4:27017
```

## Installing the dependencies

In order to set the environment variables via the `.env` file, the gotdotenv
package is needed:

```
go get github.com/joho/godotenv
```

In order to get this repo up and running you will need to grab the MongoDB
and fiber libraries and install docker.
```
go get go.mongodb.org/mongo-driver/mongo
go get github.com/gofiber/fiber/v2
```

The following snippet will spin up a MongoDB instance:
```
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

When everything is installed and the docker containers are running,
run the following command:

```
make run
```

In order to set up a test environment which will seed the database with
fake data, run:
```
make seed
```

## Mongo Express Web UI (Optional)
If you want to interact directly with the database, the Mongo Express web
user interface is excellent, installation is done with the following command:
```
docker run -it -p 8081:8081 \
        -e ME_CONFIG_MONGODB_URL="mongodb://mongodb:27017" \
        -e ME_CONFIG_BASICAUTH_USERNAME="" \
        -e ME_CONFIG_BASICAUTH_PASSWORD="" \
        -e ME_CONFIG_OPTIONS_EDITORTHEME="default" \
        -e ME_CONFIG_SITE_COOKIESECRET="cookiesecret" \
        -e ME_CONFIG_SITE_SESSIONSECRET="sessionsecret" \
        -e ME_CONFIG_SITE_SSL_ENABLED="false" \
        -e ME_CONFIG_SITE_CRT_PATH="" \
        -e ME_CONFIG_SITE_KEY_PATH="" \
        -e ME_CONFIG_MONGODB_ENABLE_ADMIN=true \
        mongo-express
```

The web interface should now be up and running on:
```
http://localhost:8081
```

## Deploying
The Makefile has a `docker` target to it so all that is needed to spin it up
is to run the following command:
```
make docker
```
