# Hydra

Hydra is a simple daemon that runs in the background and emits events to clients.
It is an overkill (some would say un-UNIX-y) solution/alternative to polling in panel scripts.
By default Hydra listens on port 9900.

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
To add support for listening to other services, you can add to the
'procs' key in the root config object:

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

upon seeing lines that match ^focus\_changed, hydra will emit the
`hc:focus` event, and similarly for the `hc:tag_change` event.
Note that the order of the matchers are important. For instance,
the following produce different output:

```
[                                  | [
  {"name": "1", "matcher": "^a"},  |   {"name": "2", "matcher": "^ab"},
  {"name": "2", "matcher": "^ab"}, |   {"name": "1", "matcher": "^a"},
]                                  | ]
```

In the first case, the `1` event is always fired regardless of whether
the matchers for `2` match. This is intended behaviour from hydra's side.
If you want your intended behaviour to be used then the second case is
probably what you want.

A more involved example for a matcher:

```json
{
  "&&": [
    {"||": ["^abc", "def$"]},
    "[0-9]",
    "tag"
  ]
}
```

this matcher matches all lines that match `[0-9]`, `tag`, and
(`^abc` or `def$`). The matchers are nestable and very powerful.
If a matcher is null then it will always match. This is an efficient
alternative to the "" matcher.