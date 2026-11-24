CREATE TABLE songs (
    id           bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    slug         text        NOT NULL UNIQUE,
    title_ja     text        NOT NULL,
    status       text        NOT NULL DEFAULT 'draft'
                             CHECK (status IN ('draft', 'published')),
    published_at timestamptz
);
