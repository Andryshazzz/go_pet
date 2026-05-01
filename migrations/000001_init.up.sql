CREATE SCHEMA myapp;

CREATE TABLE myapp.users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(15)
);