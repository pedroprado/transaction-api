# transaction-api

API for creating financial transactions

## Entities and Use Cases

![alt-text](/images/transactionAPI.jpg)

## Use Cases
### 1.Create Accounts

Create accounts for simulating the transaction

### 2.Create Transaction

Create a transaction, provisioning the resources in an intermediary account

### 3.Complete Transaction

Send the balance provisioned (transferred to the intermediary account) to the destination account.

(development) When the transference fails, created a rollback provisioning


### 4.Compensate Transaction  (development)

Rolls back the transaction, moving the balance back from the intermediary account for the source account.

## Running and Testing

### 1.Unit Tests

In the root directory run **go test ./src/...**

### 2. Integration Tests

Run the docker enviromnent: in the root director run the command **(sudo) sh compose.sh**

Run the application (using the .env file for environment variables). Run with **go run src/main.go** or using you prefered IDE.

Run the **karate features**:
