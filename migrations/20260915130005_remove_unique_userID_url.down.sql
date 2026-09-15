-- reverse: drop index "idx_id" from table: "bookmarks"
CREATE UNIQUE INDEX "idx_id" ON "public"."bookmarks" ("url", "user_id");
