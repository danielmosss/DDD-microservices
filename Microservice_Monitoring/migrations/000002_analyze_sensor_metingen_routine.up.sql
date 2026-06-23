-- Stored Procedure: analyze_sensor_metingen
-- Purpose: Analyze all unanalyzed measurements per sensor and detect anomalies
-- Returns: COUNT of detected anomalies

CREATE OR REPLACE FUNCTION analyze_sensor_metingen()
RETURNS TABLE (
    sensor_id BIGINT,
    kunstwerk_id BIGINT,
    metingen_processed BIGINT,
    afwijkingen_detected BIGINT,
    last_meting_id BIGINT,
    status VARCHAR
) AS $$
DECLARE
    v_sensor_id BIGINT;
    v_kunstwerk_id BIGINT;
    v_last_analyzed BIGINT;
    v_config RECORD;
    v_meting RECORD;
    v_metingen_count BIGINT := 0;
    v_afwijkingen_count BIGINT := 0;
    v_last_meting_id BIGINT := NULL;
    v_status VARCHAR := 'success';
    v_error_msg TEXT;
    v_status_check VARCHAR;
    v_is_warning BOOLEAN;
    v_norm_min FLOAT;
    v_norm_max FLOAT;
    v_marge FLOAT;
    v_range_margin FLOAT;

BEGIN
    -- Loop through all active sensors
    FOR v_sensor_id, v_kunstwerk_id, v_last_analyzed IN
        SELECT s.id, s.kunstwerk_id, COALESCE(s.last_analyzed_meting_id, 0)
        FROM sensor s
        WHERE s.deleted = FALSE
          AND EXISTS (
              SELECT 1
              FROM sensorconfiguratie sc
              WHERE sc.sensor_id = s.id
          )
    LOOP
        BEGIN
            -- Get sensor configuration
            SELECT sc.min_value, sc.max_value, sc.marge_percentage
            INTO v_config
            FROM sensorconfiguratie sc
            WHERE sc.sensor_id = v_sensor_id;

            IF v_config IS NULL THEN
                RAISE EXCEPTION 'No configuration found for sensor %', v_sensor_id;
            END IF;

            v_norm_min := v_config.min_value;
            v_norm_max := v_config.max_value;
            v_marge := v_config.marge_percentage;

            -- Process all metingen after last_analyzed_meting_id
            FOR v_meting IN
                SELECT m.id, m.time, m.waarde
                FROM meting m
                WHERE m.sensor_id = v_sensor_id
                  AND m.id > v_last_analyzed
                                ORDER BY m.sensor_id ASC, m.id ASC
            LOOP
                v_metingen_count := v_metingen_count + 1;
                v_last_meting_id := v_meting.id;

                -- Check thresholds using the same logic as the Go service
                v_status_check := 'oke';
                v_is_warning := FALSE;

                -- S1: Simple threshold (min_value only, no range)
                IF v_norm_max IS NULL OR v_norm_max = 0 THEN
                    IF v_marge IS NULL OR v_marge = 0 THEN
                        IF v_meting.waarde > v_norm_min THEN
                            v_status_check := 'fatal';
                        ELSEIF v_meting.waarde < v_norm_min THEN
                            v_status_check := 'fatal';
                        END IF;
                    ELSE
                        -- With margin
                        IF v_meting.waarde > (v_norm_min + (v_norm_min * v_marge / 100.0)) OR
                           v_meting.waarde < (v_norm_min - (v_norm_min * v_marge / 100.0)) THEN
                            v_status_check := 'fatal';
                        ELSEIF v_meting.waarde < (v_norm_min + (v_norm_min * v_marge / 100.0)) AND
                               v_meting.waarde > v_norm_min THEN
                            v_status_check := 'warning';
                        ELSEIF v_meting.waarde > (v_norm_min - (v_norm_min * v_marge / 100.0)) AND
                               v_meting.waarde < v_norm_min THEN
                            v_status_check := 'warning';
                        END IF;
                    END IF;

                -- S2: Range threshold (min_value AND max_value)
                ELSE
                    IF v_meting.waarde >= v_norm_min AND v_meting.waarde <= v_norm_max THEN
                        v_status_check := 'oke';
                    ELSE
                        IF v_marge IS NULL OR v_marge = 0 THEN
                            IF v_meting.waarde > v_norm_max OR v_meting.waarde < v_norm_min THEN
                                v_status_check := 'fatal';
                            END IF;
                        ELSE
                            -- With margin (apply margin based on the range between min and max)
                            -- margin amount = abs(max - min) * (marge_percentage / 100)
                            -- compute once for clarity and performance
                            v_range_margin := abs(v_norm_max - v_norm_min) * v_marge / 100.0;

                            IF v_meting.waarde >= (v_norm_min - v_range_margin) AND
                               v_meting.waarde <= v_norm_min THEN
                                v_status_check := 'warning';
                            ELSEIF v_meting.waarde <= (v_norm_max + v_range_margin) AND
                                   v_meting.waarde >= v_norm_max THEN
                                v_status_check := 'warning';
                            ELSE
                                v_status_check := 'fatal';
                            END IF;
                        END IF;
                    END IF;
                END IF;

                -- Insert anomaly if detected
                IF v_status_check != 'oke' THEN
                    v_is_warning := (v_status_check = 'warning');

                    INSERT INTO afwijking (
                        meting_id, meting_time, kunstwerk_id, sensor_id, time,
                        norm_min_waarde, norm_max_waarde, norm_marge_percentage,
                        gemeten_waarde, is_warning
                    ) VALUES (
                        v_meting.id, v_meting.time, v_kunstwerk_id, v_sensor_id, NOW(),
                        v_norm_min, COALESCE(v_norm_max, 0), v_marge,
                        v_meting.waarde, v_is_warning
                    );

                    v_afwijkingen_count := v_afwijkingen_count + 1;
                END IF;

            END LOOP;

            -- Update sensor with last analyzed meting_id
            IF v_last_meting_id IS NOT NULL THEN
                UPDATE sensor
                SET last_analyzed_meting_id = v_last_meting_id
                WHERE id = v_sensor_id;
            END IF;

            -- Return results for this sensor
            RETURN QUERY SELECT
                v_sensor_id,
                v_kunstwerk_id,
                v_metingen_count,
                v_afwijkingen_count,
                COALESCE(v_last_meting_id, v_last_analyzed),
                v_status::VARCHAR;

            -- Reset counters for next sensor
            v_metingen_count := 0;
            v_afwijkingen_count := 0;
            v_last_meting_id := NULL;
            v_status := 'success';

        EXCEPTION WHEN OTHERS THEN
            v_error_msg := SQLERRM;

            -- Log error
            INSERT INTO procedure_error_log (procedure_name, sensor_id, error_message, error_context)
            VALUES ('analyze_sensor_metingen', v_sensor_id, v_error_msg,
                    jsonb_build_object('last_analyzed_id', v_last_analyzed, 'metingen_processed', v_metingen_count));

            -- Return error status for this sensor
            RETURN QUERY SELECT
                v_sensor_id,
                v_kunstwerk_id,
                v_metingen_count,
                v_afwijkingen_count,
                COALESCE(v_last_meting_id, v_last_analyzed),
                'error: ' || v_error_msg;

            -- Continue with next sensor
            v_metingen_count := 0;
            v_afwijkingen_count := 0;
            v_last_meting_id := NULL;
            CONTINUE;

        END;
    END LOOP;

END;
$$ LANGUAGE plpgsql;
