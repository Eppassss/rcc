# rcc change log

## v8.0.11 (date: 20.1.2021)

- fix: conda-forge and default switch on templates

## v8.0.10 (date: 18.1.2021)

- fix: when there is no pip dependencies, do not try to run pip command

## v8.0.9 (date: 15.1.2021)

- fix: removing one verbosity flag from micromamba invocation

## v8.0.8 (date: 15.1.2021)

- now micromamba 0.7.8 is required
- repodata TTL is reduced to 16 hours, and in case of environment creation
  failure, fall back to 0 seconds TTL (immediate update)
- using new --retry-with-clean-cache option in micromamba

## v8.0.7 (date: 11.1.2021)

- Now rcc manages TEMP and TMP locations for its subprocesses

## v8.0.6 (date: 8.1.2021)

- Updated to robot templates
- conda channels in order for `--strict-channel-priority`
- library versions updated and strict as well (rpaframework v7.1.1)
- Added basic guides for what to do in conda.yaml for end-users.

## v8.0.5 (date: 8.1.2021)

- added robot test to validate required changes, which are common/version.go
  and docs/changelog.md

## v8.0.4 (date: 8.1.2021)

- now requires micromamba 0.7.7 at least, with version check added
- micromamba now brings --repodata-ttl, which rcc currently sets for 7 days
- and touching conda caches is gone because of repodata ttl
- can now also cleanup micromamba binary and with --all
- environment validation checks simplified (no more separate space check)

## v8.0.3 (date: 7.1.2021)

- adding path validation warnings, since they became problem (with pip) now
  that we moved to use micromamba instead of miniconda
- also validation pattern update, with added "~" and "-" as valid characters
- validation is now done on toplevel, so all commands could generate
  those warnings (but currently they don't break anything yet)

## v8.0.2 (date: 5.1.2021)

- fixing failed robot tests for progress indicators (just tests)

## v8.0.1 (date: 5.1.2021)

- added separate pip install phase progress step (just visualization)
- now `rcc env cleanup` has option to remove miniconda3 installation

## v8.0.0 (date: 5.1.2021)

- BREAKING CHANGES
- removed miniconda3 download and installing
- removed all conda commands (check, download, and install)
- environment variables `CONDA_EXE` and `CONDA_PYTHON_EXE` are not available
  anymore (since we don't have conda installation anymore)
- adding micromamba download, installation, and usage functionality
- dropping 32-bit support from windows and linux, this is breaking change,
  so that is why version series goes up to v8

## v7.1.5 (date: 4.1.2021)

- now command `rcc man changelog` shows changelog.md from build moment

## v7.1.4 (date: 4.1.2021)

- bug fix for background metrics not send when application ends too fast
- now all telemetry sending happens in background and synchronized at the end
- added this new changelog.md file

## Older versions

Versions 7.1.3 and older do not have change log entries. This changelog.md
file was started at 4.1.2021.
