CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    age INT,
    username VARCHAR(50) UNIQUE,
    email VARCHAR(100) UNIQUE,
    password VARCHAR(255),
    role VARCHAR(50)
);

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    category VARCHAR(100),
    description TEXT,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    is_paid BOOLEAN,
    price DECIMAL(10,2),
    capacity INT,
    latitude DECIMAL(9,6),
    longitude DECIMAL(9,6),
    poster_url VARCHAR(255),
    status VARCHAR(50)
);

CREATE TABLE event_attendees (
    user_id INT,
    event_id INT,
    rsvp_status VARCHAR(20),
    rsvp_date TIMESTAMP,
    PRIMARY KEY (user_id, event_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    event_id INT REFERENCES events(id),
    amount DECIMAL(10,2),
    transaction_date TIMESTAMP,
    status VARCHAR(50),
    payment_method VARCHAR(50),
    payment_gateway_transaction_id VARCHAR(100),
    notes TEXT
);