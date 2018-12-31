<img src='logo/hydra.png'/>

*hydra* is a simple daemon that emits events to clients listening over TCP.
It was originally meant as (and is still used as) a more efficient alternative
to polling in panel scripts. By default *hydra* listens on port 9900.


## Install

    $ go install github.com/eugene-eeo/hydra
    $ go install github.com/eugene-eeo/hydra/opt/hydra-head
    $ echo '{}' > ~/.hydrarc.json


## Usage

    $ nohup hydra &  # spawns processes and emits events
    $ hydra-head     # listens to events emitted by hydra


## Config

Options for *hydra* are stored in the `~/.hydrarc.json` file. Note that the
daemon needs to be restarted for changes to take effect. One way to restart
the deamon is:

    $ killall hydra
    $ nohup hydra &

An example config file is:

    {
      "nmcli": true,
      "pactl": true
    }

This listens to NetworkManager and PulseAudio events. The `nmcli` and
`pactl` events are emitted for each of those services respectively.
To add support for listening to other services, you can add to the
'procs' key in the root config object:

    {
      "nmcli": true,
      "pactl": true,
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

Think of the `match` array as a big switch statement; each element in
the array (the matcher) contains an event and regex. The first matcher
which regex matches the line would have have it's event emitted.
Thus the order of the matchers is important.


## Example

Say you wanted to send notifications whenever the networkmanager
status changes. We can match against the relevant lines in the
config file:

    {
      "nmcli": false,
      "procs": [{
        "proc": ["nmcli", "monitor"],
        "match": [
          ["nmcli:connected",    "^(.+): connected$"],
          ["nmcli:disconnected", "^(.+): disconnected$"],
          ["nmcli:unavailable",  "^(.+): unavailable"],
        ]
      }]
    }

Then we can have a little bash script that runs in the background
and monitors these events:

    #!/bin/bash
    hydra-head | while read event; do
        case "$event" in
            nmcli:connected)    notify-send 'nm-monitor' "connected to $(iwgetid -r)" ;;
            nmcli:disconnected) notify-send 'nm-monitor' 'disconnected' ;;
            nmcli:unavailable)  notify-send 'nm-monitor' 'unavailable' ;;
        esac
    done

The advantage for using hydra becomes apparent when you have multiple
of these monitor scripts; you don't have to spawn multiple instances
of `nmcli monitor`.
