-- modify "bookmarks" table
ALTER TABLE "public"."bookmarks" ALTER COLUMN "code" DROP NOT NULL, ADD COLUMN "code_int" bigserial NOT NULL;
-- create index "idx_code" to table: "bookmarks"
CREATE UNIQUE INDEX "idx_code" ON "public"."bookmarks" ("code_int");
