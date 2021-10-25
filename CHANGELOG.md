# Changelog

## [2.0.0](https://www.github.com/camptocamp/terraboard/compare/v2.0.0...v2.0.0) (2021-10-25)


### Features

* **frontend:** add loading spinner during plan fetching ([4e6d2b6](https://www.github.com/camptocamp/terraboard/commit/4e6d2b6e73c704977bc7c217a383638f6120c044))
* **frontend:** highlightjs theme ([#202](https://www.github.com/camptocamp/terraboard/issues/202)) ([e18dd8f](https://www.github.com/camptocamp/terraboard/commit/e18dd8fd714d36f98e7dc38d246ba36e209534a8))
* **frontend:** now trigger refresh on input for search's text fields ([#205](https://www.github.com/camptocamp/terraboard/issues/205)) ([cd42ac1](https://www.github.com/camptocamp/terraboard/commit/cd42ac1a3396cb05821ca9c769a8a4eb010628a2))
* **frontend:** resource filter now match any substring in res path or name ([cc707ea](https://www.github.com/camptocamp/terraboard/commit/cc707ea5d135a92f28715122299d3d889b422988))
* **frontend:** set role=button on list carets ([#200](https://www.github.com/camptocamp/terraboard/issues/200)) ([02c0022](https://www.github.com/camptocamp/terraboard/commit/02c0022a02c7d7ee7e591e5af3dbe76056a0e8eb))
* **frontend:** use pointer cursor on list group items ([#201](https://www.github.com/camptocamp/terraboard/issues/201)) ([0b87348](https://www.github.com/camptocamp/terraboard/commit/0b873487dd052223bf3fa607855272e9169e1d06))
* **frontend:** use UTC strings for dates ([#198](https://www.github.com/camptocamp/terraboard/issues/198)) ([c77d66f](https://www.github.com/camptocamp/terraboard/commit/c77d66f31f7c83e082097ccb1200e028e2c4ca95))
* **plan:** process plan's status from ci exit-code ([#215](https://www.github.com/camptocamp/terraboard/issues/215)) ([5b32184](https://www.github.com/camptocamp/terraboard/commit/5b32184086d3a9be864d418301e086cd04a30c9b))
* **plans:** add plans explorer and viewer  ([#193](https://www.github.com/camptocamp/terraboard/issues/193)) ([35d73c0](https://www.github.com/camptocamp/terraboard/commit/35d73c016b873d8162d445bd56c478632d497703))


### Bug Fixes

* **frontend:** avatar not displaying due to css issue ([#206](https://www.github.com/camptocamp/terraboard/issues/206)) ([78f8922](https://www.github.com/camptocamp/terraboard/commit/78f8922f86a36e4b9465070d56d8cf0bb7615e59))
* **frontend:** missing user's information & avatar ([#204](https://www.github.com/camptocamp/terraboard/issues/204)) ([0aec334](https://www.github.com/camptocamp/terraboard/commit/0aec3342a5041abce1cd9f56334694012dafa454))
* **frontend:** performance issue with resource filter with multiples modules ([#208](https://www.github.com/camptocamp/terraboard/issues/208)) ([cc707ea](https://www.github.com/camptocamp/terraboard/commit/cc707ea5d135a92f28715122299d3d889b422988))
* **frontend:** performence issue on plans fetching ([#212](https://www.github.com/camptocamp/terraboard/issues/212)) ([4e6d2b6](https://www.github.com/camptocamp/terraboard/commit/4e6d2b6e73c704977bc7c217a383638f6120c044))
* **frontend:** resource filter on state view now works ([#203](https://www.github.com/camptocamp/terraboard/issues/203)) ([156e544](https://www.github.com/camptocamp/terraboard/commit/156e544b3f351bf2bc6008ebc26be8d45193f933))
* **frontend:** undefined error on plan view without outputs changes ([#211](https://www.github.com/camptocamp/terraboard/issues/211)) ([a5d6050](https://www.github.com/camptocamp/terraboard/commit/a5d6050fd9c878d9773f6e1669eabea2a00069c3))
* **json:** plan's variables parsing error ([#210](https://www.github.com/camptocamp/terraboard/issues/210)) ([91b3dab](https://www.github.com/camptocamp/terraboard/commit/91b3dab01633729d0143e99629363968f9f84c0b))
* remove localhost from URLs ([#196](https://www.github.com/camptocamp/terraboard/issues/196)) ([7dd34ed](https://www.github.com/camptocamp/terraboard/commit/7dd34ed6fdd1781a6c76a6753e406fb7dc25d33f))


### Miscellaneous Chores

* release as 2.0.0 ([a2dd9b6](https://www.github.com/camptocamp/terraboard/commit/a2dd9b66b2ea5eae1a3102414261af414cce9202))

## [2.0.0-alpha.2](https://www.github.com/camptocamp/terraboard/compare/v2.0.0...v2.0.0-alpha.2) (2021-08-19)


### Features

* **frontend:** highlightjs theme ([#202](https://www.github.com/camptocamp/terraboard/issues/202)) ([e18dd8f](https://www.github.com/camptocamp/terraboard/commit/e18dd8fd714d36f98e7dc38d246ba36e209534a8))
* **frontend:** now trigger refresh on input for search's text fields ([#205](https://www.github.com/camptocamp/terraboard/issues/205)) ([cd42ac1](https://www.github.com/camptocamp/terraboard/commit/cd42ac1a3396cb05821ca9c769a8a4eb010628a2))
* **frontend:** set role=button on list carets ([#200](https://www.github.com/camptocamp/terraboard/issues/200)) ([02c0022](https://www.github.com/camptocamp/terraboard/commit/02c0022a02c7d7ee7e591e5af3dbe76056a0e8eb))
* **frontend:** use pointer cursor on list group items ([#201](https://www.github.com/camptocamp/terraboard/issues/201)) ([0b87348](https://www.github.com/camptocamp/terraboard/commit/0b873487dd052223bf3fa607855272e9169e1d06))
* **frontend:** use UTC strings for dates ([#198](https://www.github.com/camptocamp/terraboard/issues/198)) ([c77d66f](https://www.github.com/camptocamp/terraboard/commit/c77d66f31f7c83e082097ccb1200e028e2c4ca95))


### Bug Fixes

* **frontend:** avatar not displaying due to css issue ([#206](https://www.github.com/camptocamp/terraboard/issues/206)) ([78f8922](https://www.github.com/camptocamp/terraboard/commit/78f8922f86a36e4b9465070d56d8cf0bb7615e59))
* **frontend:** missing user's information & avatar ([#204](https://www.github.com/camptocamp/terraboard/issues/204)) ([0aec334](https://www.github.com/camptocamp/terraboard/commit/0aec3342a5041abce1cd9f56334694012dafa454))
* **frontend:** resource filter on state view now works ([#203](https://www.github.com/camptocamp/terraboard/issues/203)) ([156e544](https://www.github.com/camptocamp/terraboard/commit/156e544b3f351bf2bc6008ebc26be8d45193f933))
* remove localhost from URLs ([#196](https://www.github.com/camptocamp/terraboard/issues/196)) ([7dd34ed](https://www.github.com/camptocamp/terraboard/commit/7dd34ed6fdd1781a6c76a6753e406fb7dc25d33f))

## [2.0.0-alpha](https://www.github.com/camptocamp/terraboard/compare/v1.1.0...v2.0.0-alpha) (2021-07-31)


### ⚠ BREAKING CHANGES

* use lineage instead of path to link states on overview (#179)
* multiple buckets / providers support (#48) (#177)

### Features

* add plan db scheme ([#164](https://www.github.com/camptocamp/terraboard/issues/164)) ([2c1cd8f](https://www.github.com/camptocamp/terraboard/commit/2c1cd8f003016ee3c54a57cf6d679e6ee2fc5b29))
* add test environment with proper docker-compose using MinIO + add compatibility for it to AWS provider ([#165](https://www.github.com/camptocamp/terraboard/issues/165)) ([6d44ecb](https://www.github.com/camptocamp/terraboard/commit/6d44ecbf5765e33b2cfcdbef1dce9059fe7dc94a))
* **api:** add lineages get endpoint ([#176](https://www.github.com/camptocamp/terraboard/issues/176)) ([f632586](https://www.github.com/camptocamp/terraboard/commit/f6325861cca0306382ccfd897740fa25911fc27c))
* **api:** add plan submit/get endpoints   ([#175](https://www.github.com/camptocamp/terraboard/issues/175)) ([f341ffd](https://www.github.com/camptocamp/terraboard/commit/f341ffdb0c990414e349a3596221dbdd8d4668e9))
* build on Terraform v1.0.2 ([#185](https://www.github.com/camptocamp/terraboard/issues/185)) ([505ef82](https://www.github.com/camptocamp/terraboard/commit/505ef828e63d8ff61cd166ff9ab89086248e7233))
* front-end redesign to VueJS/Bootstrap 5 & server improvements ([#188](https://www.github.com/camptocamp/terraboard/issues/188)) ([5001de9](https://www.github.com/camptocamp/terraboard/commit/5001de9d4c4271bf75ba10580ba4f0038d9d8f4f))
* **gorm:** gorm v2 migration ([#170](https://www.github.com/camptocamp/terraboard/issues/170)) ([630ce3f](https://www.github.com/camptocamp/terraboard/commit/630ce3fe4044f44d66dad854430de3bee3412591))
* **migration:** add migration compatibility with states without lineage ([#183](https://www.github.com/camptocamp/terraboard/issues/183)) ([e50bbf6](https://www.github.com/camptocamp/terraboard/commit/e50bbf63807be983d1fa17b4bbcbd9873f2d0081))
* multiple buckets / providers support ([#48](https://www.github.com/camptocamp/terraboard/issues/48)) ([#177](https://www.github.com/camptocamp/terraboard/issues/177)) ([e44ebce](https://www.github.com/camptocamp/terraboard/commit/e44ebce133bc719e339446ca71390086eae50d85))
* new standalone lineage table + associated migration ([#173](https://www.github.com/camptocamp/terraboard/issues/173)) ([aa7d455](https://www.github.com/camptocamp/terraboard/commit/aa7d455cf9ef86651ad0791b02747502b8a44c4e))
* use lineage instead of path to link states on overview ([#179](https://www.github.com/camptocamp/terraboard/issues/179)) ([c576d95](https://www.github.com/camptocamp/terraboard/commit/c576d95fc1942b78e5a2b53318d22ab00778eee8))


### Bug Fixes

* terraboard crash at compose startup if db isn't fully initialized ([#174](https://www.github.com/camptocamp/terraboard/issues/174)) ([a8e6bdd](https://www.github.com/camptocamp/terraboard/commit/a8e6bdd7a9195dd54af8079a4d9c1101e6cadbb9))
* **view:** style issue which cropped overview content ([#189](https://www.github.com/camptocamp/terraboard/issues/189)) ([ee9f923](https://www.github.com/camptocamp/terraboard/commit/ee9f92309ad74ce809c121337a2098b5176d4790))


### Miscellaneous Chores

* release 2.0.0-alpha ([aa3cd25](https://www.github.com/camptocamp/terraboard/commit/aa3cd25a17fb237c495432ebb9568e8e73a84bcd))


## [1.1.0](https://github.com/camptocamp/terraboard/tree/1.1.0) (2021-04-14)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/1.0.0...1.1.0)

**Implemented enhancements:**

- View for outputs [\#50](https://github.com/camptocamp/terraboard/issues/50)
- Update Terraform to 0.14.10 [\#149](https://github.com/camptocamp/terraboard/pull/149) ([raphink](https://github.com/raphink))
- Implement AWS EXTERNAL ID feature [\#142](https://github.com/camptocamp/terraboard/pull/142) ([alemuro](https://github.com/alemuro))
- Implement Outputs [\#140](https://github.com/camptocamp/terraboard/pull/140) ([hoshsadiq](https://github.com/hoshsadiq))

**Fixed bugs:**

- Fix environment variable for APPRoleArn on README [\#141](https://github.com/camptocamp/terraboard/pull/141) ([alemuro](https://github.com/alemuro))

**Closed issues:**

- support of terraform 0.14.3 [\#137](https://github.com/camptocamp/terraboard/issues/137)

## [1.0.0](https://github.com/camptocamp/terraboard/tree/1.0.0) (2021-02-17)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.22.0...1.0.0)

**Breaking changes:**

- Fixed failing CI following recent changes introducing GitLab backend support [\#133](https://github.com/camptocamp/terraboard/pull/133) ([mvisonneau](https://github.com/mvisonneau))
- Support multiple file extensions in S3 [\#122](https://github.com/camptocamp/terraboard/pull/122) ([Moglum](https://github.com/Moglum))

**Implemented enhancements:**

- Allow setting region and force path style for s3 [\#136](https://github.com/camptocamp/terraboard/pull/136) ([hoshsadiq](https://github.com/hoshsadiq))
- option to switch go dns resolver [\#135](https://github.com/camptocamp/terraboard/pull/135) ([mihaiplesa](https://github.com/mihaiplesa))
- bump terraform to 0.13.6 [\#134](https://github.com/camptocamp/terraboard/pull/134) ([mihaiplesa](https://github.com/mihaiplesa))
- Bumped all dependencies / TF 0.14.5 / TFE 0.12.0 [\#131](https://github.com/camptocamp/terraboard/pull/131) ([mvisonneau](https://github.com/mvisonneau))
- Added support for GitLab as terraform state backend provider [\#130](https://github.com/camptocamp/terraboard/pull/130) ([mvisonneau](https://github.com/mvisonneau))
- made the goreleaser config a bit more readable [\#129](https://github.com/camptocamp/terraboard/pull/129) ([mvisonneau](https://github.com/mvisonneau))
- Fix/docker instructions [\#124](https://github.com/camptocamp/terraboard/pull/124) ([uritau](https://github.com/uritau))
- bump terraform to 0.13.5 [\#120](https://github.com/camptocamp/terraboard/pull/120) ([filiptepper](https://github.com/filiptepper))
- Bump to terraform 0.13.4 [\#119](https://github.com/camptocamp/terraboard/pull/119) ([ouranos](https://github.com/ouranos))
- Bump to terraform 0.13.3 [\#117](https://github.com/camptocamp/terraboard/pull/117) ([chelnak](https://github.com/chelnak))
- Add support for multiple instances of the same resource [\#115](https://github.com/camptocamp/terraboard/pull/115) ([Wiston999](https://github.com/Wiston999))
- Implemented gosec as part of the testing workflow [\#114](https://github.com/camptocamp/terraboard/pull/114) ([mvisonneau](https://github.com/mvisonneau))
- Fixed Dockerfiles following recent changes [\#113](https://github.com/camptocamp/terraboard/pull/113) ([mvisonneau](https://github.com/mvisonneau))
- Enhanced Makefile and testing implementation [\#112](https://github.com/camptocamp/terraboard/pull/112) ([mvisonneau](https://github.com/mvisonneau))
- Bumped to golang 1.15 and upgraded all modules to their latest version [\#111](https://github.com/camptocamp/terraboard/pull/111) ([mvisonneau](https://github.com/mvisonneau))
- Upgraded terraform to 0.13.2 [\#110](https://github.com/camptocamp/terraboard/pull/110) ([mvisonneau](https://github.com/mvisonneau))
- ADDED release URL and updated links to open in new tab [\#108](https://github.com/camptocamp/terraboard/pull/108) ([azhar22k](https://github.com/azhar22k))

**Fixed bugs:**

- Resources created using `for\_each` are missing [\#97](https://github.com/camptocamp/terraboard/issues/97)
- Fix nill pointer [\#138](https://github.com/camptocamp/terraboard/pull/138) ([Moglum](https://github.com/Moglum))
- Fixed ineffassign definition following a recent update [\#132](https://github.com/camptocamp/terraboard/pull/132) ([mvisonneau](https://github.com/mvisonneau))

**Closed issues:**

- Unable to start terraboard in Docker following Readme instructions  [\#123](https://github.com/camptocamp/terraboard/issues/123)
- v0.21.0 can't connect to RDS db when running on AWS ECS [\#118](https://github.com/camptocamp/terraboard/issues/118)
- Add support for terraform 0.13.3 [\#116](https://github.com/camptocamp/terraboard/issues/116)
- Use docker networks instead of deprecated --link. [\#45](https://github.com/camptocamp/terraboard/issues/45)

## [0.22.0](https://github.com/camptocamp/terraboard/tree/0.22.0) (2020-08-12)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.21.0...0.22.0)

**Implemented enhancements:**

- Support Terraform v0.13.0 [\#104](https://github.com/camptocamp/terraboard/pull/104) ([raphink](https://github.com/raphink))
- Add lineage to states in DB [\#103](https://github.com/camptocamp/terraboard/pull/103) ([raphink](https://github.com/raphink))
- Support Terraform 0.12.29 [\#102](https://github.com/camptocamp/terraboard/pull/102) ([raphink](https://github.com/raphink))

## [0.21.0](https://github.com/camptocamp/terraboard/tree/0.21.0) (2020-07-22)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.20.0...0.21.0)

**Implemented enhancements:**

- Use S3 compatible bucket [\#72](https://github.com/camptocamp/terraboard/issues/72)
- Update to terraform v0.12.28 [\#98](https://github.com/camptocamp/terraboard/pull/98) ([mhaley-miovision](https://github.com/mhaley-miovision))
- Add AWS ARN Support [\#95](https://github.com/camptocamp/terraboard/pull/95) ([fafalafafa](https://github.com/fafalafafa))
- GCS support [\#91](https://github.com/camptocamp/terraboard/pull/91) ([tristanvanthielen](https://github.com/tristanvanthielen))
- Add optional AWS endpoint [\#81](https://github.com/camptocamp/terraboard/pull/81) ([raphink](https://github.com/raphink))

**Closed issues:**

- panic: runtime error: invalid memory address or nil pointer dereference [\#90](https://github.com/camptocamp/terraboard/issues/90)

## [0.20.0](https://github.com/camptocamp/terraboard/tree/0.20.0) (2020-06-09)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.19.0...0.20.0)

**Implemented enhancements:**

- Update for terraform 0.12.26 [\#94](https://github.com/camptocamp/terraboard/pull/94) ([raphink](https://github.com/raphink))
- update to use terraform 0.12.24 [\#88](https://github.com/camptocamp/terraboard/pull/88) ([BlackWebWolf](https://github.com/BlackWebWolf))
- Update documentation Terraform Cloud \(TFE\) support [\#82](https://github.com/camptocamp/terraboard/pull/82) ([binlab](https://github.com/binlab))
- update to use terraform 0.12.21 [\#80](https://github.com/camptocamp/terraboard/pull/80) ([rickhlx](https://github.com/rickhlx))
- Make database sslmode default to require [\#78](https://github.com/camptocamp/terraboard/pull/78) ([raphink](https://github.com/raphink))
- docs: add --db-sslmode flag [\#77](https://github.com/camptocamp/terraboard/pull/77) ([tedder](https://github.com/tedder))

**Fixed bugs:**

- change PORT to TERRABOARD\_PORT in README.md [\#79](https://github.com/camptocamp/terraboard/pull/79) ([jimsheldon](https://github.com/jimsheldon))

**Closed issues:**

- Terraform v0.12.26 [\#93](https://github.com/camptocamp/terraboard/issues/93)
- pq sslmode [\#76](https://github.com/camptocamp/terraboard/issues/76)

## [0.19.0](https://github.com/camptocamp/terraboard/tree/0.19.0) (2020-01-24)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.18.0...0.19.0)

**Implemented enhancements:**

- Output Terraform lib version in version string [\#75](https://github.com/camptocamp/terraboard/pull/75) ([raphink](https://github.com/raphink))
- Use Terraform 0.12.20 [\#74](https://github.com/camptocamp/terraboard/pull/74) ([raphink](https://github.com/raphink))

**Fixed bugs:**

- Terraform version in changelog [\#73](https://github.com/camptocamp/terraboard/issues/73)

## [0.18.0](https://github.com/camptocamp/terraboard/tree/0.18.0) (2020-01-20)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.17.0...0.18.0)

**Merged pull requests:**

- Update go modules [\#71](https://github.com/camptocamp/terraboard/pull/71) ([raphink](https://github.com/raphink))
- Add setting for DB sslmode [\#70](https://github.com/camptocamp/terraboard/pull/70) ([raphink](https://github.com/raphink))

## [0.17.0](https://github.com/camptocamp/terraboard/tree/0.17.0) (2020-01-20)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.16.0...0.17.0)

**Merged pull requests:**

- Reduce the noise in the logs when parsing statefiles [\#68](https://github.com/camptocamp/terraboard/pull/68) ([mvisonneau](https://github.com/mvisonneau))
- Implemented support for Terraform Cloud as state provider [\#67](https://github.com/camptocamp/terraboard/pull/67) ([mvisonneau](https://github.com/mvisonneau))
- Bumped go version to 1.13 [\#66](https://github.com/camptocamp/terraboard/pull/66) ([mvisonneau](https://github.com/mvisonneau))
- Add support for newer terraform versions [\#65](https://github.com/camptocamp/terraboard/pull/65) ([slitsevych](https://github.com/slitsevych))

## [0.16.0](https://github.com/camptocamp/terraboard/tree/0.16.0) (2019-10-30)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.15.0...0.16.0)

**Closed issues:**

- Compatibility Issue with Terraform v0.12.0  [\#59](https://github.com/camptocamp/terraboard/issues/59)

**Merged pull requests:**

- New flat logo [\#69](https://github.com/camptocamp/terraboard/pull/69) ([raphink](https://github.com/raphink))
- Config file support [\#64](https://github.com/camptocamp/terraboard/pull/64) ([AndresCidoncha](https://github.com/AndresCidoncha))

## [0.15.0](https://github.com/camptocamp/terraboard/tree/0.15.0) (2019-10-24)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.14.3...0.15.0)

**Closed issues:**

- simple install issues [\#63](https://github.com/camptocamp/terraboard/issues/63)
- Terraboard unable to DNS look up DB on ECS [\#61](https://github.com/camptocamp/terraboard/issues/61)
- Don't know how to install & run [\#44](https://github.com/camptocamp/terraboard/issues/44)
- base [\#43](https://github.com/camptocamp/terraboard/issues/43)

**Merged pull requests:**

- Support Terraform 0.12 [\#62](https://github.com/camptocamp/terraboard/pull/62) ([raphink](https://github.com/raphink))
- Correct Travis fail [\#55](https://github.com/camptocamp/terraboard/pull/55) ([gliptak](https://github.com/gliptak))
- Correct Travis fail [\#54](https://github.com/camptocamp/terraboard/pull/54) ([gliptak](https://github.com/gliptak))
- Switch submodule hugo-elate-theme to https [\#53](https://github.com/camptocamp/terraboard/pull/53) ([gliptak](https://github.com/gliptak))
- Bring Go to 1.11 [\#52](https://github.com/camptocamp/terraboard/pull/52) ([gliptak](https://github.com/gliptak))
- Remove deprecated links, use compose default network [\#46](https://github.com/camptocamp/terraboard/pull/46) ([kwerey](https://github.com/kwerey))

## [0.14.3](https://github.com/camptocamp/terraboard/tree/0.14.3) (2017-12-07)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.14.2...0.14.3)

## [0.14.2](https://github.com/camptocamp/terraboard/tree/0.14.2) (2017-11-30)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.14.1...0.14.2)

**Merged pull requests:**

- Detect if state is nil [\#40](https://github.com/camptocamp/terraboard/pull/40) ([raphink](https://github.com/raphink))

## [0.14.1](https://github.com/camptocamp/terraboard/tree/0.14.1) (2017-11-29)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.14.0...0.14.1)

**Closed issues:**

- Index creation issue on blocking parsing of some statefiles [\#37](https://github.com/camptocamp/terraboard/issues/37)

**Merged pull requests:**

- Remove index for attribute value [\#39](https://github.com/camptocamp/terraboard/pull/39) ([raphink](https://github.com/raphink))
- Add a resource filter to the state view [\#31](https://github.com/camptocamp/terraboard/pull/31) ([raphink](https://github.com/raphink))

## [0.14.0](https://github.com/camptocamp/terraboard/tree/0.14.0) (2017-11-25)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.13.0...0.14.0)

**Merged pull requests:**

- Fix json field for value [\#41](https://github.com/camptocamp/terraboard/pull/41) ([raphink](https://github.com/raphink))
- Add option to change state file extension [\#38](https://github.com/camptocamp/terraboard/pull/38) ([gordonbondon](https://github.com/gordonbondon))
- Build improvements [\#35](https://github.com/camptocamp/terraboard/pull/35) ([raphink](https://github.com/raphink))
- Report errors fetching state from S3 [\#34](https://github.com/camptocamp/terraboard/pull/34) ([ant1441](https://github.com/ant1441))

## [0.13.0](https://github.com/camptocamp/terraboard/tree/0.13.0) (2017-08-16)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.12.0...0.13.0)

**Implemented enhancements:**

- Authentication info [\#29](https://github.com/camptocamp/terraboard/pull/29) ([raphink](https://github.com/raphink))
- S3 key prefix [\#25](https://github.com/camptocamp/terraboard/pull/25) ([leonidcliqz](https://github.com/leonidcliqz))

**Merged pull requests:**

- Test compare 100% coverage [\#30](https://github.com/camptocamp/terraboard/pull/30) ([cryptobioz](https://github.com/cryptobioz))
- Lint code \[WIP\] [\#28](https://github.com/camptocamp/terraboard/pull/28) ([raphink](https://github.com/raphink))
- Add Makefile and use it in .travis.yml [\#27](https://github.com/camptocamp/terraboard/pull/27) ([raphink](https://github.com/raphink))
- Add tests for utils and compare [\#26](https://github.com/camptocamp/terraboard/pull/26) ([cryptobioz](https://github.com/cryptobioz))
- Improve logging [\#19](https://github.com/camptocamp/terraboard/pull/19) ([raphink](https://github.com/raphink))

## [0.12.0](https://github.com/camptocamp/terraboard/tree/0.12.0) (2017-08-03)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.11.0...0.12.0)

**Implemented enhancements:**

- Use location in search [\#7](https://github.com/camptocamp/terraboard/issues/7)
- Add a compare view [\#6](https://github.com/camptocamp/terraboard/issues/6)

**Merged pull requests:**

- Merge compare view into state view [\#24](https://github.com/camptocamp/terraboard/pull/24) ([raphink](https://github.com/raphink))
- Add a compare API and view [\#21](https://github.com/camptocamp/terraboard/pull/21) ([raphink](https://github.com/raphink))
- Add compare API point [\#20](https://github.com/camptocamp/terraboard/pull/20) ([raphink](https://github.com/raphink))
- Fix bugs in state view [\#18](https://github.com/camptocamp/terraboard/pull/18) ([raphink](https://github.com/raphink))
- Do not import non-ASCII data [\#17](https://github.com/camptocamp/terraboard/pull/17) ([raphink](https://github.com/raphink))
- Improve search engine [\#16](https://github.com/camptocamp/terraboard/pull/16) ([cryptobioz](https://github.com/cryptobioz))

## [0.11.0](https://github.com/camptocamp/terraboard/tree/0.11.0) (2017-08-01)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.10.1...0.11.0)

**Implemented enhancements:**

- Make sparklines clickable [\#9](https://github.com/camptocamp/terraboard/issues/9)
- Get lock status [\#3](https://github.com/camptocamp/terraboard/issues/3)

**Merged pull requests:**

- Make sparklines clickable [\#15](https://github.com/camptocamp/terraboard/pull/15) ([raphink](https://github.com/raphink))
- Sort versions using an array of objects [\#14](https://github.com/camptocamp/terraboard/pull/14) ([raphink](https://github.com/raphink))
- Fix lock chart async [\#13](https://github.com/camptocamp/terraboard/pull/13) ([cryptobioz](https://github.com/cryptobioz))
- Add charts [\#12](https://github.com/camptocamp/terraboard/pull/12) ([cryptobioz](https://github.com/cryptobioz))
- Add lock info [\#11](https://github.com/camptocamp/terraboard/pull/11) ([raphink](https://github.com/raphink))
- state: get version list from api/state/activity [\#10](https://github.com/camptocamp/terraboard/pull/10) ([raphink](https://github.com/raphink))
- Add sparklines to overview [\#5](https://github.com/camptocamp/terraboard/pull/5) ([raphink](https://github.com/raphink))

## [0.10.1](https://github.com/camptocamp/terraboard/tree/0.10.1) (2017-07-27)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.10.0...0.10.1)

## [0.10.0](https://github.com/camptocamp/terraboard/tree/0.10.0) (2017-07-27)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.9.0...0.10.0)

**Fixed bugs:**

- Can't get terraboard running locally \(without docker\) [\#1](https://github.com/camptocamp/terraboard/issues/1)

## [0.9.0](https://github.com/camptocamp/terraboard/tree/0.9.0) (2017-07-26)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.8.0...0.9.0)

**Merged pull requests:**

- Port to PostgreSQL [\#2](https://github.com/camptocamp/terraboard/pull/2) ([raphink](https://github.com/raphink))

## [0.8.0](https://github.com/camptocamp/terraboard/tree/0.8.0) (2017-07-25)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.7.2...0.8.0)

## [0.7.2](https://github.com/camptocamp/terraboard/tree/0.7.2) (2017-07-24)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.7.1...0.7.2)

## [0.7.1](https://github.com/camptocamp/terraboard/tree/0.7.1) (2017-07-24)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.7.0...0.7.1)

## [0.7.0](https://github.com/camptocamp/terraboard/tree/0.7.0) (2017-07-24)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.6.0...0.7.0)

## [0.6.0](https://github.com/camptocamp/terraboard/tree/0.6.0) (2017-07-22)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.5.0...0.6.0)

## [0.5.0](https://github.com/camptocamp/terraboard/tree/0.5.0) (2017-07-20)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.4.0...0.5.0)

## [0.4.0](https://github.com/camptocamp/terraboard/tree/0.4.0) (2017-07-20)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.3.3...0.4.0)

## [0.3.3](https://github.com/camptocamp/terraboard/tree/0.3.3) (2017-07-20)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.3.2...0.3.3)

## [0.3.2](https://github.com/camptocamp/terraboard/tree/0.3.2) (2017-07-20)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.3.1...0.3.2)

## [0.3.1](https://github.com/camptocamp/terraboard/tree/0.3.1) (2017-07-20)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.3.0...0.3.1)

## [0.3.0](https://github.com/camptocamp/terraboard/tree/0.3.0) (2017-07-20)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.2.0...0.3.0)

## [0.2.0](https://github.com/camptocamp/terraboard/tree/0.2.0) (2017-07-19)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.1.1...0.2.0)

## [0.1.1](https://github.com/camptocamp/terraboard/tree/0.1.1) (2017-07-19)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/0.1.0...0.1.1)

## [0.1.0](https://github.com/camptocamp/terraboard/tree/0.1.0) (2017-07-19)

[Full Changelog](https://github.com/camptocamp/terraboard/compare/a1e76f6fe37cc64e01b4a142d6f749cabd6f9170...0.1.0)



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
