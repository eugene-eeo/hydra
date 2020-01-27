package main

import "runtime"
import "os"
import "github.com/xshellinc/tools/lib/inotify"

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
	runtime.GOMAXPROCS(1)
	watcher, err := inotify.New()
	must(err)
	for _, file := range files {
		_, err := watcher.Add(file, inotify.InAccess|inotify.InModify)
		must(err)
	}
	defer watcher.Close()
	for {
		<-watcher.C
		_, _ = os.Stdout.Write([]byte("battery\n"))
	}
}
