DROP INDEX IF EXISTS songs_published_at_idx;

ALTER TABLE songs
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS footnotes,
    DROP COLUMN IF EXISTS lyrics,
    DROP COLUMN IF EXISTS media_links,
    DROP COLUMN IF EXISTS composer,
    DROP COLUMN IF EXISTS lyricist,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS album,
    DROP COLUMN IF EXISTS year,
    DROP COLUMN IF EXISTS cover_url,
    DROP COLUMN IF EXISTS title_id,
    DROP COLUMN IF EXISTS title_romaji;
