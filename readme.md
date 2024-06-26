# AssetFindr Test
## Posts and tags application


### summary

this application is created with go, gin, gorm and postgres sql. it also follows base practice for creating simple api,error handling, testing, and mocking.

### installation

to run the stacks, the project needs to create `.env` file, you can copy from example.env.

and run the docker compose, please make sure you install the `docker` and `make` 
```
make stack-up
```

to run the application need to run the go, please make sure you install the `go` language 
```
make run-http
```

and application ready to serve in desired port

### List of API

this projects is using http api

- get list of post 

to get list of post 
```
GET {{API_ENDPOINT}}/api/posts
```
can be invoked with
```curl
curl --location 'http://{{API_ENDPOINT}}/api/posts'
``` 
and the response will look like this
```json
{
    "data": [
        {
            "id": 86,
            "title": "a",
            "content": "a",
            "tags": [
                "v",
                "b"
            ]
        },
        {
            "id": 85,
            "title": "f",
            "content": "g",
            "tags": [
                "h",
                "i"
            ]
        }
    ],
    "result": "ok"
}
```

- get post by id

to get 1 post by id 
```
GET {{API_ENDPOINT}}/api/posts/1
```
can be invoked with
```curl
curl --location 'http://{{API_ENDPOINT}}/api/posts/1'
``` 
and the response will look like this
```json
{
    "data": {
        "id": 76,
        "title": "Lorem 12",
        "content": "a",
        "tags": [
            "Lorema",
            "a"
        ]
    },
    "result": "ok"
}
```

- create post

to create post by id 
```
POST {{API_ENDPOINT}}/api/posts
```
can be invoked with
```curl
curl --location 'http://{{API_ENDPOINT}}/api/posts' \
--header 'Content-Type: application/json' \
--data '{
 "title": "Lorem",
"content": "test",
 "tags":["ipsum"]
}'
``` 
and the response will look like this
```json
{
    "result": "created"
}
```

- update post

to update post by id 
```
PUT {{API_ENDPOINT}}/api/posts/1
```
can be invoked with
```curl
curl --location --request PUT 'http://{{API_ENDPOINT}}/api/posts/86' \
--header 'Content-Type: application/json' \
--data '{
 "title": "Upda",
 "content": "Upda",
 "tags": [ "Ipsum1000", "ac"]
}'
```
and the response will look like this
```json
{
    "result": "updated"
}
```


- delete post

to delete post by id 
```
DELETE {{API_ENDPOINT}}/api/posts/1
```
can be invoked with
```curl
curl --location --request DELETE 'http://{{API_ENDPOINT}}/api/posts/86' \
--header 'Content-Type: application/json' \
``` 
and the response will look like this
```json
{
    "result": "deleted"
}
```


