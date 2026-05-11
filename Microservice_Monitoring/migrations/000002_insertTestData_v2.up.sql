-- Seed data voor init_v2 schema

INSERT INTO kunstwerktype (id, naam, beschrijving, deleted)
VALUES (1, 'Brug', 'Standaard verkeersbrug', FALSE);

INSERT INTO kunstwerk (id, beheeridentifier, naam, geolocation, kunstwerktype_id, beschrijving, last_send_dh_update, deleted)
VALUES (1, 'BRG-001', 'Brug bij Barneveld', NULL, 1, 'Testbrug voor lokaal ontwikkelen', NULL, FALSE);

INSERT INTO onderdelen (id, kunstwerk_id, naam, parent_id, deleted)
VALUES (1, 1, 'Brugdek Oost', NULL, FALSE);

INSERT INTO sensortype (id, naam, eenheid, drempel_is_range)
VALUES (1, 'Trillingssensor', 'mm/s', TRUE);

INSERT INTO sensor (id, kunstwerk_id, onderdeel_id, geolocation, sensortype_id, last_analyzed_meting_id, deleted)
VALUES (1, 1, 1, NULL, 1, NULL, FALSE);

INSERT INTO sensorconfiguratie (sensor_id, min_value, max_value, marge_percentage)
VALUES (1, 0.0, 10.0, 10.0);

INSERT INTO sensortype (id, naam, eenheid, drempel_is_range)
VALUES (2, 'Temperatuursensor', '°C', FALSE);

INSERT INTO sensor (id, kunstwerk_id, onderdeel_id, geolocation, sensortype_id, last_analyzed_meting_id, deleted)
VALUES (2, 1, 1, NULL, 2, NULL, FALSE);

INSERT INTO sensorconfiguratie (sensor_id, min_value, max_value, marge_percentage)
VALUES (2, 30.0, 0.0, 5.0);
