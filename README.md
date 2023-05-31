# Hotel backend

This repo contains a demo hotel backend written in go and is powered by fiber
and uses MongoDB as its database. Authentication is done via JWT.

All URL routes are prefixed with *http://localhost:3000/v1/*

NOTE: The backend is still missing a few routes and error handling needs to be tidied
up a bit before this project is considering feature complete.

See the [API Documentation](DOCUMENTATION.md) for a complete map of all available routes.

## Setting up the environment
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

### Mongo Express Web UI (Optional)
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
