CREATE TABLE objects(
    object_id BIGSERIAL PRIMARY KEY, 
    code VARCHAR(64) NOT NULL UNIQUE,
    title VARCHAR(64) NOT NULL,
    cover_path VARCHAR(128) NOT NULL,
    audio_path VARCHAR(128) NOT NULL);