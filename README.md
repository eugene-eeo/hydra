# hydra

Hydra is a simple daemon that runs in the background and emits events to clients.
It is an overkill (some would say un-UNIX-y) solution/alternative to polling in panel scripts.

## Install

```sh
$ go install github.com/eugene-eeo/hydra
$ go install github.com/eugene-eeo/hydra/opt/hydra-head
$ echo '{}' > ~/.hydrarc.json
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
To add support for other services, you can add more to the 'procs' key in the root config object:

```json
{
  ...
  "procs": [
    {
      "proc": ["herbstluftwm", "--idle"],
      "matchers": [
        {"name": "hc:focus",      "matcher": "^focus_changed"},
        {"name": "hc:tag_change", "matcher": "^tag_changed"},
      ]
    }
  ]
}
```

upon seeing lines that match ^focus\_changed, hydra will emit the hc:focus event,
and similarly for the hc:tag\_change event.  A more involved example for a matcher:

```json
{
    "&&": [
        {"||": ["^abc", "def$"]},
        "[0-9]",
        "tag"
    ]
}
```

this matcher matches all lines that match `[0-9]`, `tag`, and (`^abc` or `def$`).
The matchers are nestable and very powerful.
