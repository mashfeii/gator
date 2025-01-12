# Gator 🐊

![Rss](https://img.shields.io/badge/rss-F88900?style=for-the-badge&logo=rss&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)

Foobar is a small RSS aggregator written in Go using Postgres Database.

## Installation

To run the program, you need to install [Postgres](https://www.postgresql.org/download/) and [Golang](https://go.dev/doc/install).

Then, to install application generally, run:

```bash
go install github.com/mashfeii/gator
```

Program will search for config file inside `$HOME` directory by `.gatorconfig.json` as default file name, there you should specify `db_url` field.

## Usage

```bash
# Clear users table (other tables will be cleared as well)
gator reset

# Register user (automatically become logged in)
gator register masfheii
gator register byaka

# Login user (in order to switch them)
gator login mashfeii

# List all the users
gator users

# Add new feed to current user
gator addfeed '$name' '$URL'

# List all the feeds for all users
gator feeds

# List feeds for the current user
gator following

# Remove feed from current user
gator unfollow '$URL'

# With selected duration check the source for new posts
gator agg '$duration (5s/1m/10h)'

# Show selected number of posts from database for current user
gator browse $limit
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
