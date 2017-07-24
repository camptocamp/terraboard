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
