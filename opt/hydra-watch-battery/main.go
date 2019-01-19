package main

import "time"
import "os"
import "github.com/xshellinc/tools/lib/inotify"

var battery = []byte("battery\n")
var files = []string{
	"/sys/class/power_supply/BAT0/charge_now",
	"/sys/class/power_supply/AC/uevent",
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	watcher, err := inotify.New()
	must(err)
	for _, file := range files {
		_, err := watcher.Add(file, inotify.InAccess|inotify.InCloseWrite)
		must(err)
	}
	defer watcher.Close()
	t0 := time.Now()
	for {
		<-watcher.C
		t := time.Now()
		if t.Sub(t0).Seconds() > 0.01 {
			_, _ = os.Stdout.Write(battery)
			t0 = t
		}
	}
}
