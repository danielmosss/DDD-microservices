INSERT INTO kunstwerk (beheeridentifier, naam, geolocation, kunstwerktype_id, beschrijving, deleted) VALUES
                                                                                                         ('KW-2026-BRUG-01', 'Erasmus-Stijl Kabelbrug', '51.9225,4.4792', (SELECT id FROM kunstwerktype WHERE naam = 'Brug'), 'Grote tuibrug met actieve monitoring', false),
                                                                                                         ('KW-2026-TUNN-02', 'Spoorvlakte Tunnel', '52.0907,5.1214', (SELECT id FROM kunstwerktype WHERE naam = 'Tunnel'), 'Ondergrondse spoortunnel met diepe compartimentering', false),
                                                                                                         ('KW-2026-SLUI-03', 'Delta Sluis West', '53.2109,6.5679', (SELECT id FROM kunstwerktype WHERE naam = 'Sluis'), 'Geautomatiseerd sluizencomplex met waterkering', false);

-- 3. ONDERDELEN (Met subqueries voor parent_id en kunstwerk_id)
-- Niveau 1: Hoofdsecties (parent_id = NULL)
INSERT INTO onderdelen (kunstwerk_id, naam, parent_id, deleted) VALUES
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), 'Bovenbouw (Superstructure)', NULL, false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), 'Onderbouw (Substructure)', NULL, false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-TUNN-02'), 'Tunnelbuis Noord', NULL, false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-TUNN-02'), 'Technische Ruimtes', NULL, false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-SLUI-03'), 'Sluisdeur Complex', NULL, false);

