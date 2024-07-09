DROP DATABASE IF EXISTS mediascraper;
CREATE DATABASE mediascraper;

\c mediascraper

CREATE TABLE IF NOT EXISTS url_metadata (
    id SERIAL PRIMARY KEY,
    real_url TEXT NOT NULL,
    stored_url TEXT NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);