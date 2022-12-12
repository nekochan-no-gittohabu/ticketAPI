CREATE TABLE tickets (
    ticket_id SERIAL NOT NULL PRIMARY KEY,
    name varchar,
    descr varchar,
    allocation int CHECK (allocation >= 0),
);