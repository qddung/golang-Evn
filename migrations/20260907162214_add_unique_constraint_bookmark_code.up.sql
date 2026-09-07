-- create index "idx_bookmark_code" to table: "bookmarks"
CREATE UNIQUE INDEX "idx_bookmark_code" ON "public"."bookmarks" ("code");
