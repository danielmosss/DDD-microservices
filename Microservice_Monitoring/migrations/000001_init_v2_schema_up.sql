CREATE TABLE kunstwerktype (
    id SERIAL PRIMARY KEY,
    naam VARCHAR(255) NOT NULL,
    beschrijving TEXT,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE kunstwerk (
    id BIGSERIAL PRIMARY KEY,
    beheeridentifier VARCHAR(255) UNIQUE NOT NULL,
    naam VARCHAR(255) NOT NULL,
    geolocation VARCHAR(255),
    kunstwerktype_id BIGINT REFERENCES kunstwerktype(id),
    beschrijving TEXT,
    last_send_dh_update TIMESTAMPTZ,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE onderdelen (
    id BIGSERIAL PRIMARY KEY,
    kunstwerk_id BIGINT REFERENCES kunstwerk(id) NOT NULL,
    naam VARCHAR(255) NOT NULL,
    parent_id BIGINT REFERENCES onderdelen(id),
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE sensortype (
    id SERIAL PRIMARY KEY,
    naam VARCHAR(255) NOT NULL,
    eenheid VARCHAR(50),
    drempel_is_range BOOLEAN NOT NULL
);

CREATE TABLE sensor (
    id BIGSERIAL PRIMARY KEY,
    kunstwerk_id BIGINT REFERENCES kunstwerk(id) NOT NULL,
    onderdeel_id BIGINT REFERENCES onderdelen(id),
    geolocation VARCHAR(255),
    sensortype_id INT REFERENCES sensortype(id) NOT NULL,
    last_analyzed_meting_id BIGINT REFERENCES meting(id) DEFAULT NULL,
    deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE sensorconfiguratie (
    id SERIAL PRIMARY KEY,
    sensor_id BIGINT REFERENCES sensor(id) NOT NULL UNIQUE,
    min_value FLOAT,
    max_value FLOAT,
    marge_percentage FLOAT
);

CREATE TABLE meting (
    id BIGSERIAL PRIMARY KEY,
    time TIMESTAMPTZ NOT NULL,
    sensor_id BIGINT REFERENCES sensor(id),
    kunstwerk_id BIGINT REFERENCES kunstwerk(id) NOT NULL,
    waarde FLOAT NOT NULL,
    is_handmatig BOOLEAN NOT NULL DEFAULT FALSE,
    inspectie_id VARCHAR(255)
    INDEX id_SensorId_idx (id, sensor_id) 
);

SELECT create_hypertable('meting', 'time');

CREATE TABLE afwijking (
    id BIGSERIAL PRIMARY KEY,
    meting_id BIGINT REFERENCES meting(id) NOT NULL,
    meting_time TIMESTAMPTZ NOT NULL,
    kunstwerk_id BIGINT REFERENCES kunstwerk(id) NOT NULL,
    sensor_id BIGINT REFERENCES sensor(id),
    time TIMESTAMPTZ NOT NULL,
    norm_min_waarde FLOAT NOT NULL,
    norm_max_waarde FLOAT,
    norm_marge_percentage FLOAT,
    gemeten_waarde FLOAT NOT NULL,
    is_warning BOOLEAN NOT NULL,
    FOREIGN KEY (meting_time, meting_id) REFERENCES meting(time, id)
);
