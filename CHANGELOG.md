# Changelog

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
