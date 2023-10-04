BEGIN;

CREATE TABLE covers(
    cover_id BIGSERIAL PRIMARY KEY,
    object_id BIGINT NOT NULL,
    index INT NOT NULL,
    path VARCHAR(128) NOT NULL,
    UNIQUE (object_id, index));

INSERT INTO covers (object_id, index, path)
    SELECT object_id, 0, cover_path
    FROM objects;   

ALTER TABLE objects
    DROP COLUMN cover_path;

END;