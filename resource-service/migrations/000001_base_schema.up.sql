CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE job_titles (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL UNIQUE
);

CREATE TABLE workgroups (
    id BIGSERIAL PRIMARY KEY,
    workgroup TEXT NOT NULL UNIQUE
);

CREATE TABLE resources (
    id INTEGER NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    workgroup_id INTEGER NOT NULL,
    job_title_id INTEGER NOT NULL,
    manager_id INTEGER NOT NULL,
    CONSTRAINT fk_workgroup FOREIGN KEY (workgroup_id) REFERENCES workgroups(id),
    CONSTRAINT fk_job_title FOREIGN KEY (job_title_id) REFERENCES job_titles(id)
);