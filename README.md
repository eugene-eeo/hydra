<p align='center'><img src='logo/hydra.png' height='250px'/></p>

*hydra* is a simple daemon that emits events to clients listening over TCP.
It was originally meant as (and is still used everyday as) a more efficient
alternative to polling in panel scripts. *hydra* listens on port 9900.

*hydra-head* connects to *hydra*'s socket and listens for emitted events.


## Install

    $ go install github.com/eugene-eeo/hydra
    $ go install github.com/eugene-eeo/hydra/opt/hydra-head
    $ echo '{}' > ~/.hydrarc.json


## Usage

    $ nohup hydra &  # spawns processes and emits events
    $ hydra-head     # listens to events emitted by hydra
    $ nc localhost 9900  # alternative to hydra-head


## Config

Options for *hydra* are stored in the `~/.hydrarc.json` file. Note that the
daemon needs to be restarted for changes to take effect (you can just kill
the hydra process). An example config file is:

    {
      "nmcli": true,
      "pactl": true
    }

This listens to NetworkManager and PulseAudio events. The `nmcli` and
`pactl` events are emitted for each of those services respectively.
To listen to services, add more declarations to the `procs` key. Say you
wanted to send events whenever the networkmanager status changes.
We can match and emit on the relevant lines:

    {
      "procs": [{
        "proc": "nmcli monitor",
        "match": [
          ["nmcli:connected",    "^.+: connected$"],
          ["nmcli:disconnected", "^.+: disconnected$"],
          ["nmcli:unavailable",  "^.+: unavailable$"]
        ]
      }]
    }

Think of the `match` array as a big switch statement; each element in
the array (the matcher) contains an event and regex. The output of the
process is processed line-by-line; The first matcher whice regex matches
the line has its event emitted. Thus the order of the matchers is important.
If the `match` array is empty, then all output from the process is forwarded
line by line to all clients.

Then we can have a little bash script that runs in the background
and monitors these events:

    #!/bin/bash
    hydra-head | while read event; do
        case "$event" in
            nmcli:connected)    notify-send 'nm' "connected to $(iwgetid -r)" ;;
            nmcli:disconnected) notify-send 'nm' 'disconnected' ;;
            nmcli:unavailable)  notify-send 'nm' 'unavailable' ;;
        esac
    done

The advantage for using hydra becomes apparent when you have multiple
of these monitor scripts; you don't have to spawn multiple instances
of `nmcli monitor`.

## Tips

Even more dynamicism is possible using the `hydra-emit` tool:

    $ go install github.com/eugene-eeo/hydra/opt/hydra-emit
    $ cat ~/.hydrarc.json
    {
      "procs": [
        {"proc": ["hydra-emit", "server"]}
      ]
    }

Then you can emit events to all listening clients using:

    $ hydra-emit emit 'abc'

