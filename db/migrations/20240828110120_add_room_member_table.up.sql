CREATE TABLE "room_member" (
    "id" bigserial PRIMARY KEY,
    "room_id" bigint NOT NULL,
    "user_id" bigint NOT NULL,
    "join_date" timestamp NOT NULL DEFAULT now(),
    "last_online" timestamp NOT NULL DEFAULT now(),
    FOREIGN KEY ("room_id") REFERENCES "room" ("id"),
    FOREIGN KEY ("user_id") REFERENCES "users" ("id")
)
