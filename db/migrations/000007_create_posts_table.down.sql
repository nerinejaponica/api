DELETE FROM content_tags WHERE post_id IS NOT NULL;

ALTER TABLE content_tags
    DROP CONSTRAINT content_tags_check,
    DROP COLUMN post_id,
    ADD CONSTRAINT content_tags_check CHECK (num_nonnulls(song_id, literature_id) = 1);

DROP TABLE IF EXISTS posts;
