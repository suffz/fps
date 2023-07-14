package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

type PowerShell struct {
	powerShell string
}

var WIN_CREATE_SHORTCUT = `$WshShell = New-Object -comObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("$HOME\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\main.lnk")
$Shortcut.TargetPath = "PLACEHOLDER"
$Shortcut.Save()`

// New create new session
func New() *PowerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		powerShell: ps,
	}
}

func (p *PowerShell) execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.powerShell, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}

// enableAutostartWin creates a shortcut to MyAPP in the shell:startup folder
func enableAutostartWin() {
	ps := New()
	fff, _ := os.Executable()
	exec_path := fff
	WIN_CREATE_SHORTCUT = strings.Replace(WIN_CREATE_SHORTCUT, "PLACEHOLDER", exec_path, 1)
	stdOut, stdErr, err := ps.execute(WIN_CREATE_SHORTCUT)
	log.Printf("CreateShortcut:\nStdOut : '%s'\nStdErr: '%s'\nErr: %s",
		strings.TrimSpace(stdOut), stdErr, err)
}
