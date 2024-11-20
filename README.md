# bootdev_blog_aggregator
a blog aggregator written in go

This is a lightly tested implementation of the RSS feed aggregator server/CLI gator thought up by the folks over at Boot.dev. With it, you can:

- Register multiple user accounts
- Add feeds
- Follow feeds across multiple accounts
- View available feeds to follow
- Browse your most recent X aggregated posts (with a customizable value for X)
- Unfollow feeds

# Requirements

gator requires Postgres 16.x (written with 16.5) and Golang 1.23.x (written with 1.23.3) or higher.

To check if you have Postgres installed, run `psql --version` in your command line. Your result should appear similar to the following:

```bash
bash-3.2$ psql --version
psql (PostgreSQL) 16.5 (Homebrew)
```

To check if you have Golang installed, perform the above with `go version` and check for something like the following:

```bash
bash-3.2$ go version
go version go1.23.3 darwin/arm64
```

# Installation

To install, run `go install https://github.com/0x4D5352/bootdev_blog_aggregator`. That's it!

# Configuration

1. Edit the `.gatorconfig.json` file, replacing `REPLACE_WITH_YOUR_USERNAME` with your username. e.g. if your account name is `johncena`, use `postgres://johncena:@localhost:5432/gator?sslmode=disable`
2. Copy the file to your home directory, typically accessible at the alias `~` or at the location specified by the `$XDG_CONFIG_HOME` environment variable.
3. Possibly mess with the postgres DB? i don't understand how y'all are supposed to get the DB when you install the binary directly. 

# Available Commands

> [NOTE] Note:
> All commands are to be prefixed with the binary, e.g. `bootdev_blog_aggregator browse 1`. Mandatory values are surrounded by square brackets `[]`, optional values by parentheses `()`.

- `login [username]` -  set the speciied user as the current user.
- `register [username]` - register a new user and set them as the current user.
- `users` - list all registered users, with the current user indicated.
- `addfeed [Feed Name] [Feed URL]` - add the RSS feed at the given URL to the available feeds for all users and follow it.
- `agg [time]` - begin a process to fetch RSS feed posts from all available feeds, checking a new feed every `[time]` interval. Use duration strings such as "1s", "30m", "24h", etc.
- `feeds` - list all available feeds 
- `follow [Feed URL]` - follow a feed added by another user, specified by the URL.
- `following` - list all currently followed feeds
- `unfollow [Feed URL]` - unfollow a feed, specified by the URL
- `browse (limit)` - view up to `(limit)` posts, aggregated across your feeds. when unspecified, limit defaults to 2.
