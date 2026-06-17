CREATE TABLE IF NOT EXISTS projects (
    id          VARCHAR PRIMARY KEY,
    name        VARCHAR NOT NULL,
    owner_id    VARCHAR NOT NULL,
    created_at  TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_projects_owner ON projects(owner_id);
