# your-money

### Deployment
Find the deployed app [here](https://buying-frenzy-dipto.fly.dev) .

### Run

#### Test
```shell
make test
```

#### Test with coverage report
This is open a browser tab with graphical coverage report
```shell
make test-coverage
```

#### Server
This app will be running on port 8080.
1. Run this command to run this project in docker.
   ```shell
    make development-serve
   ```

2. DB init script is already given. So tables and initial data will be loaded to db automatically.

3. `6d7750a1-c3f2-4765-bf8f-33bc80f3f809` this user is automatically inserted with the help of db init script and
    can be used to test the apis.
   ```json
     {
       "id": "6d7750a1-c3f2-4765-bf8f-33bc80f3f809",
       "created_at": "2023-04-09 17:00:42.705392+00",
       "updated_at": "2023-04-09 17:00:42.705392+00",
       "name": "Test",
       "balance": 100
     }
   ```

### Api
#### Add Balance
This endpoint is used to add balance to a user's account. this endpoint adds the balance in a transaction and does
a `select` query with `for update` expression to lock the selected rows for update so that no other concurrent 
connection doesn't read dirty row.

---
Method : `POST`
> /users/{uid}/add 

Query Params:
> N/A

##### Response - 202
```json
{
   "success": true,
   "message": "transaction successful!",
   "status_code": 202,
   "data": {
      "current_balance": 160
   }
}
```

##### Response - 422
* when transaction id was already processed. this is to handle
   if internal services can't guarantee caller will call only once.
```json
{
   "success": false,
   "message": "transaction was already processed",
   "status_code": 422
}
```

##### Response - 400
* invalid transaction id
```json
{
   "success": false,
   "message": "valid transaction id required",
   "status_code": 400
}
```

#### Balance
Provides a list of restaurant which have items which maintains the condition low_price >= price <= high_price
and the count of the items in the range for a restaurant is > more_than or count < less_than

---
Method : `GET`
> /users/{uid}/balance

Query Params:
> Optional:
> N/A

> Required:
> N/A

##### Response - 200
```json
{
   "success": true,
   "message": "request successful!",
   "status_code": 200,
   "data": {
      "balance": 160
   }
}
```

##### Response - 404
* User not found
```json
{
   "success": false,
   "message": "user not found",
   "status_code": 404
}
```

#### History
Provides a list of paginated transaction history for a user

---
Method : `GET`
> /users/{uid}/history

#### Query Params:

#### Required:

> page_size (ex - 10)

#### Optional:
> * page (ex - `Mw==`)

##### Response - 200
```json
{
   "success": true,
   "message": "request successful!",
   "status_code": 200,
   "data": {
      "total": 3,
      "page_size": 20,
      "next_page": "Mw==",
      "histories": [
         {
            "created_at": "2023-04-09T17:54:38.305802Z",
            "amount": 20,
            "transaction_id": "tx_1as4ndakdab"
         },
         {
            "created_at": "2023-04-09T17:06:21.624738Z",
            "amount": 20,
            "transaction_id": "tx_1as4ndakda"
         },
         {
            "created_at": "2023-04-09T17:06:06.202434Z",
            "amount": 20,
            "transaction_id": "tx_1as4ndakd"
         }
      ]
   }
}
```

##### Response - 

* Invalid `page_size` query param
```json
{
   "success": false,
   "message": "page size should be a valid integer",
   "status_code": 400
}
```
* Invalid `page` query param. should be a valid base64 encoded string.
```json
{
    "success": false,
    "message": "invalid pagination cursor",
    "status_code": 400
}
```
