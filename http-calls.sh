# Create user1 and top-up for 100.11
curl --location --request POST 'http://localhost:8080/api/balance/add' \
--header 'Content-Type: application/json' \
--data-raw '{
  "account_id": "user1",
  "amount": 100.11
}'

# Get user1 balance
curl --location --request GET 'http://localhost:8080/api/balance/get/user1'

# Transfer from user1 to created user2
curl --location --request POST 'http://localhost:8080/api/balance/transfer' \
--header 'Content-Type: application/json' \
--data-raw '{
  "sender_id": "user1",
  "receiver_id": "user2",
  "amount": 10.11
}'

# New order reserve for user1
curl --location --request POST 'http://localhost:8080/api/order/reserve' \
--header 'Content-Type: application/json' \
--data-raw '{
  "target_id": "user1",
  "service_id": "Cleaning",
  "order_id": 12345678903,
  "amount": 50.99
}'

# New order reserve for user2
curl --location --request POST 'http://localhost:8080/api/order/reserve' \
--header 'Content-Type: application/json' \
--data-raw '{
  "target_id": "user2",
  "service_id": "Cleaning",
  "order_id": 1134780699,
  "amount": 10.03
}'

# Approve user1 order
curl --location --request POST 'http://localhost:8080/api/order/withdraw' \
--header 'Content-Type: text/plain' \
--data-raw '{
  "order_id": 12345678903
}'

# Cancel user2 order
curl --location --request POST 'http://localhost:8080/api/order/cancel' \
--header 'Content-Type: text/plain' \
--data-raw '{
  "order_id": 1134780699
}'

# Create new report
curl --location --request POST 'http://localhost:8080/api/accounting/add' \
--header 'Content-Type: text/plain' \
--data-raw '{
  "month": "11",
  "year": "2022"
}'

# Get created report
curl --location --request GET 'http://localhost:8080/api/accounting/get/163bd0a74e'

# Get user1 history
curl --location --request POST 'http://localhost:8080/api/history' \
--header 'Content-Type: application/json' \
--data-raw '{
  "userID": "user1",
  "sortBy": "date",
  "sortOrder": "desc",
  "quantity": 10
}
'