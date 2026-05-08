ALTER TABLE afwijking
ADD COLUMN IF NOT EXISTS sensor_id BIGINT;

DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'fk_sensor'
	) THEN
		ALTER TABLE afwijking
		ADD CONSTRAINT fk_sensor FOREIGN KEY (sensor_id) REFERENCES sensor(id);
	END IF;
END $$;
