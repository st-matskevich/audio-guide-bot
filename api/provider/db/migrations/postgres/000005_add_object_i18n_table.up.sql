BEGIN;

CREATE TABLE objects_i18n(
    i18n_id BIGSERIAL PRIMARY KEY,
    object_id BIGINT NOT NULL,
    language VARCHAR(2) NOT NULL,
    title VARCHAR(64) NOT NULL,
    audio_path VARCHAR(128) NOT NULL,
    UNIQUE (object_id, language));

INSERT INTO objects_i18n (object_id, language, title, audio_path)
    SELECT object_id, 'en', title, audio_path
    FROM objects;   

ALTER TABLE objects
    DROP COLUMN title,
    DROP COLUMN audio_path;

END;