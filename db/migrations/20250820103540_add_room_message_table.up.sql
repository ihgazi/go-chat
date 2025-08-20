CREATE TABLE "room_message" (
    "id" bigserial PRIMARY KEY,
    "room_id" bigint NOT NULL,
    "user_id" bigint NOT NULL,
    "message" text NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    FOREIGN KEY ("room_id") REFERENCES "room" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
); 
