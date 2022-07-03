CREATE TABLE IF NOT EXISTS urls (
    id integer NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    uid varchar(36) NOT NULL,
    short_url_id varchar(8) UNIQUE,
    original_url varchar(2000) NOT NULL,
    created_at timestamp DEFAULT current_timestamp
)