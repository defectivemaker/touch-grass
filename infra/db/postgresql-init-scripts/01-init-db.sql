CREATE EXTENSION IF NOT EXISTS postgis;

DROP TABLE IF EXISTS entries;

CREATE TABLE entries (
  id             SERIAL PRIMARY KEY,
  deviceUUID     VARCHAR(36) NOT NULL,
  payphoneID     VARCHAR(40) NOT NULL,
  payphoneMAC    VARCHAR(17) NOT NULL,
  payphoneTime   INT NOT NULL,
  recordedTime   TIMESTAMP NOT NULL,
  mapUUID        VARCHAR(40),
  mapLatitude    DOUBLE PRECISION,
  mapLongitude   DOUBLE PRECISION,
  mapLocation    geography(POINT, 4326)
);

DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  pub_key BYTEA,
  email VARCHAR(255),
  uuid VARCHAR(36),
  token VARCHAR(20),
  username VARCHAR(255),
  created_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  hasDevice BOOLEAN DEFAULT FALSE,
  UNIQUE (uuid)
);

DROP TABLE IF EXISTS user_statistics;
CREATE TABLE user_statistics (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    total_payphones INT DEFAULT 0,
    total_entries INT DEFAULT 0,
    total_maps INT DEFAULT 0,
    payphone_rank INT DEFAULT 0,
    entry_rank INT DEFAULT 0,
    map_rank INT DEFAULT 0,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (uuid),
    UNIQUE (user_id)
);

-- Create an index on the user_id column for faster lookups
CREATE INDEX idx_user_statistics_user_id ON user_statistics (user_id);

DROP TABLE IF EXISTS telstra_hotspots;

CREATE TABLE telstra_hotspots (
  uuid VARCHAR(36) PRIMARY KEY,
  location geography(POINT, 4326) NOT NULL,
  street_address VARCHAR(255),
  alias VARCHAR(255)
);

CREATE INDEX idx_location ON telstra_hotspots USING gist (location);

CREATE TEMPORARY TABLE tmp_telstra_hotspots (
  latitude DOUBLE PRECISION,
  longitude DOUBLE PRECISION,
  street_address VARCHAR(255),
  uuid VARCHAR(36),
  alias VARCHAR(255)
);

BEGIN;

\copy tmp_telstra_hotspots(latitude, longitude, street_address, uuid, alias) FROM '/var/lib/postgres-files/telstra_hotspots.csv' WITH (FORMAT csv, DELIMITER '|', HEADER true);

INSERT INTO telstra_hotspots (uuid, location, street_address, alias)
SELECT 
  uuid,
  ST_Point(longitude, latitude),
  street_address, 
  alias
FROM tmp_telstra_hotspots;

DROP TABLE tmp_telstra_hotspots;

COMMIT;
