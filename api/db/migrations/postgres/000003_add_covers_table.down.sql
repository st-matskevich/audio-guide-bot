BEGIN;

ALTER TABLE objects
    ADD cover_path VARCHAR(128);

UPDATE objects
    SET cover_path = covers.path
    FROM covers
    WHERE covers.object_id = objects.object_id
    AND covers.index = 0;

DROP TABLE covers;

END;