<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should have the following
format:

* [#PullRequestNumber](PullRequestLink) message

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState
given same genesisState and txList.
Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

## Unreleased

### Bug fixes

- [#200](https://github.com/babylonlabs-io/babylon-sdk/pull/200) ci: add `check-no-babylon-deps` check
- [#199](https://github.com/babylonlabs-io/babylon-sdk/pull/199) dep: avoid using babylon utilities in demo/ and x/
- [#198](https://github.com/babylonlabs-io/babylon-sdk/pull/198) fix: fix ineffective error propagation

### Improvements

- [#205](https://github.com/babylonlabs-io/babylon-sdk/pull/205) chore: remove custom handlers
- [#203](https://github.com/babylonlabs-io/babylon-sdk/pull/203) chore: refactor querier

## v0.12.0

### Improvements

- [#196](https://github.com/babylonlabs-io/babylon-sdk/pull/196) chore: bump versions and release v0.12.0

## v0.11.0

### State Breaking

- [#165](https://github.com/babylonlabs-io/babylon-sdk/pull/165) refactor: improve contracts storage & remove init
- [#191](https://github.com/babylonlabs-io/babylon-sdk/pull/191) chore: add MaxGasEndBlocker in params and fix tests

### Improvements

- [#194](https://github.com/babylonlabs-io/babylon-sdk/pull/194) chore: bump babylon to v3.0.0-rc.0 and update contracts version
- [#189](https://github.com/babylonlabs-io/babylon-sdk/pull/189) e2e: Add rewards test
- [#188](https://github.com/babylonlabs-io/babylon-sdk/pull/188) chore: refactor error handling and sudo call to avoid chain halt
- [#187](https://github.com/babylonlabs-io/babylon-sdk/pull/187) chore: bump contracts with finality signature fix
- [#186](https://github.com/babylonlabs-io/babylon-sdk/pull/186) e2e: test for fee collector
- [#182](https://github.com/babylonlabs-io/babylon-sdk/pull/182) e2e: reenable and fix e2e tests
- [#172](https://github.com/babylonlabs-io/babylon-sdk/pull/172) chore: Fix e2e tests
- [#108](https://github.com/babylonlabs-io/babylon-sdk/pull/108) chore: Upgrade wasmd to 0.55
- [#116](https://github.com/babylonlabs-io/babylon-sdk/pull/116) add docs for
  the x/babylon module
- [#143](https://github.com/babylonlabs-io/babylon-sdk/pull/143) Re-enable e2e tests
- [#162](https://github.com/babylonlabs-io/babylon-sdk/pull/162) Add standard modules cli in bcd
- [#171](https://github.com/babylonlabs-io/babylon-sdk/pull/171) Upgrade contracts data v0.15.0
- [#164](https://github.com/babylonlabs-io/babylon-sdk/pull/164) Distribute rewards per block from fee collector
