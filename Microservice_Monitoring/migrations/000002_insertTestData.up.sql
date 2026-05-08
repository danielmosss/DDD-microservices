-- 1. Maak een type kunstwerk aan
INSERT INTO kunstwerktype (id, naam, beschrijving)
VALUES (1, 'Brug', 'Standaard verkeersbrug');

-- 2. Maak een Kunstwerk aan met een vaste UUID
INSERT INTO kunstwerk (id, BeheerIdentifier, naam, kunstwerktype_id, beschrijving)
VALUES (1, 'BRG-001', 'Brug bij Barneveld', 1, 'Testbrug voor lokaal ontwikkelen');

-- 3. Voeg een onderdeel toe aan de brug
INSERT INTO onderdelen (id, kunstwerk_id, naam)
VALUES (1, 1, 'Brugdek Oost');

-- 4. Maak een sensortype aan (bijv. een trillingssensor)
INSERT INTO sensortype (id, naam, eenheid, drempel_is_range)
VALUES (1, 'Trillingssensor', 'mm/s', true);

-- 5. Koppel de sensor aan het onderdeel en het kunstwerk met een vaste UUID
INSERT INTO sensor (id, kunstwerk_id, onderdeel_id, sensortype_id)
VALUES (1, 1, 1, 1);

-- 6. Stel de configuratie (bedrijfsregels) in voor deze specifieke sensor
-- Hier zeggen we: normaal is tussen 0 en 10. Marge is 10%. (Dus 10 tot 11 is een warning, daarboven een harde afwijking).
INSERT INTO sensorconfiguratie (sensor_id, min_value, max_value, marge_percentage)
VALUES (1, 0.0, 10.0, 10.0);