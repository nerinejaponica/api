CREATE TABLE artists (
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    slug        text NOT NULL UNIQUE,
    name_ja     text NOT NULL,
    name_romaji text
);

CREATE TABLE sources (
    id          bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    slug        text NOT NULL UNIQUE,
    name_ja     text NOT NULL,
    name_romaji text,
    kind        text NOT NULL DEFAULT 'lainnya'
                     CHECK (kind IN ('anime', 'game', 'album', 'lainnya'))
);

CREATE TABLE tags (
    id   bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    slug text NOT NULL UNIQUE,
    name text NOT NULL
);
