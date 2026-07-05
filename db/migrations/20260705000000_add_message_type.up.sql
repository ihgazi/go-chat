ALTER TABLE room_message ADD COLUMN message_type VARCHAR(50) DEFAULT 'user';
ALTER TABLE room_message ADD COLUMN event VARCHAR(255);
