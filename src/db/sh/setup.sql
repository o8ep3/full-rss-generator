-- importing uuid-ossp extension to be able to use uuid_generate_v4 function
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE FeedInfoTable
(
    id          UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
    title       TEXT,
    rss_url     TEXT UNIQUE NOT NULL,
    xpath       TEXT NOT NULL
);
