CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--
-- In order to make it possible to remove completed transcript chunks from the DB,
-- the author contributions need to be tracked separately.
--

CREATE TABLE "rank"
(
    id     TEXT PRIMARY KEY,
    name   TEXT UNIQUE,
    points DECIMAL
);

INSERT INTO rank (id, name, points)
VALUES (uuid_generate_v4(), 'Smelly Eyebrows', 1),
       (uuid_generate_v4(), 'Posh Student', 2),
       (uuid_generate_v4(), 'Camfield', 3),
       (uuid_generate_v4(), 'Lanzagrotty', 4),
       (uuid_generate_v4(), 'Snail Expert', 5),
       (uuid_generate_v4(), 'Cockroach Expert', 10),
       (uuid_generate_v4(), 'Fish Shop Owner' , 20),
       (uuid_generate_v4(), 'Wendy Robinson', 30),
       (uuid_generate_v4(), 'Saucer Drinker', 40),
       (uuid_generate_v4(), 'Cat Shaver', 50),
       (uuid_generate_v4(), 'Lanky Co-writer', 60),
       (uuid_generate_v4(), 'Producer', 70),
       (uuid_generate_v4(), 'Cheeky Freak', 80),
       (uuid_generate_v4(), 'Doh Nutter', 90),
       (uuid_generate_v4(), 'Rickydiculous', 100);
