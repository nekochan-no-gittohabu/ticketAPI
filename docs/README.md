# TicketAPI

## SetUp

 * Create .env file that contains variables: 
    * DB_USER
    * DB_PASSWORD
    * DB_NAME
    * DB_CONTAINER
    * DB_URL
 * Start db: docker compose up -d
 * Create db tables: docker exec -it {container_id} migrate -path ./sql-migrations -database {database_url} up
    * db URL Format: postgres://{DB_USER}:{DB_PASSWORD}@{DB_CONTAINER}:5432/{DB_NAME}?sslmode=disable


 ## Testing

 * Unit tests for service are in ticketAPI/ticket/service_test.go
 * The API can be tested E2E with the requests in ticketAPI/api/ticket.rest. On VS Code, REST Client extension can be used.