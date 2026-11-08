> # Current status: BROKEN

# Naka

![License: MIT](https://img.shields.io/github/license/rl404/naka)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/rl404/naka)
[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/rl404/naka)](https://hub.docker.com/r/rl404/naka)
[![CI](https://github.com/rl404/naka/actions/workflows/ci.yml/badge.svg?event=push)](https://github.com/rl404/naka/actions/workflows/ci.yml)

Naka is a discord bot to play song from youtube.

[Naka](https://en.wikipedia.org/wiki/Japanese_cruiser_Naka)'s name is taken from japanese cuiser. Also, [exists](https://kancolle.fandom.com/wiki/Naka) in Kantai Collection games and anime as a fleet's idol. To live up its name, this bot will 'sing' for you like an idol.

## Features

- Play song from youtube url.
- Search song from youtube.
- Song queue system.
  - Pause song.
  - Resume song.
  - Stop song.
  - Next song.
  - Previous song.
  - Skip/jump to specific song in queue.
  - Remove 1 or more song from queue.
  - Delete queue.

## Requirement

- [Discord bot](https://discordpy.readthedocs.io/en/latest/discord.html) and its token
- [Go](https://golang.org/)
- [Youtube API key](https://developers.google.com/youtube/v3/getting-started)
- [Redis](https://redis.io/) (optional)
- [Docker](https://docker.com) + [Docker compose](https://docs.docker.com/compose/) (optional)

## Steps

1. Git clone this repo.

   ```bash
   git clone github.com/rl404/naka
   ```

2. Rename `.env.sample` to `.env` and modify according to your configuration.

   | Env                    |  Default   | Description                       |
   | ---------------------- | :--------: | --------------------------------- |
   | `NAKA_DISCORD_TOKEN`\* |            | Discord bot token.                |
   | `NAKA_DISCORD_PREFIX`  |    `=`     | Discord bot prefix command.       |
   | `NAKA_CACHE_DIALECT`   | `inmemory` | Cache type (`inmemory`, `redis`). |
   | `NAKA_CACHE_ADDRESS`   |            | Cache address.                    |
   | `NAKA_CACHE_PASSWORD`  |            | Cache password.                   |
   | `NAKA_CACHE_TIME`      |   `24h`    | Cache duration.                   |
   | `NAKA_YOUTUBE_KEY`\*   |            | Youtube API key.                  |
   | `NAKA_LOG_LEVEL`       |    `-1`    | Log level                         |
   | `NAKA_LOG_JSON`        |  `false`   | Log json format                   |
   | `NAKA_LOG_COLOR`       |   `true`   | Log color                         |

3. Run.

   ```bash
   make bot

   # or using docker
   make docker-bot
   # to stop docker
   make docker-stop
   ```

4. Invite the bot to your server.
5. Join a voice channel.
6. Try `=help`.
7. Have fun.

## Bot Commands

### Play Song

```bash
# Play song in queue.
=play

# Will put the song to queue and play.
=play https://www.youtube.com/watch?v=dQw4w9WgXcQ

# Will search, put the song to queue, and play.
=play suisei stellar
```

### Queue Song

```bash
# See queue list.
=queue

# Will put the song to queue.
=queue https://www.youtube.com/watch?v=dQw4w9WgXcQ

# Will search and put the song to queue.
=queue suisei stellar
```

### Other Commands

```bash
=join       # Bot join to voice channel.
=leave      # Bot leave voice channel.
=pause      # Pause song.
=resume     # Resume song.
=stop       # Stop song.
=next       # Go to next song in queue.
=prev       # Go to previous song in queue.
=skip 2     # Go to song number 2 in queue.
=remove 1 2 # Remove song number 1 and 2 from queue.
=purge      # Remove all songs from queue.
```

## License

MIT License

Copyright (c) 2022 Axel
