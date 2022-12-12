# TicketAPI

## SetUp

 * Start db: docker compose up -d
 * Create db tables
 * Host the API: go run main.go

 ## Testing

 * Unit tests for service are in ticketAPI/ticket/service_test.go
 * The API can be tested E2E with the requests in ticketAPI/api/ticket.rest. On VS Code, REST Client extension can be used.