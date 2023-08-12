CREATE TABLE "files" (
  "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "file_id" varchar NOT NULL,
  "bucket_id" varchar NOT NULL,
  "owner" uuid,
  "name" varchar NOT NULL,
  "size" varchar NOT NULL,
  "favourite" bool NOT NULL DEFAULT false,
  "file_type" varchar NOT NULL,
  "last_modified" varchar NOT NULL DEFAULT (now()),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "files" ("owner");

-- ALTER TABLE "files" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id") ON DELETE SET NULL;
