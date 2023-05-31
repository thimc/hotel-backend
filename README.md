# Hotel backend

This repo contains a demo hotel backend written in go and is powered by fiber
and uses MongoDB as its database. Authentication is done via JWT.

All URL routes are prefixed with *http://localhost:3000/v1/*

NOTE: The backend is still missing a few routes and error handling needs to be tidied
up a bit before this project is considering feature complete.

## Routes

### /auth
* Request type: POST
* Requires admin access: No

This route should be the first step in any interaction with the API and is
expecting the user to provide their email and password.

The following fields are required:
- email
- password

The response will contain the users information along with a token that will
be valid for 4 hours.

### /user
* Request type: POST
* Requires admin access: No

This route will create a new user and return the information if it succeeds.

The following fields are required:
- firstName
- lastName
- email
- password

### /user
* Request type: GET
* Requires admin access: no

This route will return all the users

The following fields are returneed:
- \_id
- firstName
- lastName
- email
- encryptedPassword
- isAdmin

### /user/*id*
* Request type: GET
* Requires admin access: no

This route will return one specific user and return the data in the similar
manner as the previous route.

### /user/*id*
* Request type: POST
* Requires admin access: no

This route allows the user to update any users personal information.

The following fields are valid:
- firstName
- lastName

### /user/*id*
* Request type: DELETE
* Requires admin access: no

This route will delete the user if it exists along with their personal information.

**NOTE: any reservations made by the user will not be removed**

### /hotel
* Request type: GET
* Requires admin access: no

This route will return all hotels that are known to the database.

The following fields are returned:
- \_id
- name
- location
- rooms (array)
- rating

### /hotel/*id*
* Request type: GET
* Requires admin access: no

This route will return one specific hotel in the same way as the previous route.

### /hotel/*id*/rooms
* Request type: GET
* Requires admin access: no

This route will return all known rooms to one specific hotel.

The following fields are returned:
- \_id
- seaside
- size
- price
- hotel \_id

### /room
* Request type: GET
* Requires admin access: no

This route will return all known rooms.

The following fields are returned:
- \_id
- seaside
- size
- price
- hotel \_id

### /room/*id*/book
* Request type: POST
* Requires admin access: no

This route will reserve the room that has the *id* and return the data structure
if it succeeds.

The following fields are required:
- \_id
- userID
- roomID
- numPersons
- fromDate
- untilDate
- canceled

### /booking/*id*
* Request type: GET
* Requires admin access: no

This route will return any reservations that matches the *id*.

The following fields are returned:
- \_id
- userID
- roomID
- numPersons
- fromDate
- untilDate
- canceled

### /booking/*id*/cancel
* Request type: POST
* Requires admin access: no

This route will cancel any reservations that matches the *id*.

### /booking
* Request type: GET
* Requires admin access: yes

This route will return any known reservations.

The following fields are returned:
- \_id
- userID
- roomID
- numPersons
- fromDate
- untilDate
- canceled

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
