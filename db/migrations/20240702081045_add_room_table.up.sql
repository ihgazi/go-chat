CREATE TABLE "room" (
    "id" bigserial PRIMARY KEY,
    "name" varchar NOT NULL,
    "create_date" timestamp NOT NULL DEFAULT now()
)
