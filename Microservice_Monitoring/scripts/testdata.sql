-- ========== Scenario 1: Min-only sensor (no max_value) ==========
-- Configuration: min_value = 50, marge_percentage = 10
-- Expectation:
--  - waarde = 30  -> fatal (well below min - margin)
--  - waarde = 48  -> warning (within margin below min)
--  - waarde = 60  -> fatal (above min + margin)

INSERT INTO kunstwerk (id, beheeridentifier, naam) VALUES (10001, 'test-kw-10001', 'Test KW 10001') ON CONFLICT (id) DO NOTHING;
INSERT INTO sensortype (id, naam, eenheid, drempel_is_range) VALUES (10001, 'test-sensortype-1', 'u', false) ON CONFLICT (id) DO NOTHING;
INSERT INTO sensor (id, kunstwerk_id, sensortype_id, deleted) VALUES (10001, 10001, 10001, false) ON CONFLICT (id) DO NOTHING;
INSERT INTO sensorconfiguratie (sensor_id, min_value, max_value, marge_percentage) VALUES (10001, 50.0, NULL, 10.0)
    ON CONFLICT (sensor_id) DO UPDATE SET min_value = EXCLUDED.min_value, max_value = EXCLUDED.max_value, marge_percentage = EXCLUDED.marge_percentage;

INSERT INTO meting (time, sensor_id, kunstwerk_id, waarde) VALUES
(NOW() - INTERVAL '30 seconds', 10001, 10001, 30.0),
(NOW() - INTERVAL '20 seconds', 10001, 10001, 48.0),
(NOW() - INTERVAL '10 seconds', 10001, 10001, 60.0);

SELECT * FROM analyze_sensor_metingen() WHERE sensor_id = 10001;
SELECT meting_id, gemeten_waarde, is_warning, norm_min_waarde, norm_marge_percentage
FROM afwijking WHERE sensor_id = 10001 ORDER BY meting_time;
SELECT id, last_analyzed_meting_id FROM sensor WHERE id = 10001;

-- ========== Scenario 2: Range sensor with margin (min & max present) ==========
-- Configuration: min_value = 20, max_value = 30, marge_percentage = 10
-- Range length = 10 -> margin = 10% of 10 = 1
-- Expectation:
--  - waarde = 25   -> ok
--  - waarde = 19   -> warning (min - margin = 19)
--  - waarde = 30.5 -> warning (max + margin = 31)
--  - waarde = 40   -> fatal

INSERT INTO kunstwerk (id, beheeridentifier, naam) VALUES (10002, 'test-kw-10002', 'Test KW 10002') ON CONFLICT (id) DO NOTHING;
INSERT INTO sensortype (id, naam, eenheid, drempel_is_range) VALUES (10002, 'test-sensortype-2', 'u', true) ON CONFLICT (id) DO NOTHING;
INSERT INTO sensor (id, kunstwerk_id, sensortype_id, deleted) VALUES (10002, 10002, 10002, false) ON CONFLICT (id) DO NOTHING;
INSERT INTO sensorconfiguratie (sensor_id, min_value, max_value, marge_percentage) VALUES (10002, 20.0, 30.0, 10.0)
    ON CONFLICT (sensor_id) DO UPDATE SET min_value = EXCLUDED.min_value, max_value = EXCLUDED.max_value, marge_percentage = EXCLUDED.marge_percentage;

INSERT INTO meting (time, sensor_id, kunstwerk_id, waarde) VALUES
(NOW() - INTERVAL '35 seconds', 10002, 10002, 25.0),
(NOW() - INTERVAL '25 seconds', 10002, 10002, 19.0),
(NOW() - INTERVAL '15 seconds', 10002, 10002, 30.5),
(NOW() - INTERVAL '5 seconds', 10002, 10002, 40.0);

SELECT * FROM analyze_sensor_metingen() WHERE sensor_id = 10002;
SELECT meting_id, gemeten_waarde, is_warning, norm_min_waarde, norm_max_waarde, norm_marge_percentage
FROM afwijking WHERE sensor_id = 10002 ORDER BY meting_time;
SELECT id, last_analyzed_meting_id FROM sensor WHERE id = 10002;

