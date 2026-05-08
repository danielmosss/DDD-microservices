ALTER TABLE kunstwerk
ALTER COLUMN last_send_dh_update TYPE TIME
USING last_send_dh_update::time;
