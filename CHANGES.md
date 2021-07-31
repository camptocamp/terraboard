# Changelog

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
