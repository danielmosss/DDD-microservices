ALTER TABLE kunstwerk
ALTER COLUMN last_send_dh_update TYPE TIMESTAMPTZ
USING CURRENT_DATE + last_send_dh_update;
