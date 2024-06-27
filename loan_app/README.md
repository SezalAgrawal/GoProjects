# Loan App

It is an app that allows authenticated users to go through a loan application. It doesn’t have to contain too many fields, but at least “amount
required” and “loan term.” All the loans will be assumed to have a “weekly” repayment frequency.
After the loan is approved, the user must be able to submit the weekly loan repayments. It can be a simplified repay functionality, which won’t
need to check if the dates are correct but will just set the weekly amount to be repaid.

## Getting Started

### Prerequisites

The following are the tools and packages required for local development:

1. Install and setup Golang local development environment.
1. Install and setup Postgres
1. Install and setup Ruby (Use `rvm` to manage ruby version)

### Project Setup

1. Clone the repo to your local machine

    ```sh
    > git clone https://github.com/goProjects/loan_app.git
    ```

1. Create a dev environment file from the sample file

    ```sh
    > cp env.sample development.env
    ```

   Modify the values in `development.env` to match your needs.

1. Use helper bash script to run locally

    ```sh
    > bash ./scripts/run-local-server.sh
    ```

### Running Test Suite

Create a dev environment file from the sample file

   ```sh
   > cp env.sample test.env
   ```

Modify the values in `test.env` to match your needs.

Use helper bash script to run tests

```sh
# install gotestsum
> go install gotest.tools/gotestsum@latest
```

```sh
> sh ./scripts/run-local-tests.sh

# Alternatively, you can also run all the tests from the project root
> go test -v ./...

# Also, for running specific tests,
# use -run flag while running the test from the package path
> go test -v -run <test_name>
```

## Schema Migrations

### Applying Migrations

Once the local postgres database as been set up, run the following commands to apply the migration

```sh
> cd migrations

# Install the dependencies from the GemFile
> bundle install

# Load values from env config in console
```sh
loadenv()
{
  echo "Loading $1"
  for i in $(cat $1); do
    export $i
  done
}

loadenv development.env
```

# skip this command if the database has already been created
> bundle exec rake db:create

# Run the migrations
> bundle exec rake db:migrate

# Rollback the rollback
bundle exec rake db:rollback
```

Check out the documentation on [Active Record Migrations](https://guides.rubyonrails.org/active_record_migrations.html)
for more information on creating, running and changing migrations.

### Adding New Migrations

```sh
> cd migrations
> bundle install
> bundle exec rake db:migrate:new name=InitSetup
```

## Code formatting

To ensure lint/format rules are maintained, run the following on the codebase prior to committing files

```sh
# Fixing formatting issues
> go fmt ./...

# Lint checking using golangci-lint (https://github.com/golangci/golangci-lint)
> golangci-lint -D errcheck run
```

## Thumb Rules

To be added as and when needed.

## Assumptions

- Have created installments and their due date calculation on approval of loan rather than creation
- Remainder loan amount is distributed over all the last remaining loan installments instead of just last one. Generally in my opinion, initial installments should have the extra remainder amount. Can be tweaked as needed.
- one_time_settlements table exists in separate repayment system
- Cents is used as lower denomination. It can be tweaked as needed.
- Seed role as needed. Attached image for sample data.

## Pending things to add

- More test coverage
- Add user lock to avoid concurrency edge cases
- Add loan interest/penalty support
- Add some instrumentation library like newRelic
- Add events
- Along with roles, permissions can also be added
- CRON job to move installment status to payment_pending
- Dockerfile is not working, need more time to debug

## Apis

Refer postman collection for list of all apis

## Roles Seeded

![alt text](image.png)

## Flow to test apis

1. Make sure you create roles. I did not have bandwidth to expose apis to create role. For now, you can create role as following by directly running these SQL queries in loan_development database:

```sql
INSERT INTO roles(
 id, name, created_at, updated_at)
 VALUES ('rol_123', 'USER', now(), now()),
        ('rol_124', 'ADMIN', now(), now());
```

![alt text](image-1.png)

1. Once this is done, we are ready to test the apis. Lets start by calling create user api `POST /api/v1/users`. We need to create 2 users, one customer and one admin. You can use the following payload. This api will create user and its associated role.

customer

```json
{
    "name": "john",
    "password": "john@123",
    "role_id": "rol_123"
}
```

Admin user

```json
{
    "name": "admin",
    "password": "admin@123",
    "role_id": "rol_124"
}
```

1. Once user is created, we need to generate access_token which is used in subsequent apis to authenticate the user. You can call login api `POST /api/v1/users/login` to get user's access token. Please pass the name and password used while creating user here. You need to create access_token for both customer and admin user. Save both of them with you. You would need admin's access token for calling approve loan api.

customer login

```json
{
    "name": "john",
    "password": "john@123"
}
```

1. Once you get customer's access token, you can call create loan api to create loan `POST /api/v1/loans`. Make sure to update access token you received in the previous call in `Access-Token` header of this api. In the response, you would get loan.id. Save that for further api calls.

sample request

```json
{
    "amount_in_cents": 100,
    "term": 3,
    "frequency_in_days": 7
}
```

1. Once loan is created, you can call get loans api to list down loans of a user `GET /api/v1/loans`. Make sure to update access token you received for customer in `Access-Token` header of this api.

1. Now, that loan is created, admin needs to approve the loan. You can call approve loan api `POST /api/v1/loans/:id/approve` to get this done. Make sure to update access token you received for `admin` in `Access-Token` header of this api. If you use customer's access token, it would give you 401 unauthenticated as user doesn't have the necessary role. Make sure to update `:id` path param with the loan_id you received in create loan api. This api would result in all installment creation with their amount split and due_dates.

1. Now, once the loan is approved, customer can start paying in parts/full. You can call pay loan api `POST /api/v1/loans/:id/payment`.Make sure to update access token you received for customer in `Access-Token` header of this api. Make sure to update `:id` path param with the loan_id you received in create loan api. You can call this api multiple times for a loan until all tme amount has been paid. Pass any dummy one_time_settlement_id for now.

sample request

```json
{
    "one_time_settlement_id": "ots_123",
    "amount_in_cents": 46
}
```

1. Finally you can cleanup user's access token by calling logout api `DELETE /api/v1/users/logout`. Make sure to update access token you received for customer in `Access-Token` header of this api. Once this is done, user will not be able to call loan related apis. User has to login again to get new access token.
