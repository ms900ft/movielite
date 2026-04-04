# movielite

**movielite** is a simple **personal movie database**, written in Go and Vue.js.

You can specify local directories watched for new movies. If new movies are found, metadata and images are added from TMDB. [TMDB](https://www.themoviedb.org).

![](./docs/images/screencapture-main.png)

## Features

* List your movies with metadata
* Search movies by title, description and persons
* Play movies with your local player
* View detailed movie info
* Automatically add new movies and metadata
* Uses SQLite — no DB server required
* Watchlist
* Download movies
* Move movies to desired directory
* User management with admin roles
* Multiple choice movie selection for ambiguous titles
* On-screen keyboard for search
* Genre, country, and person filtering
* Recently added sorting

## Getting Started (Demo)

- Clone the Repository
- Run in this directory: `docker-compose up`
- Goto http://localhost:8000/movie2/
- Login with admin/password

You will see some example data without media files.

## Installation

- Get the API Key from https://developers.themoviedb.org/3/

- Clone Repository: `git clone https://github.com/ms900ft/movielite.git`

### Build

```bash
# install npm packages
make npminstall
# build single binary (includes frontend)
make
```

You find the binary `movielite` in the root directory.

### Frontend Development

For frontend development, use the Vite dev server with API proxy:

```bash
cd new-ui
npm run dev
```

The dev server proxies `/api`, `/images`, and `/login` to `http://localhost:8001`.

## Configuration

The main configuration file is `movielite.yaml`, but it is not included by default. Simply make a copy of `movielite.yaml.tmpl` and rename it to `movielite.yaml`. See `example_movielite.yaml` for more information about configurations.

```yaml
# devel, prod
Mode: prod
# Port to run
Port: 8000
# Enable sql debugging
SQLDebug: false
# Language to get Metadata
Language: "en-US"
# Location of the sqlite database
Database:
  Dbname: "./movielite.db"
# Secret to sign token
Secret: "dsdhsjhsdhr8q73z478274z3qhui4"
# Initial Admin password, can be changed later
InitialAdminPassword: "password"
# Player to play movies default is vlc (optional)
# Player: "QuickTime Player"

TMDB:
  # Your API key. Without API key you can't get metadata from tmdb
  ApiKey: #add your api key
  # Directory for the image cache.
  ImageDir: ./images

# Base directory of your movie lib.
TargetDirectory: ./movies

# URL of this server
MovieServerUrl: http://localhost:8000

# Use system trash instead of delete
TrashCan: true

# Watch here for new movies
WatchDir: "./movies/new"

# Additional regex's to find movie name in filename
# Key is only for debugging. Build in can be found in
# models/files.go
FilterRegEx:
  OTR: ^(.+?)_\d{2}\.\d{2}\.\d{2}_\d{2}-\d{2}_
  MT: ^.+?-(.+?)-\d+-\d+
  MT_Zdf: ^.+?-(.+?)_-
  MT_ZDF_SEASON: ^.+?-(.+?)_\((\d+)\)-\d+_\w+_\d+

# enable webdav (admin/test123)
WebDav: true
```

### Start server

```bash
./movielite start
```

Now browse to the app at http://localhost:8000/movie2/.
Login with admin/password.

If you use the default config a watcher is set on the directory **./example/movies**.
To add a new movie, copy the movie into the watched folder.

### Add your library

The server must be running. To add your existing movies to movielite type:

```bash
movielite scan -d [directory to scan]
```

## API

The API is documented with Swagger. When the server is running, browse to http://localhost:8000/swagger/ for the interactive API documentation.
