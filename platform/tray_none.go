package platform

// noTray is the notification area of a desktop that has none: every system
// without a backend, and a Linux session without a StatusNotifierWatcher.
type noTray struct{}

func (noTray) available() bool              { return false }
func (noTray) show(trayState)               {}
func (noTray) hide()                        {}
func (noTray) notify(Notification)          {}
func (noTray) authorization() Authorization { return Unsupported }
func (noTray) requestAuthorization()        {}
func (noTray) close()                       {}
