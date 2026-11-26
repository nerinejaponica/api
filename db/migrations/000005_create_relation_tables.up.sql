CREATE TABLE song_artists (
    song_id   bigint NOT NULL REFERENCES songs (id) ON DELETE CASCADE,
    artist_id bigint NOT NULL REFERENCES artists (id) ON DELETE RESTRICT,
    PRIMARY KEY (song_id, artist_id)
);
CREATE INDEX song_artists_artist_id_idx ON song_artists (artist_id);

CREATE TABLE song_sources (
    song_id   bigint NOT NULL REFERENCES songs (id) ON DELETE CASCADE,
    source_id bigint NOT NULL REFERENCES sources (id) ON DELETE RESTRICT,
    PRIMARY KEY (song_id, source_id)
);
CREATE INDEX song_sources_source_id_idx ON song_sources (source_id);

CREATE TABLE content_tags (
    id            bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    tag_id        bigint NOT NULL REFERENCES tags (id) ON DELETE RESTRICT,
    song_id       bigint REFERENCES songs (id) ON DELETE CASCADE,
    literature_id bigint REFERENCES literature (id) ON DELETE CASCADE,
    CONSTRAINT content_tags_check CHECK (num_nonnulls(song_id, literature_id) = 1)
);
CREATE UNIQUE INDEX content_tags_song_uidx ON content_tags (tag_id, song_id) WHERE song_id IS NOT NULL;
CREATE UNIQUE INDEX content_tags_literature_uidx ON content_tags (tag_id, literature_id) WHERE literature_id IS NOT NULL;
CREATE INDEX content_tags_song_id_idx ON content_tags (song_id);
CREATE INDEX content_tags_literature_id_idx ON content_tags (literature_id);
