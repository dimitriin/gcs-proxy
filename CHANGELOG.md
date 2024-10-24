# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2024-10-24

### Added

- Ability to set timeout before handle shutdown signal, could be useful for container synchronization;
- Set `131` exit code on `QUIT` signal;
- Ability to set custom exit codes for different signals (`TERM`, `INT`, `QUIT`).

## [1.0.0] - 2024-09-22

### Added

- Initial release.
