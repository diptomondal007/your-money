# your-money

### Deployment
Find the deployed app [here](https://your-money.fly.dev/) .

### Q/A
* How could the service be developed to handle thousands of concurrent users
  with hundreds of transactions each?

    **Ans** - Scalability is important consideration for our ledger service as our system
will handle thousands of concurrent users with hundreds of transaction each. Some key factors we need to
consider to make this system scalable -
  * We need to make this system horizontally scalable as vertical scalability does have a limit. Like using container orchestration
  systems can be used to make the system easily scalable.
  * Database needs to be designed efficiently. Like adding proper indexing,
  partitioning, sharding can heavily improve the scalability of the system.
  * As this is going to be a financial service `ACID` needs to be maintained
  * Using caching mechanism where it can be used can improve the performance also can
reduce the pressure on DB also makes the response time low hence improves the UX.
  * When the system becomes big it become impossible to monitor by only logs. Need to 
do proper monitoring and alerting with tools like `prometheus/grafana`. We can use
`Sentry` like tools which can help to track bugs in production environments.

* What has to be paid attention to if we assume that the service is going to run in
  multiple instances for high availability?
  
  **Ans** -
  * If this service runs in multiple instances to scale, it's so important
  to maintain `ACID` properties because data needs to be consistent
  across all the instances. Otherwise, for example if we fail to maintain
  `ACID` if user's account has a balance of 100$ and multiple instance tries to add
  10$ one instance will see the balance as 100$ and one will see as 110$ which will end up being a mess.
  * To grow confidence and reliability we need to design load testing,
  unit testing, integration testing, disaster recovery testing etc.

* How does the the add endpoint have to be designed, if the caller cannot
  guarantee that it will call exactely once for the same money transfer?
  
  **Ans** - 
  * If the caller can't guarantee that it will call exactly once for the same money transfer
it's important to design the endpoint to be idempotent. That will make sure
how many times the endpoint is being called for the transaction, the 
transaction will be processed only once.  We can generate a unique hash
for every transaction and save that to our transaction history and if the endpoint is called
multiple times for the same transaction identifier we can simply return `422` error
or just ignore the transaction.
  * This endpoint needs to handle potential error in any step of the whole process of the transaction.
This endpoint should process the whole flow or discard all. That can be done by db `transaction` system.
This will help to guarantee there is no inconsistency in data because of any error in any step of the whole process
which we call `ACID`.

    

### Run

#### Test
```shell
make test
```

#### Test with coverage report
This is open a browser tab with graphical coverage report
```shell
make test.coverage
```

#### E2E Test
This will spin up docker containers and will run the e2e tests
```shell
make test.integration
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
