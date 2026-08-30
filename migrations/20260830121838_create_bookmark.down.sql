-- reverse: create index "idx_id" to table: "bookmarks"
DROP INDEX "public"."idx_id";
-- reverse: create index "idx_bookmarks_deleted_at" to table: "bookmarks"
DROP INDEX "public"."idx_bookmarks_deleted_at";
-- reverse: create "bookmarks" table
DROP TABLE "public"."bookmarks";
