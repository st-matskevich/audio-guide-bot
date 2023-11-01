BEGIN;

ALTER TABLE objects
    ADD title VARCHAR(64),
    ADD audio_path VARCHAR(128);

UPDATE objects
    SET title = objects_i18n.title, audio_path = objects_i18n.audio_path
    FROM objects_i18n
    WHERE objects_i18n.object_id = objects.object_id
    AND objects_i18n.language = 'en';

DROP TABLE objects_i18n;

END;