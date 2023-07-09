CREATE TABLE IF NOT EXISTS compositions (
   id BIGSERIAL PRIMARY KEY,
   composer_id BIGINT NOT NULL,
   name TEXT NOT NULL,
   version BIGINT NOT NULL DEFAULT 1,
   CONSTRAINT fk_composer
      FOREIGN KEY(composer_id)
         REFERENCES composers(id)
);
