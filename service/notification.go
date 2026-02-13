package service

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/wibecoderr/Reminder-2.git/model"
)

type Notifier struct{}

func NewNotifier() *Notifier {
	return &Notifier{}
}
func (n *Notifier) Notify(reminder model.Reminder) {
	title := "Reminder"
	message := reminder.Message
	switch runtime.GOOS {
	case "darwin":
		n.NotifyMacOs(title, message)

	case "windows":
		n.NotifyWindows(title, message)
	case "linux":
		n.NotifyLinux(title, message)
	}
	n.PlaySound()

}
func (n *Notifier) NotifyMacOs(title string, message string) {
	cmd := exec.Command("osascript", "-e",
		fmt.Sprintf(`display notification "%s" with title "%s"`, message, title))
	cmd.Run()
}

func (n *Notifier) NotifyLinux(title, message string) {
	cmd := exec.Command("notify-send", title, message)
	cmd.Run()
}

func (n *Notifier) NotifyWindows(title, message string) {
	// Using PowerShell for Windows notifications
	cmd := exec.Command("powershell",
		"-Command",
		fmt.Sprintf(`New-BurntToastNotification -Text "%s", "%s"`, title, message))
	cmd.Run()
}

func (n *Notifier) PlaySound() {
	switch runtime.GOOS {
	case "darwin":
		exec.Command("afplay", "/System/Library/Sounds/Ping.aiff").Run()
	case "linux":
		exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/complete.oga").Run()
	case "windows":
		exec.Command("powershell", "-Command", "[System.Media.SystemSounds]::Asterisk.Play()").Run()
	}

}
