go-imgurroulette
=======
> yep, it's basically what it sounds like

This is a roulette server for imgur. By default it hosts itself on port 8080 over
plain HTTP at localhost. If you want to turn this into a public-facing server, beware that
you're not using an API of any kind, it's completely unauthenticated. This means
that imgur might get mad. Since it's for personal use anyways, I don't care.
You should also probably put it behind nginx and wrap everything in TLS or something.

## Installation
- Set up a `$GOPATH` if you don't have one already. (`mkdir ~/.go`, `export GOPATH=~./go`, the latter preferably in your shell .rc)
- `go get gitlab.com/Niesch/go-imgurroulette`
- The binary is now installed in `~/.go/bin/`. Add it to your path, move it somewhere fun, I don't really care. Bare in mind it's using relative paths for assets and static, so you probably want to run it from its workspace (`~/.go/src/gitlab.com/Niesch/go-imgurroulette`) anyways, or modify the source to your liking.

## Running
- See `go-imgurroulette --help` or `go-imgurroulette -h` for details
- The `--workers` flag may increase performance massively. Higher values should be used with a larger `--cachesize`.
- Generally, the `--workers` flag should be increased with multiple clients viewing. The `--cachesize` flag should generally be increased with faster viewers.
- Imgur may ban you if you excessively use their services, it's probably in their ToS somewhere. I dunno, I don't read that shit.

## Usage
- Visit http://localhost:8080 in your browser. No, really, it's that simple.

## Security
Even though the local webserver is running on plain HTTP, all imgur connections are HTTPS. The served web page embeds a random imgur image over HTTPS as a png.