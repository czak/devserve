# devserve: A minimal webserver with livereload

`devserve` is a small development webserver with built-in live-reloading for HTML files.

## Features

- Serve static files from specified directory (or CWD by default)
- Watch files for changes
- Automatically reload HTML files

## How does it work?

HTML files are automatically injected with this snippet just before `</body>`:

```js
<script>
new EventSource('/events').addEventListener('change', () => location.reload())
</script>
```

The server recursively watches files in specified directory and notifies of any writes via `text/event-stream` events on the `/events` endpoint.

Any other files are served by `http.ServeFile` directly.

## Installation

```sh
$ go install github.com/czak/devserve@latest
```

## Usage

By default `devserve` will serve files from current directory on `:8080`:

```sh
$ devserve
2025/03/27 16:00:34 Serving files from "." on ":8080"
```

Directory and address can be overridden with `-dir` and `-addr` respectively:

```sh
$ devserve -dir public -addr "127.0.0.1:3000"
2025/03/27 16:01:25 Serving files from "public" on "127.0.0.1:3000"
```

All flags:

```sh
$ devserve -h
Usage of devserve:
  -addr string
        network address (default ":8080")
  -dir string
        directory to serve from (default ".")
```

## Notes

* Should work on all Go platforms, but only tested on Linux
