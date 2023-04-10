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
