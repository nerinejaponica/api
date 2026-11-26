CREATE TABLE posts (
    id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    slug         text        NOT NULL UNIQUE,
    title        text        NOT NULL,
    excerpt      text,
    cover_url    text,
    body         text        NOT NULL DEFAULT '',
    status       text        NOT NULL DEFAULT 'draft'
                             CHECK (status IN ('draft', 'published')),
    published_at timestamptz,
    updated_at   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX posts_published_at_idx ON posts (published_at DESC) WHERE status = 'published';

ALTER TABLE content_tags
    ADD COLUMN post_id bigint REFERENCES posts (id) ON DELETE CASCADE,
    DROP CONSTRAINT content_tags_check,
    ADD CONSTRAINT content_tags_check CHECK (num_nonnulls(song_id, literature_id, post_id) = 1);

CREATE UNIQUE INDEX content_tags_post_uidx ON content_tags (tag_id, post_id) WHERE post_id IS NOT NULL;
CREATE INDEX content_tags_post_id_idx ON content_tags (post_id);
