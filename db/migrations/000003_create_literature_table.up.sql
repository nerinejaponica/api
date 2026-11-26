CREATE TABLE literature (
    id              bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    slug            text        NOT NULL UNIQUE,
    title_ja        text        NOT NULL,
    title_id        text,
    original_author text,
    year            smallint,
    aozora_url      text,
    cover_url       text,
    body            text        NOT NULL DEFAULT '',
    footnotes       jsonb       NOT NULL DEFAULT '[]'::jsonb,
    status          text        NOT NULL DEFAULT 'draft'
                                CHECK (status IN ('draft', 'published')),
    published_at    timestamptz,
    updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX literature_published_at_idx ON literature (published_at DESC) WHERE status = 'published';
