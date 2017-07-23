* Store states in a sqlite3 DB with (at least) 4 tables:
   - `states`
   - `modules`
   - `resources`
   - `attributes`

* replace `buildCache()` with a `feedDB()` function
* have `feedDB()` watch the S3 bucket regularly to update the DB
* use the sqlite DB to query states
* use the sqlite DB for search by statefile, module name, resource type/id,
  attribute name/value
