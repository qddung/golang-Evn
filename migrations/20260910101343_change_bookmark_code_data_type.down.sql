-- reverse: create index "idx_code" to table: "bookmarks"
DROP INDEX "public"."idx_code";
-- reverse: modify "bookmarks" table
ALTER TABLE "public"."bookmarks" DROP COLUMN "code_int", ALTER COLUMN "code" SET NOT NULL;