-- Niveau 2: Subsystemen
INSERT INTO onderdelen (kunstwerk_id, naam, parent_id, deleted) VALUES
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), 'Brugdek Compartimenten', (SELECT id FROM onderdelen WHERE naam = 'Bovenbouw (Superstructure)'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), 'Tuikabel Systeem', (SELECT id FROM onderdelen WHERE naam = 'Bovenbouw (Superstructure)'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), 'Pylonen (Hoofdpilaar)', (SELECT id FROM onderdelen WHERE naam = 'Onderbouw (Substructure)'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-TUNN-02'), 'Rijvloer Segmenten', (SELECT id FROM onderdelen WHERE naam = 'Tunnelbuis Noord'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-TUNN-02'), 'Ventilatiesysteem Axiaal', (SELECT id FROM onderdelen WHERE naam = 'Tunnelbuis Noord'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-SLUI-03'), 'Hydraulisch Aandrijvingssysteem', (SELECT id FROM onderdelen WHERE naam = 'Sluisdeur Complex'), false);

-- Niveau 3: Elementen
INSERT INTO onderdelen (kunstwerk_id, naam, parent_id, deleted) VALUES
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), 'Stalen Kokerliggers SE-01', (SELECT id FROM onderdelen WHERE naam = 'Brugdek Compartimenten'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), 'Kabelverankering Top Pylon', (SELECT id FROM onderdelen WHERE naam = 'Tuikabel Systeem'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), 'Fundering Pylon Voet', (SELECT id FROM onderdelen WHERE naam = 'Pylonen (Hoofdpilaar)'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-TUNN-02'), 'Segment 12 Dilatatievoeg', (SELECT id FROM onderdelen WHERE naam = 'Rijvloer Segmenten'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-TUNN-02'), 'Jet Fan Behuizing JF-04', (SELECT id FROM onderdelen WHERE naam = 'Ventilatiesysteem Axiaal'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-SLUI-03'), 'Hydraulische Cilinder HC-200', (SELECT id FROM onderdelen WHERE naam = 'Hydraulisch Aandrijvingssysteem'), false);

-- Niveau 4: Bouwdelen
INSERT INTO onderdelen (kunstwerk_id, naam, parent_id, deleted) VALUES
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), 'Slijtlaag Asfaltbedding', (SELECT id FROM onderdelen WHERE naam = 'Stalen Kokerliggers SE-01'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), 'Dwarsdragers Intern', (SELECT id FROM onderdelen WHERE naam = 'Stalen Kokerliggers SE-01'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), 'Betonwapening Matrix', (SELECT id FROM onderdelen WHERE naam = 'Fundering Pylon Voet'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), 'Diepwand Palenrij', (SELECT id FROM onderdelen WHERE naam = 'Fundering Pylon Voet'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-TUNN-02'), 'Rubber Glijprofiel', (SELECT id FROM onderdelen WHERE naam = 'Segment 12 Dilatatievoeg'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-TUNN-02'), 'Ophangbeugels & Trillingsdempers', (SELECT id FROM onderdelen WHERE naam = 'Jet Fan Behuizing JF-04'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-SLUI-03'), 'Zuigerstang Seal-behuizing', (SELECT id FROM onderdelen WHERE naam = 'Hydraulische Cilinder HC-200'), false);

-- Niveau 5: Componenten
INSERT INTO onderdelen (kunstwerk_id, naam, parent_id, deleted) VALUES
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), 'Lasnaad Verbinding L-234', (SELECT id FROM onderdelen WHERE naam = 'Dwarsdragers Intern'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), 'Boorpaal Kop BP-08', (SELECT id FROM onderdelen WHERE naam = 'Diepwand Palenrij'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-TUNN-02'), 'Ankerbout Bevestiging M24', (SELECT id FROM onderdelen WHERE naam = 'Ophangbeugels & Trillingsdempers'), false),
                                                                    ((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-SLUI-03'), 'Primaire Drukafdichting Ring', (SELECT id FROM onderdelen WHERE naam = 'Zuigerstang Seal-behuizing'), false);

-- 4. SENSORTYPEN
INSERT INTO sensortype (naam, eenheid, drempel_is_range) VALUES
                                                             ('Spanning (Strain Gauge)', 'µε', true),
                                                             ('Inclinometer (Helling)', 'mrad', true),
                                                             ('Druk', 'bar', true),
                                                             ('Akoestische Emissie', 'dB', true);

-- 5. SENSORS
INSERT INTO sensor (kunstwerk_id, onderdeel_id, geolocation, sensortype_id, deleted) VALUES
-- Sensoren Brug
((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), (SELECT id FROM onderdelen WHERE naam = 'Lasnaad Verbinding L-234'), '51.9225,4.4792', (SELECT id FROM sensortype WHERE naam = 'Spanning (Strain Gauge)'), false),
((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), (SELECT id FROM onderdelen WHERE naam = 'Slijtlaag Asfaltbedding'), '51.9225,4.4795', (SELECT id FROM sensortype WHERE naam = 'Temperatuur'), false),
((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), (SELECT id FROM onderdelen WHERE naam = 'Kabelverankering Top Pylon'), '51.9220,4.4790', (SELECT id FROM sensortype WHERE naam = 'Inclinometer (Helling)'), false),
((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), (SELECT id FROM onderdelen WHERE naam = 'Boorpaal Kop BP-08'), '51.9222,4.4788', (SELECT id FROM sensortype WHERE naam = 'Accelerometer'), false),
((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-BRUG-01'), (SELECT id FROM onderdelen WHERE naam = 'Bovenbouw (Superstructure)'), '51.9225,4.4792', (SELECT id FROM sensortype WHERE naam = 'Accelerometer'), false),

-- Sensoren Tunnel
((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-TUNN-02'), (SELECT id FROM onderdelen WHERE naam = 'Segment 12 Dilatatievoeg'), '52.0907,5.1214', (SELECT id FROM sensortype WHERE naam = 'Spanning (Strain Gauge)'), false),
((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-TUNN-02'), (SELECT id FROM onderdelen WHERE naam = 'Ankerbout Bevestiging M24'), '52.0909,5.1216', (SELECT id FROM sensortype WHERE naam = 'Accelerometer'), false),
((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-TUNN-02'), (SELECT id FROM onderdelen WHERE naam = 'Ventilatiesysteem Axiaal'), '52.0907,5.1214', (SELECT id FROM sensortype WHERE naam = 'Temperatuur'), false),

-- Sensoren Sluis
((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-SLUI-03'), (SELECT id FROM onderdelen WHERE naam = 'Hydraulische Cilinder HC-200'), '53.2109,6.5679', (SELECT id FROM sensortype WHERE naam = 'Druk'), false),
((SELECT id FROM kunstwerk WHERE beheeridentifier = 'KW-2026-SLUI-03'), (SELECT id FROM onderdelen WHERE naam = 'Primaire Drukafdichting Ring'), '53.2110,6.5682', (SELECT id FROM sensortype WHERE naam = 'Akoestische Emissie'), false);

-- 6. SENSORCONFIGURATIES
INSERT INTO sensorconfiguratie (sensor_id, min_value, max_value, marge_percentage) VALUES
                                                                                       ((SELECT id FROM sensor WHERE onderdeel_id = (SELECT id FROM onderdelen WHERE naam = 'Lasnaad Verbinding L-234') AND sensortype_id = (SELECT id FROM sensortype WHERE naam = 'Spanning (Strain Gauge)')), -500, 500, 5),
                                                                                       ((SELECT id FROM sensor WHERE onderdeel_id = (SELECT id FROM onderdelen WHERE naam = 'Slijtlaag Asfaltbedding') AND sensortype_id = (SELECT id FROM sensortype WHERE naam = 'Temperatuur')), -15, 60, 10),
                                                                                       ((SELECT id FROM sensor WHERE onderdeel_id = (SELECT id FROM onderdelen WHERE naam = 'Kabelverankering Top Pylon') AND sensortype_id = (SELECT id FROM sensortype WHERE naam = 'Inclinometer (Helling)')), -5, 5, 2),
                                                                                       ((SELECT id FROM sensor WHERE onderdeel_id = (SELECT id FROM onderdelen WHERE naam = 'Boorpaal Kop BP-08') AND sensortype_id = (SELECT id FROM sensortype WHERE naam = 'Accelerometer')), 0, 150, 8),
                                                                                       ((SELECT id FROM sensor WHERE onderdeel_id = (SELECT id FROM onderdelen WHERE naam = 'Bovenbouw (Superstructure)') AND sensortype_id = (SELECT id FROM sensortype WHERE naam = 'Accelerometer')), 0, 50, 5),
                                                                                       ((SELECT id FROM sensor WHERE onderdeel_id = (SELECT id FROM onderdelen WHERE naam = 'Segment 12 Dilatatievoeg') AND sensortype_id = (SELECT id FROM sensortype WHERE naam = 'Spanning (Strain Gauge)')), 0, 20, 10),
                                                                                       ((SELECT id FROM sensor WHERE onderdeel_id = (SELECT id FROM onderdelen WHERE naam = 'Ankerbout Bevestiging M24') AND sensortype_id = (SELECT id FROM sensortype WHERE naam = 'Accelerometer')), 0, 200, 12),
                                                                                       ((SELECT id FROM sensor WHERE onderdeel_id = (SELECT id FROM onderdelen WHERE naam = 'Ventilatiesysteem Axiaal') AND sensortype_id = (SELECT id FROM sensortype WHERE naam = 'Temperatuur')), 0, 45, 5),
                                                                                       ((SELECT id FROM sensor WHERE onderdeel_id = (SELECT id FROM onderdelen WHERE naam = 'Hydraulische Cilinder HC-200') AND sensortype_id = (SELECT id FROM sensortype WHERE naam = 'Druk')), 0, 320, 4),
                                                                                       ((SELECT id FROM sensor WHERE onderdeel_id = (SELECT id FROM onderdelen WHERE naam = 'Primaire Drukafdichting Ring') AND sensortype_id = (SELECT id FROM sensortype WHERE naam = 'Akoestische Emissie')), 0, 40, 15);