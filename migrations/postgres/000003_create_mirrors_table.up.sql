CREATE TABLE IF NOT EXISTS composition_mirrors (
   id BIGSERIAL PRIMARY KEY,
   composition_id BIGINT NOT NULL,
   link TEXT NOT NULL,
   version BIGINT NOT NULL DEFAULT 1,
   CONSTRAINT fk_composition
      FOREIGN KEY(composition_id) 
         REFERENCES compositions(id)
);
