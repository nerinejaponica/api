ALTER TABLE songs
    ADD COLUMN title_romaji text,
    ADD COLUMN title_id     text,
    ADD COLUMN cover_url    text,
    ADD COLUMN year         smallint,
    ADD COLUMN album        text,
    ADD COLUMN description  text,
    ADD COLUMN lyricist     text,
    ADD COLUMN composer     text,
    ADD COLUMN media_links  jsonb       NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN lyrics       jsonb       NOT NULL DEFAULT '{"schemaVersion": 1, "stanzas": []}'::jsonb,
    ADD COLUMN footnotes    jsonb       NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN updated_at   timestamptz NOT NULL DEFAULT now();

CREATE INDEX songs_published_at_idx ON songs (published_at DESC) WHERE status = 'published';
