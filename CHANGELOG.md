## 0.13.0 (X X, 2017)

FEATURES:

- Server:

  * Add S3 key prefix [GH #25]


## 0.12.0 (August 3, 2017)

FEATURES:

- Server:

  * Check versions and path existence in memory to limit DB queries
  * Add log-level and log-format options
  * Log with fields [GH #19]
  * Add a compare API point [GH #20]
  * Split types into a types package
  * Add tf_versions API point
  * Add version API point
  * Sanitize raw SQL queries

- UI:

  * Remove sb-admin theme
  * Make charts in overview clickable, linking to search view
  * Fix display bugs in state view [GH #18]
  * Select first resource of first module on state view load
  * Add a compare function to state view [GH #21]
  * Use $routeParams instead of parsing $location.url()
  * Make state view work without reloading the page
  * Order resource attributes in state view
  * Display long resource attributes and titles with ellipsis
  * Support permalinks and fix form in search view [GH #16]
  * Add tf_version filtering to search view
  * Allow to clear filters in search view
  * Remove unused sorting in tables

FIXES:

- Server:

  * Do not import non-ASCII attribute values [GH #17]
  * Remove --no-sync from docker-compose.yml


## 0.11.0 (August 1, 2017)

FEATURES:

- Server:

  * Optimize various database queries
  * Add state activity to API
  * Add --no-sync flag to disable S3 syncing
  * Get version list from activity API point [GH #10]
  * Retire legacy history API point
  * Add locks API point [GH #11]
  * Add tfversion and types counts to API

- UI:

  * Use a non-fluid container and fix the margins
  * Move index.html to static/ directory
  * List each path only once in overview, with most recent version
  * Add state activity to overview [GH #5]
  * Sort state files in navbar select
  * Add lock information to overview and state view [GH #11]
  * Add charts to overview [GH #12]


## 0.10.1 (July 27, 2017)

FIXES:

- UI:

  * Rename Main link to Overview
  * Use relative path to get back to Overview

## 0.10.0 (July 27, 2017)

FEATURES:

- Server:

  * Change default port to 8080
  * Add flags and help to command line
  * Add version
  * List states from the DB instead of S3

- UI:

  * Add list of state updates to main page
  * Do search on clear and when search page is loaded
  * Use sb-admin Bootstrap theme


FIXES:

- Server:

  * Crash if HTTP Listenandserve fails [GH #1]

## 0.9.0 (July 26, 2017)

**Warning**: The database was ported from SQLite3 to PostgreSQL and needs to be
rebuilt!

FEATURES:

- Server:

  * Port to PostgreSQL [GH #2]
  * Improve and fix search API
  * Add flags and environment variables
  * Improve DB refresh idempotence
  * Migrate all API calls to DB

- UI

  * Use container-fluid instead of container
  * Improve and fix search interface


## 0.8.0 (July 26, 2017)

FEATURES:

* Add a landing page
* Add a search interface


## 0.7.2 (July 24, 2017)

FIXES:

- Server:

  * Do not UpdateState(), select default version instead
  * Do not insert states without versions, select default version instead

- UI:

  * Select default version from state


## 0.7.1 (July 24, 2017)

FIXES:

- Server:

  * Add indices on tables to improve performance

- UI:

  * Fix selected class on resource


## 0.7.0 (July 24, 2017)

FEATURES:

- Server:

  * Use a SQLite3 database to store states locally


## 0.6.0 (July 22, 2017)

FEATURES:

- Server:

  * Cache versioned states
  * Add finer internal methods
  * Add buildCache() (not activated yet)

- UI:

  * Adjust resource list height
  * Mark selected resource


## 0.5.0 (July 21, 2017)

FEATURES:

* Use standard notation `<module>.<resource>`.
* Use hash in URL to track resources.
* Remove `/state/` prefix from breadcrumbs.

IMPROVEMENTS:

- Server:

  * Split and refactor code.
  * Log errors.

- UI:

  * Split navbar from index.html.
  * Split javascript code.


## 0.4.0 (July 20, 2017)

FEATURES:

* Use locale date for versions.


## 0.3.3 (July 20, 2017)

FIXES:

* Fix vertical scrolling issues in UI.


## 0.3.2 (July 20, 2017)

FIXES:

* Use `BASE_URL` in style sheet.


## 0.3.1 (July 20, 2017)

FIXES:

* Fix `BASE_URL` support in API.


## 0.3.0 (July 20, 2017)

FEATURES:

* Add support for `BASE_URL` to change app base URL.


## 0.2.0 (July 20, 2017)

FEATURES:

* Add history to API.
* Add version selection in UI.
* Add screenshot to README.

## 0.1.1 (July 19, 2017)

FIXES:

* Add missing files and README instructions.

## 0.1.0 (July 19, 2017)

* Initial release.
