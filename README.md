# Media Mapper
Media Mapper is a tool designed to help manage the naming of your media files.
It can be used to identify movies and tv shows based on the commonly used
naming conventions detailed [here](https://en.wikipedia.org/wiki/Pirated_movie_release_types)
and will subsequently format the titles such that they are compatible with
services like Plex.

This project is intended to offer a free alternative to software like filebot,
and verifies titles through TMDB.


## How to use Media Mapper
Currently, Media Mapper can only be run via the command line. The supported
arguments can be listed by running it with the `-h` flag. Below is an example
of how this to run the program:

```console
foo@bar:~$ media-mapper -location /root-dir/of/mediafiles/to/format/
```

If you execute the program with the command above, the program will run for all
of the supported media files contained recursively within the path
`/root-dir/of/mediafiles/to/format/`.

*Note: Before changing any file names, the program will display a list of the
changes and wait for permission to proceed.*

## Supported files
Media Mapper currently supports the following file types:

**Video Formats**
- .asf
- .avi
- .mov
- .mp4
- .ts
- .mkv
- .wmv