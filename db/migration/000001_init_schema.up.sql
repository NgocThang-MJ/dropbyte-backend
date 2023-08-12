CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
