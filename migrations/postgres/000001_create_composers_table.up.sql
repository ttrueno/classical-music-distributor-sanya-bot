CREATE TABLE IF NOT EXISTS composers (
   id BIGSERIAL PRIMARY KEY,
   first_name TEXT NOT NULL,
   last_name TEXT NOT NULL,
   image_link TEXT,
   description TEXT,
   version BIGINT NOT NULL DEFAULT 1
);
