-- create "bookmarks" table
CREATE TABLE "public"."bookmarks" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "code" text NOT NULL,
  "description" text NULL,
  "url" text NOT NULL,
  "user_id" uuid NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_bookmarks" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_bookmarks_deleted_at" to table: "bookmarks"
CREATE INDEX "idx_bookmarks_deleted_at" ON "public"."bookmarks" ("deleted_at");
-- create index "idx_id" to table: "bookmarks"
CREATE UNIQUE INDEX "idx_id" ON "public"."bookmarks" ("url", "user_id");
