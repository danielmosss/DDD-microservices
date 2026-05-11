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
    last_analyzed_meting_id BIGINT DEFAULT NULL,
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
    id BIGSERIAL NOT NULL,
    time TIMESTAMPTZ NOT NULL,
    sensor_id BIGINT REFERENCES sensor(id),
    kunstwerk_id BIGINT REFERENCES kunstwerk(id) NOT NULL,
    waarde FLOAT NOT NULL,
    is_handmatig BOOLEAN NOT NULL DEFAULT FALSE,
    inspectie_id VARCHAR(255)
);

SELECT create_hypertable('meting', 'time');
CREATE INDEX idx_meting_id ON meting(id);
CREATE INDEX idx_meting_time ON meting(time DESC);
CREATE INDEX idx_meting_sensor_id_id ON meting(sensor_id, id);

CREATE TABLE afwijking (
    id BIGSERIAL PRIMARY KEY,
    meting_id BIGINT NOT NULL,
    meting_time TIMESTAMPTZ NOT NULL,
    kunstwerk_id BIGINT REFERENCES kunstwerk(id) NOT NULL,
    sensor_id BIGINT REFERENCES sensor(id),
    time TIMESTAMPTZ NOT NULL,
    norm_min_waarde FLOAT NOT NULL,
    norm_max_waarde FLOAT,
    norm_marge_percentage FLOAT,
    gemeten_waarde FLOAT NOT NULL,
    is_warning BOOLEAN NOT NULL
);

CREATE TABLE procedure_error_log (
    id BIGSERIAL PRIMARY KEY,
    logged_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    procedure_name VARCHAR(255) NOT NULL,
    sensor_id BIGINT REFERENCES sensor(id),
    error_message TEXT NOT NULL,
    error_context JSONB
);

CREATE INDEX idx_error_log_procedure ON procedure_error_log(procedure_name);
CREATE INDEX idx_error_log_sensor ON procedure_error_log(sensor_id);