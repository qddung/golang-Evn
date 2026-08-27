-- create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "display_name" text NULL,
  "email" text NOT NULL,
  "password" text NOT NULL,
  "user_name" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_user_name" UNIQUE ("user_name")
);
