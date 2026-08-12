# ⚠️ This project has moved to Codeberg

**Development on GitHub has been permanently halted.** This repository is now
archived and serves as a read-only historical record for existing dependencies.

### New Home

`media-mapper` is still being maintained and updated. It can now be found at
https://codeberg.org/rustedturnip/media-mapper.

### Why the move?
My confidence in GitHub has been repeatedly knocked of late; given growing
downtime, the recent
[merge queue issue](https://www.githubstatus.com/incidents/zsg1lk7w13cf) causing
silent reverts, and most recently the
[source code breach](https://github.blog/security/investigating-unauthorized-access-to-githubs-internal-repositories/)
and the risk of future vulnerability discovery that it brings, I have decided to
migrate my projects to [Codeberg](codeberg.org).

Codeberg is a non-profit, privacy-respecting, and open-source platform powered
by Forgejo.
[In their own words](https://github.blog/security/investigating-unauthorized-access-to-githubs-internal-repositories/):

> The platform we choose for hosting says a lot about the values of our
> ecosystem. Codeberg's open-source roots and non-profit status help assure us
> that their interests are in the collaboration. Further, Codeberg's tooling is
> moving toward federated interoperability which will make it even easier to
> work across services in the future.

### What you need to do:

If you are using any part of `media-mapper` as an import, to make use of the
latest version of `media-mapper`, update your imports from
```go
import "github.com/rustedturnip/media-mapper"
```

to

```go
  import "codeberg.org/rustedturnip/media-mapper"
```

followed by running

```shell
go mod tidy
```

to remove the now-redundant dependency.

---

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
