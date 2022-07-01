CREATE TABLE IF NOT EXISTS urls (
    id integer NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    uid varchar(36),
    short_url varchar(40),
    original_url varchar(2000),
    created_at timestamp DEFAULT current_timestamp
);
CREATE INDEX IF NOT EXISTS original_url_index ON urls USING btree (original_url);