#! bash
curl --location --request PUT 'http://localhost:8080/dashboard' \
--header 'Content-Type: application/json' \
--data-raw '{"id":1,"first_name":"Valaree","last_name":"Easom","email":"veasom0@buzzfeed.com","gender":"Female","ip_address":"89.162.180.56","animal":"Heloderma horridum"}' -v