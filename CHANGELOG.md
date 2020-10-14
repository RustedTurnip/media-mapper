# Changelog
All notable changes to this project will be documented in this file.

*Note: Additional information on this changelog can be found in the [footnote](#a-namefootnotefootnotea).*

## v0.2.0 - 2020-10-14
### Added
- Added a break in operation in between identifying name changes and applying
the name changes where user input is required. This offers users a chance to
prevent unwanted changes.
- Populated the [README](https://github.com/RustedTurnip/media-mapper/blob/master/README.md).

### Fixed
- Added check for empty `location` argument to prevent unwanted behaviour.



## v0.1.0 - 2020-09-16
### Added
- Added changelog (using the convention defined in https://keepachangelog.com/en/1.0.0/).
- File renaming has now been added so this tool can now be used as intended.
- Added `-version` commandline flag.

### Changed
- Better logging has been implemented replacing the few non-descriptive log
lines.
- The command line flag specifying auth location has been changed from 
`-authentication` to `-auth`.

### Removed
- The fatal panics scattered throughout the program have been removed as they
were unnecessary.



## v0.0.0 - 2020-09-14 (unreleased)
### Added
- Command-line arguments including:
    - `-database` - to specify the online media database titles are checked
    against.
    - `-authentication` - to pass the location of the auth file (for database
    access).
    - `-location` - path to the root directory of the file tree to be
    searched.
- Recursive file scanning so that all media files contained within the
specified root directory.
- [TMDB API](https://www.themoviedb.org/) support so that media titles can be
verified.

## <a name="footnote">*Footnote*</a>
- *The changelog format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).*
- *This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).*

[v0.0.0]: https://github.com/RustedTurnip/media-mapper/releases/tag/v0.0.0