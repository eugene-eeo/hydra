# Hydra

Hydra is a simple daemon that runs in the background and emits events to clients.
It was originally meant as a solution/alternative to polling in panel scripts.
By default Hydra listens on port 9900.

## Install

```sh
$ go install github.com/eugene-eeo/hydra
$ go install github.com/eugene-eeo/hydra/opt/hydra-head
$ echo '{}' > ~/.hydrarc.json
```

## Usage

```sh
$ nohup hydra &  # spawns processes and emits events
$ hydra-head     # listens to events emitted by hydra
```

## Config

Hydra is configured by changing the `~/.hydrarc.json` file.

```json
{
  "nmcli": true,
  "pactl": true
}
```

This listens to and outputs sensible events using `nmcli` and `pactl`.
To add support for listening to other services, you can add to the
'procs' key in the root config object:

```json
{
  "nmcli": true,
  "pactl": true
  "procs": [
    {
      "proc": ["herbstluftwm", "--idle"],
      "match": [
        ["hc:focus",      "^focus_changed"],
        ["hc:tag_change", "^tag_changed"]
      ]
    }
  ]
}
```

Think of the `match` array as a big switch statement; the first regex
that matches the line would have have it's event emitted. Thus the
order of the matchers (entries in the match array) are important.
