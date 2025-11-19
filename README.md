Gator is a simple CLI RSS reader that uses PostgreSQL to store feeds into a database.

## Requirements
- Go 1.21+ installed
- PostgreSQL 16+ running and accessible

## Setup
1. Create a databse user in Postgres (or use an existing one)

2. Create the database for gator:
```bash
createdb gator
```

3. Create a .gatorconfig.json file in your home directory with the following body:
```json
{
    "db_url":"postgres://USERNAME:PASSWORD@localhost:5432/gator?sslmode=disable"
}
```

4. Run the schema SQL files on the database using these commands:
```bash
psql -d gator -f sql/schema/001_users.sql
psql -d gator -f sql/schema/002_feeds.sql
psql -d gator -f sql/schema/003_feed_follows.sql
psql -d gator -f sql/schema/004_feeds.sql
psql -d gator -f sql/schema/005_posts.sql
```
5. Install the binaries from the root:
```bash
go install .
```

## Usage
From any directory once installed you can enter the commands:
```bash
gator register <username>
gator login <username>
gator users
gator reset
gator feeds
gator addfeed <name> <url>
gator follow <url>
gator following
gator unfollow <url>
gator agg <optional refresh time e.g 30s>
gator browse
```

`register` - Register a new user.
`login` - Log in with a user.
`users` - Show the list of registered users.
`reset` - Delete all users from the database.
`feeds` - Show the list of all feeds in the database.
`addfeed` - Add a feed to the database.
`follow` - Follow a particular feed for the current user.
`following` - Show currently followed feeds for the current user.
`unfollow` - Remove a feed from current users follow list.
`agg` - Fetch the feeds at an optionally specified refresh rate.
`browse` - Read the user's followed feeds.