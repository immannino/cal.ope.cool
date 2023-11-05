CREATE TABLE IF NOT EXISTS calendar (
        id         INTEGER PRIMARY KEY AUTOINCREMENT,
        name_key   TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CREATE TABLE IF NOT EXISTS cache (
--         id         INTEGER PRIMARY KEY AUTOINCREMENT,
--         url        TEXT,
--         content    TEXT,
--         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- CREATE TABLE IF NOT EXISTS events (
--         id          INTEGER PRIMARY KEY AUTOINCREMENT,
--         google_id   TEXT NOT NULL,
--         title       TEXT NOT NULL,
--         description TEXT,
--         date        TEXT NOT NULL,
--         location    TEXT,
--         region      TEXT,
--         category    TEXT,
--         start       TIMESTAMP,
--         end         TIMESTAMP,
--         url         TEXT,
--         organizer   TEXT,
--         created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--         updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

--         UNIQUE(google_id)
-- );