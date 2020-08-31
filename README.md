vehicleSimulator Code Test
=================

## Summary

A simple web server, that can store information about users, create, update and delete user. 
By default, it will start on port 8080. 
Samples of requests are below

## Start service
  ### Docker
  In order to srart the service, a `Docker` should be installed. 

 How to install `Docker` 
  ``` 
  https://docs.docker.com/get-docker/
```
  If `Docker` is installed, please, run command,s from service root directory
  ```
  docker-compose -f docker/docker-compose.yaml up
``` 
  or 
  ```
  start_docker.sh
```

  This command will run 3 docker containers on a local computer. 
  - container with `postgres` (service DB)
  - container with `dependencies` (will run migrations)
  - container with `service` (expose port :8080)

  ### Localy (macos and linux)
With default config it strts on port 8080 by running `start_local.sh` script.  
This script will start a docker container, install migration tool `dbmate`, run migrations and will start a server it'self. 


## Example JSON Requests
  ### Health
  List of servise dependencies (ak databases, queues etc.). Accepts query parameter `service`, to retrive health of only one registered dependency.

  Request:
   ```GET: http://localhost:8080/v1/health```

  Response:
  ```javascript
[
    {
        "name": "postgres",
        "healthy": true,
        "time": "2020-07-20T02:59:34.105781Z"
    }
]
```

  ### Countries
  List of countries registered in service. Supposably be used on Create/Update user API calls, to prevent users from naming Countries differently.
  Stored in separate table in database and referenced by user table. 
  As an improvement, admit API should be added, to be ablr to manage list of countries. Also cache should be valid, before making a request to backend.
  Improve validation

  Request:
```GET: http://localhost:8080/v1/countries```

  Response:
```javascript
[
    {
        "ID": 1,
        "ISO2": "AF",
        "Name": "Afghanistan"
    },
    {
        "ID": 2,
        "ISO2": "AL",
        "Name": "Albania"
    }
    ...
]
```

  ### Create user
  Create user accepts json body with user parameters. Countries should be passed as an integer value (id) to reduse load on server and manage necessary 
  relations in DB. Password will not bre retrived in response, as it was hashed, salted and after that saved in DB. For future password checks, same procedure 
  will happened, and hashed and salted values will be compared. 
  As an improvement, fields validation should check max lenght for each field and email must be confirmed.

  Request:
```POST: http://localhost:8080/v1/user```

  Body:
```javascript
{
   "first_name":"David",
   "last_name":"Bowie",
   "nick_name":"star man",
   "email":"davidbowie@gmail.com",
   "password":"vanillaice",
   "country":231
}
```

Response: 
```javascript
{
   "first_name":"David",
   "last_name":"Bowie",
   "nick_name":"star man",
   "email":"davidbowie@gmail.com",
   "country":"GB"
}
```
  

  ### Update user
  Update user accepts json body with user parameters. Countries should be passed as an integer value (id) to reduse load on server and manage necessary 
  relations in DB. Password will be used to confirm user's identity. Password update will be made with a different REST call
  As an improvement, fields validation should check max lenght for each field and email must be confirmed.

  Request:
```PUT: http://localhost:8080/v1/user/{id}```

  Body:
```javascript
{
   "first_name":"David",
   "last_name":"Bowie",
   "nick_name":"star man",
   "email":"davidbowie@gmail.com",
   "password":"vanillaice",
   "country":231
}
```

  Response: 
  ```
  Status Code
```

  ### Update user's password
  Update user's password accepts json body with old and new passwords. Before update, service is checking if password, stored in DB mathes old password from 
  request, and if so, proceeds with update.

  Request:
```PUT: http://localhost:8080/v1/user/{id}/password```

  Body:
```javascript
{
   "new_password":"abcd",
   "old_password":"cdba"
}
```

  Response: 
  ```
  Status Code
```

  ### Delete user
  Deletes user's info frem database.
  As an improvement, user identity confirmation should be added, and admin access to remove user.

  Request:
```DELETE: http://localhost:8080/v1/user/{id}```

  Response: 
  ```
  Status Code
```

  ### Get All users
  Get all users retrives information about all users, stored in database. Users can be filtered by `counrty`, `first_nae`, `last_name` and `nick_name`,
  as request accepts query parameters `title` and `filter`. 
  As an improvement, paggination should be added, and possibility to use more, than one filter by request.

  Request:
```GET: http://localhost:8080/v1/user```

Response: 
```javascript
[
   {
      "first_name":"David",
      "last_name":"Bowie",
      "nick_name":"star man",
      "email":"davidbowie@gmail.com",
      "country":"GB"
   },
   {
      "first_name":"Amy",
      "last_name":"Lee",
      "nick_name":"Gothic princess",
      "email":"amylee@gmail.com",
      "country":"US"
   }
]
```
  
  ### Get One user
  Get on users retrives information about one, by it's id users, stored in database.

  Request:
```GET: http://localhost:8080/v1/user/{id}```

Response: 
```javascript
   {
      "first_name":"David",
      "last_name":"Bowie",
      "nick_name":"star man",
      "email":"davidbowie@gmail.com",
      "country":"GB"
   }
``` 

## Assumptions during development
  ## Database

  Based on project description and data, service will be working with, PostgreSQL was chosen as a database engine. 
  It can provide a good mechanism for working with user information and any other related tada, and quick access to it.
  In current implementation, 3 tables are used. 
    
  Table `countries` to store list of countries, just to keep country management for all users the same. Presumably, during
  `create` or `update` users requests, UI should make an http call to get list of all countries, and allow user to chose 
  one from it

  Table `users` stores publicly availabe user information (can be adjusted).

  Table `users_password` stores hashed and salthed users password, that is referencing to `users` table by field `user_id`. 
  Since publicly available user's information, presumably, will be accessed more often, than password data, separation of those tables
  will speed up reads of publicly available user's information. 

  ## Service
  A web service is desinged to provide a CRUD operations on `User` entity. A multiple nodes can be added, with load balancer and 
  distributed cache before service (no cache implementation at this moment). It is designed to be able to change Database engine, if there 
  is a need for that, without changing buisenes logick of the service.   

  ## Notifier
  Notifier package providing an interface, which will allow to notify other services about events, that have happened in current service.
  Based on configuration and interface implementation, differet approaches and protocols can be used, to comunicate with different services.
  `httpclient` implementation, which is provided, creates a simle `http.Client` and notifies services via `http`.

## Improvement on servise
  Add unit tests. Add config with enviropment variables. Add Authentefication and authorisation mechanisms.
  Add load balancing before service, cache for static data. Add `swagger`.


