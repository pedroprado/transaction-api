# transaction-api

API for creating financial transactions

## 1.Entities Diagram

![alt-text](/images/transactionAPI.jpg)

## 2.Use Cases

### i.Create Accounts

Create accounts for simulating the transaction

### ii.Create Transaction

Create a transaction, provisioning the resources in an intermediary account

### iii.Complete Transaction

Send the balance provisioned (transferred to the intermediary account) to the destination account.

(development) When the transference fails, created a rollback provisioning

### iv.Compensate Failed Transactions  (development)

Rolls back the transaction, moving the balance back from the intermediary account for the source account.

## 3.API Documentation

### i.Swagger

The swagger can be accessed by url (application must be running):

http://localhost:8098/snapfi/swagger/index.html

### ii.Postman Collection

A collection is available in the root directory with the name  **TransactionAPI.postman_collection.json** .

## 4.Running and Testing

### i.Unit Tests

In the root directory run **go test ./src/...**

### ii. Integration Tests

The Integration Tests are Broad Integration Tests, as described in https://martinfowler.com/articles/practical-test-pyramid.html.
Thus, the tests should cover all the behaviors of the whole application and the integration with database, external apis and services.

First, set up the docker environment: in the root director run the command **(sudo) sh compose.sh**

Second, run the application (using the .env file for environment variables):
 * Run with **go run src/main.go** or using you preferred IDE.
 * Run with docker using the command **(sudo) sh run.sh** in the root directory

#### Testing scenarios manually using Postman:

1. Create a Transaction
    * Create two accounts (origin and destination) using the account create account api (POST /account)
    * Create a transaction using the create transaction api ( POST /transactions)
    * Verify a status OPEN was created for the transaction using the api (GET /transaction_status&transaction_id={transaction_id}) 
    * Verify a balanceProvision OPEN was created for the transaction using the api (GET /balance_provisions&transaction_id={transaction_id})
    * Verify the balance was moved from origin account to intermediary account, using api (GET /account/{account_id})

2. Complete a Transaction
    * Recreate Scenario 1
    * Complete the transaction using the complete transaction api (POST /transaction/{transaction_id}/complete)
    * Verify the transaction status is BOOKED using the api (GET /transaction_status&transaction_id={transaction_id}) 
    * Verify the balanceProvision was CLOSED the api (GET /balance_provisions&transaction_id={transaction_id})
    * Verify the balance was moved from intermediary account to destination account, using api (GET /account/{account_id})

3. Failure Completing Transaction
   * Recreate Scenario 1
   * Remove funds from intermediary account using the patch account api (PATCH /account)
   * Complete the transaction using the complete transaction api (POST /transaction/{transaction_id}/complete)
   * Verify that the transaction status is set to FAILED using the api (GET /transaction_status&transaction_id={transaction_id})
   * Verify that the balanceProvision was of type ADD is set to CLOSED the api (GET /balance_provisions&transaction_id={transaction_id})
   * Verify that a new balanceProvision was of type VOID in OPEN status was created, using the api (GET /balance_provisions&transaction_id={transaction_id})

4. Compensate a Failed Transaction

#### Testing scenarios manually using Karate


#### Testing scenarios automated

You can run all the tests scenarios described automatically.

In the root directory, use run the script:

**(sudo) sh integration_test.sh**


