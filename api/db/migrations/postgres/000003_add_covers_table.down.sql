BEGIN;

ALTER TABLE objects
    ADD cover_path VARCHAR(128);

UPDATE objects
    SET cover_path = covers.path
    FROM covers
    WHERE objects.object_id = covers.object_id;

DROP TABLE covers;

END;