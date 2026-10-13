# Blog Aggregator

## Overview

This command-line tool helps users collect and explore posts from RSS feeds across the internet. Its features include:

- Adding RSS feeds for the service to collect.
- Storing collected posts in a PostgreSQL database for later browsing.
- Discovering feeds added by other users and following or unfollowing them as interests change.
- Viewing summaries of aggregated posts in the terminal, with links to the original articles.

## Technology stack

- **Language:** Go
- **Database:** PostgreSQL

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/sharipovr/blog-aggregator-go.git
   cd blog-aggregator-go
   ```
2. Install application:
   ```bash
   cd blog-aggregator-go
   go install .
   mv "$(go env GOPATH)/bin/blog-aggregator-go" \
   "$(go env GOPATH)/bin/gator"
   ```

## Setup configuration and database
1. Create a configuration file (`~/.gatorconfig.json`) 

2. Open the configuration file in a text editor and update it with your database credentials i.e:
    ```json
{"db_url":"postgres://yourname:@localhost:5432/gator?sslmode=disable","current_user_name":"yourname"}
    ```

3. Using psql or any PostgreSQL client, create the database specified in your configuration file. For example:
   ```bash
   psql -U yourname -h localhost -d postgres -c "CREATE DATABASE gator;"
   ```
4.  Set the `DATABASE_URL` environment variable to point to your database.
    ``` bash 
    export DATABASE_URL="postgres://yourname:@localhost:5432/gator?sslmode=disable"
    ```
5. Create the necessary tables in the database by running the scripts provided in the `sql/schema` directory.
For example:
   ```bash
   psql "$DATABASE_URL" -f sql/schema/001_users.sql
   psql "$DATABASE_URL" -f sql/schema/002_feeds.sql
   psql "$DATABASE_URL" -f sql/schema/003_feed_follows.sql
   psql "$DATABASE_URL" -f sql/schema/004_feed_lastfetched.sql
   psql "$DATABASE_URL" -f sql/schema/005_posts.sql
   ```
  
## How to use / run
Application has several commands to interact with the blog aggregator. Some of the key commands include:
- `gator register <name>`: Register a new user with the specified name.
- `gator addfeed <feed_url>`: Add a new RSS feed to the current user's collection.
- `gator agg 30s`: Aggregate posts from all followed feeds every 30 seconds.
- `gator feeds`: List all the RSS feeds you have added.
- `gator follow <feed_url>`: Follow a feed added by another user.
- `gator unfollow <feed_url>`: Unfollow a feed you are currently following.
- `gator browse <limit>`: View posts with a specified limit.

