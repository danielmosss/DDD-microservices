-- Insert bulk test data for comprehensive testing

-- Kunstwerktype
INSERT INTO kunstwerktype (naam, beschrijving, deleted) VALUES
    ('Brug', 'Voetgangers en fietsersbrug', false),
    ('Viaduct', 'Snelwegviaduct', false),
    ('Tunnel', 'Verkeerstunnel', false),
    ('Sluis', 'Waterkering/sluis', false);

-- Kunstwerken (structures being monitored)
INSERT INTO kunstwerk (beheeridentifier, naam, geolocation, kunstwerktype_id, beschrijving, deleted) VALUES
    ('KW-2024-001', 'Amsterdam Brug Noord', '52.3676,4.9041', 1, 'Nieuwe brug over de Amstel', false),
    ('KW-2024-002', 'Rotterdam Viaduct A15', '51.9225,4.4792', 2, 'Viaduct op snelweg A15', false),
    ('KW-2024-003', 'Utrecht Tunnel Centraal', '52.0907,5.1214', 3, 'Tunnel onder centraal station', false),
    ('KW-2024-004', 'Groningen Sluis West', '53.2109,6.5679', 4, 'Industriële sluis', false),
    ('KW-2024-005', 'Den Haag Brug Oud', '52.0705,4.3007', 1, 'Historische brug', false);

-- Onderdelen (structural components)
INSERT INTO onderdelen (kunstwerk_id, naam, parent_id, deleted) VALUES
    (1, 'Hoofdspanten', NULL, false),
    (1, 'Opleggingen', NULL, false),
    (1, 'Dek', NULL, false),
    (2, 'Pijlers', NULL, false),
    (2, 'Dek', NULL, false),
    (2, 'Expansieverbindingen', NULL, false),
    (3, 'Wandconstructie', NULL, false),
    (3, 'Ventilatie', NULL, false),
    (4, 'Palen', NULL, false),
    (4, 'Poorten', NULL, false),
    (5, 'Horizontale balken', NULL, false),
    (5, 'Verticale pilaren', NULL, false);

-- Sensortypes
INSERT INTO sensortype (naam, eenheid, drempel_is_range) VALUES
    ('Accelerometer', 'mm/s²', true),
    ('Temperatuur', '°C', true),
    ('Vochtigheid', '%', true),
    ('Spanning', 'kN', true),
    ('Doorbuiging', 'mm', true),
    ('Rotatie', 'mrad', true);

-- Sensors (many per kunstwerk with variety)
INSERT INTO sensor (kunstwerk_id, onderdeel_id, geolocation, sensortype_id, deleted) VALUES
    -- Amsterdam Brug Noord (Kunstwerk 1)
    (1, 1, '52.3676,4.9041', 1, false), -- Accelerometer on main span
    (1, 1, '52.3676,4.9041', 5, false), -- Deflection on main span
    (1, 2, '52.3676,4.9041', 4, false), -- Tension in bearings
    (1, 3, '52.3676,4.9041', 2, false), -- Temperature on deck
    
    -- Rotterdam Viaduct A15 (Kunstwerk 2)
    (2, 4, '51.9225,4.4792', 1, false), -- Accelerometer on pillar
    (2, 4, '51.9225,4.4792', 5, false), -- Settlement monitoring
    (2, 5, '51.9225,4.4792', 2, false), -- Temperature on deck
    (2, 5, '51.9225,4.4792', 3, false), -- Humidity on deck
    (2, 6, '51.9225,4.4792', 4, false), -- Expansion joint stress
    
    -- Utrecht Tunnel (Kunstwerk 3)
    (3, 7, '52.0907,5.1214', 2, false), -- Temperature in tunnel
    (3, 7, '52.0907,5.1214', 3, false), -- Humidity in tunnel
    (3, 8, '52.0907,5.1214', 1, false), -- Vibration from traffic
    (3, 8, '52.0907,5.1214', 2, false), -- Ventilation temp
    
    -- Groningen Sluis (Kunstwerk 4)
    (4, 9, '53.2109,6.5679', 5, false), -- Pile settlement
    (4, 9, '53.2109,6.5679', 4, false), -- Pile stress
    (4, 10, '53.2109,6.5679', 1, false), -- Door vibration
    
    -- Den Haag Brug Oud (Kunstwerk 5)
    (5, 11, '52.0705,4.3007', 1, false), -- Vibration monitoring
    (5, 11, '52.0705,4.3007', 5, false), -- Deflection
    (5, 12, '52.0705,4.3007', 4, false), -- Compression in pillar
    (5, 12, '52.0705,4.3007', 2, false); -- Temperature

-- Sensor configurations (thresholds and margins)
INSERT INTO sensorconfiguratie (sensor_id, min_value, max_value, marge_percentage) VALUES
    -- Accelerometer sensors
    (1, 0, 100, 10),
    (5, 0, 100, 10),
    (12, 0, 100, 10),
    (16, 0, 100, 10),
    (17, 0, 100, 10),

    -- Deflection sensors
    (2, 0, 50, 15),
    (6, 0, 50, 15),
    (14, 0, 50, 15),
    (18, 0, 50, 15),

    -- Tension sensors
    (3, 100, 500, 8),
    (9, 100, 500, 8),
    (15, 100, 500, 8),
    (19, 100, 500, 8),

    -- Temperature sensors
    (4, 5, 35, 20),
    (7, 5, 35, 20),
    (10, 5, 35, 20),
    (13, 5, 35, 20),
    (20, 5, 35, 20),

    -- Humidity sensors
    (8, 30, 80, 5),
    (11, 30, 80, 5);
