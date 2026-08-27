-- create "bookmarks" table
CREATE TABLE "public"."bookmarks" (
  "id" uuid NOT NULL,
  "code" text NULL,
  "description" text NULL,
  "url" text NULL,
  "user_id" uuid NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_bookmarks" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_id" to table: "bookmarks"
CREATE UNIQUE INDEX "idx_id" ON "public"."bookmarks" ("url", "user_id");
