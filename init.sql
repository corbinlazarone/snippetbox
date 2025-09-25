-- Create a `snippets` table.
CREATE TABLE snippets (
    id SERIAL NOT NULL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created TIMESTAMP NOT NULL,
    expires TIMESTAMP NOT NULL
);

-- Add an index on the created column.
CREATE INDEX idx_snippets_created ON snippets(created);

-- Add some dummy records.
INSERT INTO snippets (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    CURRENT_TIMESTAMP AT TIME ZONE 'UTC',
    CURRENT_TIMESTAMP AT TIME ZONE 'UTC' + INTERVAL '365 days'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    CURRENT_TIMESTAMP AT TIME ZONE 'UTC',
    CURRENT_TIMESTAMP AT TIME ZONE 'UTC' + INTERVAL '365 days'
);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    CURRENT_TIMESTAMP AT TIME ZONE 'UTC',
    CURRENT_TIMESTAMP AT TIME ZONE 'UTC' + INTERVAL '7 days'
);
