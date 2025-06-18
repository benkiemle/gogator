# GoGator

gator is a go cli for aggregating rss feeds, allowing users to follow certain rss feeds, and get information from those feeds as new items are published

## Installation

gogator requires both Postgres and Go installed locally

After cloning the repository, navigate to the installation directory and run the following:

```cmd
go install
```

The program expects a configuration file to exist in your home directory called `.gatorconfig.json`. This file should look like the following:

```json
{
    "current_user_name":"",
    "connection_string":"POSTGRES CONNECTION STRING"
}
```

## Usage

The following commands are available from the command line. The format of commands are: `gogator [command_name] [arguments..]`

| Command Name | Arguments | Desciption |
| ------------ | --------- | ---------- |
| login | `username` | Logs the given user in as the current user, if that user exists |
| register | `username` | Creates a new user with the given username, and logs that user in |
| users | - | Lists the current users in gogator |
| agg | `duration` | Starts the rss feed aggregator, that will run in the background indefinitely. The duration given must be of the format: `\dh\dm\ds`, ex: `30s` would run every 30 seconds. The user must use `Ctrl+C` to force the process to stop. |
| addfeed | `name` `url` | Adds a new feed to the system, with the given `name` and `url` |
| feeds | - | Lists all feeds available in the system |
| follow | `url` | Follows the given rss feed url |
| following | - | Lists all feeds being followed by the current user |
| unfollow | `url` | Stops following the given rss feed url |
| browse | OPTIONAL `limit` | Displays to most recent rss feed items the current user is subscribed to, up to the `limit` (default is 2) |

