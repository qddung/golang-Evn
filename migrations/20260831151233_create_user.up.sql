DROP TABLE IF EXISTS "public"."users";

-- create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "display_name" text NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  "user_name" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_user_name" UNIQUE ("user_name")
);
-- create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
