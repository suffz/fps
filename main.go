package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/martinlindhe/notify"
)

var Home string = `%v\AppData\Local\Roblox\Versions`
var Homedir string
var CurrentVersion string

func init() {
	if name, err := user.Current(); err == nil {
		Homedir = name.HomeDir
		Home = fmt.Sprintf(Home, name.HomeDir)
	}

	Checkstartup()

	notify.Notify("FPS Unlocker", "Hello!", "I really hope your day is good!", "")
}

func main() {
	if Logs := Returnlogs(); len(Logs) > 0 || len(Logs) == 1 {
		if Logs[0].IsDir() {
			CurrentVersion = Logs[0].Name()
			Apply(Home + `\` + CurrentVersion + `\ClientSettings`)
		}
	}
	CheckConsistent()
}

func CheckConsistent() {
	for {
		time.Sleep(time.Minute)
		if Logs := Returnlogs(); len(Logs) > 0 || len(Logs) == 1 {
			if Logs[0].IsDir() && CurrentVersion != Logs[0].Name() {
				CurrentVersion = Logs[0].Name()
				Apply(Home + `\` + CurrentVersion + `\ClientSettings`)
			}
		}
	}
}

func Apply(Dir string) {
	if err := os.Mkdir(Dir, 0655); err == nil {
		if file, err := os.Create(Dir + `\ClientAppSettings.json`); err == nil {
			file.WriteString(`{"DFIntTaskSchedulerTargetFps":999}`)
			file.Close()
			notify.Notify("FPS Unlocker", "Applied! ᕕ( ᐛ )ᕗ", fmt.Sprintf(`Succesfully applied fps unlocker to version %v`, CurrentVersion), "")
		} else {
			notify.Notify("FPS Unlocker", "Error", err.Error(), "")
		}
	}
}

func Returnlogs() (Logs []fs.FileInfo) {
	ff, _ := os.ReadDir(Home)
	for _, f := range ff {
		if info, err := f.Info(); err == nil {
			Logs = append(Logs, info)
		}
	}
	sort.Slice(Logs, func(i, j int) bool {
		return Logs[i].ModTime().After(Logs[j].ModTime())
	})
	return
}

func Checkstartup() {
	var Dir string = Homedir + `\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`
	ff, _ := os.ReadDir(Dir)
	var found bool
	for _, f := range ff {
		fff, _ := os.Executable()
		if f.Name() == strings.Replace(filepath.Base(fff), ".exe", ".lnk", -1) {
			found = true
		}
	}
	if !found {
		enableAutostartWin()
	}
}
