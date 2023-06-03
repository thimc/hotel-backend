# API Documentation

All URL routes are prefixed with *http://localhost:3000/v1/*

## /auth
* Request type: POST
* Requires admin access: No

This route should be the first step in any interaction with the API and is
expecting the user to provide their email and password.

The following fields are required:
- email
- password

The response will contain the users information along with a token that will
be valid for 4 hours.

## /user
* Request type: POST
* Requires admin access: No

This route will create a new user and return the information if it succeeds.

The following fields are required:
- firstName
- lastName
- email
- password

## /user
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

## /user/*id*
* Request type: GET
* Requires admin access: no

This route will return one specific user and return the data in the similar
manner as the previous route.

## /user/*id*
* Request type: POST
* Requires admin access: no

This route allows the user to update any users personal information.

The following fields are valid:
- firstName
- lastName

## /user/*id*
* Request type: DELETE
* Requires admin access: no

This route will delete the user if it exists along with their personal information.

**NOTE: any reservations made by the user will not be removed**

## /hotel
* Request type: GET
* Requires admin access: no

This route will return all hotels that are known to the database.

The following fields are returned:
- \_id
- name
- location
- rooms (array)
- rating

## /hotel/*id*
* Request type: GET
* Requires admin access: no

This route will return one specific hotel in the same way as the previous route.

Database pagination is enabled for this route and is required.

`..../v1/hotel?page=1&limit=5&rating=3`

The result will be a list of 5 hotels starting from index 1 with a 3 star rating.

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

## /room
* Request type: GET
* Requires admin access: no

This route will return all known rooms.

The following fields are returned:
- \_id
- seaside
- size
- price
- hotel \_id

## /room/*id*
* Request type: GET
* Requires admin access: no

This route will return the room matching `id`.

The following fields are returned:
- \_id
- seaside
- size
- price
- hotel \_id


## /room/*id*/book
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

## /booking/*id*
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

## /booking/*id*/cancel
* Request type: POST
* Requires admin access: no

This route will cancel any reservations that matches the *id*.

## /booking
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
