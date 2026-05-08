CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE kunstwerktype (
                               id SERIAL PRIMARY KEY,
                               naam VARCHAR(255) NOT NULL,
                               beschrijving TEXT
);

CREATE TABLE kunstwerk (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           naam VARCHAR(255) NOT NULL,
                           geolocation VARCHAR(255),
                           kunstwerktype_id INT REFERENCES kunstwerktype(id),
                           beschrijving TEXT,
                           deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE onderdelen (
                            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                            kunstwerk_id UUID REFERENCES kunstwerk(id) NOT NULL,
                            naam VARCHAR(255) NOT NULL,
                            parent_id UUID REFERENCES onderdelen(id)
);

CREATE TABLE sensortype (
                            id SERIAL PRIMARY KEY,
                            naam VARCHAR(255) NOT NULL,
                            eenheid VARCHAR(50),
                            drempel_is_range BOOLEAN NOT NULL
);

CREATE TABLE sensor (
                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                        kunstwerk_id UUID REFERENCES kunstwerk(id) NOT NULL,
                        onderdeel_id UUID REFERENCES onderdelen(id),
                        geolocation VARCHAR(255),
                        sensortype_id INT REFERENCES sensortype(id) NOT NULL
);

CREATE TABLE sensorconfiguratie (
                                    id SERIAL PRIMARY KEY,
                                    sensor_id UUID REFERENCES sensor(id) NOT NULL UNIQUE,
                                    min_value FLOAT,
                                    max_value FLOAT,
                                    marge_percentage FLOAT
);

CREATE TABLE meting (
                        time TIMESTAMPTZ NOT NULL,
                        id BIGSERIAL,
                        sensor_id UUID REFERENCES sensor(id),
                        kunstwerk_id UUID REFERENCES kunstwerk(id) NOT NULL,
                        waarde FLOAT NOT NULL,
                        is_afwijking BOOLEAN DEFAULT FALSE,
                        is_handmatig BOOLEAN DEFAULT FALSE,
                        inspectie_id INT,
                        PRIMARY KEY (time, id)
);

SELECT create_hypertable('meting', 'time');

CREATE TABLE afwijking (
                           id BIGSERIAL PRIMARY KEY,
                           meting_id BIGINT NOT NULL,
                           meting_time TIMESTAMPTZ NOT NULL,
                           kunstwerk_id UUID REFERENCES kunstwerk(id) NOT NULL,
                           time TIMESTAMPTZ NOT NULL,
                           norm_waarde FLOAT NOT NULL,
                           gemeten_waarde FLOAT NOT NULL,
                           is_warning BOOLEAN NOT NULL,
                           FOREIGN KEY (meting_time, meting_id) REFERENCES meting(time, id)
);