-- ========== Scenario 3: Range sensor without margin (marge = NULL) ==========
-- Configuration: min_value = 0, max_value = 100, marge_percentage = NULL
-- Expectation:
--  - waarde = -10  -> fatal
--  - waarde = 50   -> ok

INSERT INTO kunstwerk (id, beheeridentifier, naam) VALUES (10003, 'test-kw-10003', 'Test KW 10003') ON CONFLICT (id) DO NOTHING;
INSERT INTO sensortype (id, naam, eenheid, drempel_is_range) VALUES (10003, 'test-sensortype-3', 'u', true) ON CONFLICT (id) DO NOTHING;
INSERT INTO sensor (id, kunstwerk_id, sensortype_id, deleted) VALUES (10003, 10003, 10003, false) ON CONFLICT (id) DO NOTHING;
INSERT INTO sensorconfiguratie (sensor_id, min_value, max_value, marge_percentage) VALUES (10003, 0.0, 100.0, NULL)
    ON CONFLICT (sensor_id) DO UPDATE SET min_value = EXCLUDED.min_value, max_value = EXCLUDED.max_value, marge_percentage = EXCLUDED.marge_percentage;

INSERT INTO meting (time, sensor_id, kunstwerk_id, waarde) VALUES
(NOW() - INTERVAL '40 seconds', 10003, 10003, -10.0),
(NOW() - INTERVAL '30 seconds', 10003, 10003, 50.0);

SELECT * FROM analyze_sensor_metingen() WHERE sensor_id = 10003;
SELECT meting_id, gemeten_waarde, is_warning, norm_min_waarde, norm_max_waarde, norm_marge_percentage
FROM afwijking WHERE sensor_id = 10003 ORDER BY meting_time;
SELECT id, last_analyzed_meting_id FROM sensor WHERE id = 10003;

-- ========== Scenario 4: Range sensitivity demonstration ==========
-- Configuration: min_value = 5, max_value = 100, marge_percentage = 20
-- Range length = 95 -> margin = 20% of 95 = 19
-- This shows why applying the percentage to min/max individually is wrong:
--  - old min-based margin: 5 * 20% = 1 (too small)
--  - new range-based margin: 19 (applied symmetrically)
-- Expectation:
--  - waarde = -10  -> warning (>= min - margin => -14)
--  - waarde = 4    -> warning (within margin below min)
--  - waarde = 118  -> warning (<= max + margin => 119)
--  - waarde = 130  -> fatal

INSERT INTO kunstwerk (id, beheeridentifier, naam) VALUES (10004, 'test-kw-10004', 'Test KW 10004') ON CONFLICT (id) DO NOTHING;
INSERT INTO sensortype (id, naam, eenheid, drempel_is_range) VALUES (10004, 'test-sensortype-4', 'u', true) ON CONFLICT (id) DO NOTHING;
INSERT INTO sensor (id, kunstwerk_id, sensortype_id, deleted) VALUES (10004, 10004, 10004, false) ON CONFLICT (id) DO NOTHING;
INSERT INTO sensorconfiguratie (sensor_id, min_value, max_value, marge_percentage) VALUES (10004, 5.0, 100.0, 20.0)
    ON CONFLICT (sensor_id) DO UPDATE SET min_value = EXCLUDED.min_value, max_value = EXCLUDED.max_value, marge_percentage = EXCLUDED.marge_percentage;

INSERT INTO meting (time, sensor_id, kunstwerk_id, waarde) VALUES
(NOW() - INTERVAL '50 seconds', 10004, 10004, -10.0),
(NOW() - INTERVAL '40 seconds', 10004, 10004, 4.0),
(NOW() - INTERVAL '30 seconds', 10004, 10004, 118.0),
(NOW() - INTERVAL '20 seconds', 10004, 10004, 130.0);

SELECT * FROM analyze_sensor_metingen() WHERE sensor_id = 10004;
SELECT meting_id, gemeten_waarde, is_warning, norm_min_waarde, norm_max_waarde, norm_marge_percentage
FROM afwijking WHERE sensor_id = 10004 ORDER BY meting_time;
SELECT id, last_analyzed_meting_id FROM sensor WHERE id = 10004;