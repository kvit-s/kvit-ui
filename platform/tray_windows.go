package platform

// The notification area on Windows: Shell_NotifyIconW, with a hidden window
// of the tray's own that receives the icon's messages. The window lives on a
// goroutine locked to its own OS thread, with its own message loop, so the
// tray works whatever unison's thread is doing; the UI thread hands it work
// through a queue and a posted message (do). The menu is a popup menu shown
// with TrackPopupMenuEx, and a notification is the icon's balloon (NIF_INFO),
// which Windows 10 and 11 show as a toast; clicking it sends
// NIN_BALLOONUSERCLICK. The icon is made from the picture's pixels at the
// size Windows draws small icons (SM_CXSMICON).

import (
	"image"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	shell32  = windows.NewLazySystemDLL("shell32.dll")
	gdi32    = windows.NewLazySystemDLL("gdi32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")

	procShellNotifyIcon       = shell32.NewProc("Shell_NotifyIconW")
	procGetModuleHandle       = kernel32.NewProc("GetModuleHandleW")
	procRegisterClassEx       = user32.NewProc("RegisterClassExW")
	procCreateWindowEx        = user32.NewProc("CreateWindowExW")
	procDefWindowProc         = user32.NewProc("DefWindowProcW")
	procDestroyWindow         = user32.NewProc("DestroyWindow")
	procGetMessage            = user32.NewProc("GetMessageW")
	procTranslateMessage      = user32.NewProc("TranslateMessage")
	procDispatchMessage       = user32.NewProc("DispatchMessageW")
	procPostMessage           = user32.NewProc("PostMessageW")
	procPostQuitMessage       = user32.NewProc("PostQuitMessage")
	procRegisterWindowMessage = user32.NewProc("RegisterWindowMessageW")
	procCreatePopupMenu       = user32.NewProc("CreatePopupMenu")
	procAppendMenu            = user32.NewProc("AppendMenuW")
	procTrackPopupMenuEx      = user32.NewProc("TrackPopupMenuEx")
	procDestroyMenu           = user32.NewProc("DestroyMenu")
	procSetForegroundWindow   = user32.NewProc("SetForegroundWindow")
	procGetCursorPos          = user32.NewProc("GetCursorPos")
	procGetSystemMetrics      = user32.NewProc("GetSystemMetrics")
	procCreateIconIndirect    = user32.NewProc("CreateIconIndirect")
	procDestroyIcon           = user32.NewProc("DestroyIcon")
	procSetTimer              = user32.NewProc("SetTimer")
	procKillTimer             = user32.NewProc("KillTimer")
	procGetDC                 = user32.NewProc("GetDC")
	procReleaseDC             = user32.NewProc("ReleaseDC")
	procCreateDIBSection      = gdi32.NewProc("CreateDIBSection")
	procCreateBitmap          = gdi32.NewProc("CreateBitmap")
	procDeleteObject          = gdi32.NewProc("DeleteObject")
)

const (
	nimAdd        = 0
	nimModify     = 1
	nimDelete     = 2
	nimSetVersion = 4

	nifMessage = 0x01
	nifIcon    = 0x02
	nifTip     = 0x04
	nifInfo    = 0x10
	nifShowTip = 0x80

	niifInfo           = 0x01
	notifyIconVersion4 = 4

	wmNull        = 0x0000
	wmDestroy     = 0x0002
	wmContextMenu = 0x007B
	wmTimer       = 0x0113
	wmUser        = 0x0400
	wmApp         = 0x8000

	// The icon's messages, and the queue's.
	wmTrayIcon = wmApp + 1
	wmTrayRun  = wmApp + 2

	ninSelect           = wmUser + 0
	ninKeySelect        = wmUser + 1
	ninBalloonUserClick = wmUser + 5

	wsPopup        = 0x80000000
	wsExToolWindow = 0x00000080

	mfString    = 0x0000
	mfGrayed    = 0x0001
	mfSeparator = 0x0800

	tpmRightButton = 0x0002
	tpmRightAlign  = 0x0008
	tpmBottomAlign = 0x0020
	tpmNoNotify    = 0x0080
	tpmReturnCmd   = 0x0100

	smMenuDropAlignment = 40
	smCxSmIcon          = 49

	biBitfields  = 3
	dibRGBColors = 0

	// trayIconID is the icon's number among this window's icons; there is one.
	trayIconID = 1
	// balloonTimer tells a balloon's click apart from a click on the icon
	// made while the balloon shows, which Windows also reports as the
	// balloon's;  waits the same way.
	balloonTimer = 1
	balloonWait  = 80 * time.Millisecond
)

// trayWindowClass is the class of the tray's hidden window, which a script
// can find it by.
const trayWindowClass = "KvitTrayWindow"

// notifyIconData is NOTIFYICONDATAW, 976 bytes on 64-bit Windows.
type notifyIconData struct {
	Size             uint32
	Wnd              uintptr
	ID               uint32
	Flags            uint32
	CallbackMessage  uint32
	Icon             uintptr
	Tip              [128]uint16
	State            uint32
	StateMask        uint32
	Info             [256]uint16
	TimeoutOrVersion uint32
	InfoTitle        [64]uint16
	InfoFlags        uint32
	GUIDItem         windows.GUID
	BalloonIcon      uintptr
}

type wndClassEx struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   uintptr
	Icon       uintptr
	Cursor     uintptr
	Background uintptr
	MenuName   *uint16
	ClassName  *uint16
	IconSm     uintptr
}

type winMsg struct {
	Wnd      uintptr
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       winPoint
	LPrivate uint32
}

type winPoint struct{ X, Y int32 }

// bitmapV5Header is BITMAPV5HEADER, which can say the fourth byte of each
// pixel is alpha.
type bitmapV5Header struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
	RedMask       uint32
	GreenMask     uint32
	BlueMask      uint32
	AlphaMask     uint32
	CSType        uint32
	Endpoints     [9]int32
	GammaRed      uint32
	GammaGreen    uint32
	GammaBlue     uint32
	Intent        uint32
	ProfileData   uint32
	ProfileSize   uint32
	Reserved      uint32
}

type iconInfo struct {
	Icon     int32
	XHotspot uint32
	YHotspot uint32
	Mask     uintptr
	Color    uintptr
}

var (
	trayClassOnce  sync.Once
	trayClassErr   error
	taskbarCreated uint32   // the message Explorer broadcasts when it starts again
	trayWindows    sync.Map // hidden window to its *winTray
)

// winTray is the notification area on Windows.
type winTray struct {
	events trayEvents
	hwnd   uintptr       // the hidden window; set before newTrayBackend returns
	done   chan struct{} // closed when the tray's thread has ended

	mu    sync.Mutex
	queue []func()

	// What follows belongs to the tray's thread.
	added     bool      // the icon is in the notification area
	state     trayState // what it shows
	icon      uintptr
	balloon   string    // the ID of the notification last shown
	pending   string    // a balloon click waiting out balloonWait
	clickedAt time.Time // the last click on the icon
}

// newTrayBackend starts the tray's thread and waits for its window.
func newTrayBackend(events trayEvents) trayBackend {
	t := &winTray{events: events, done: make(chan struct{})}
	ready := make(chan struct{})
	go t.run(ready)
	<-ready
	if t.hwnd == 0 {
		return noTray{}
	}
	return t
}

// run is the tray's thread: its window and its message loop.
func (t *winTray) run(ready chan struct{}) {
	runtime.LockOSThread()
	defer close(t.done)
	hwnd := createTrayWindow()
	if hwnd == 0 {
		close(ready)
		return
	}
	trayWindows.Store(hwnd, t)
	defer trayWindows.Delete(hwnd)
	t.hwnd = hwnd
	close(ready)
	var m winMsg
	for {
		if r, _, _ := procGetMessage.Call(uintptr(unsafe.Pointer(&m)), 0, 0, 0); int32(r) <= 0 {
			return
		}
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&m)))
	}
}

func utf16Ptr(s string) *uint16 {
	p, _ := windows.UTF16PtrFromString(s)
	return p
}

// createTrayWindow makes a hidden window, registering its class the first
// time. It is a top-level window rather than a message-only one because the
// menu needs a window that can be put in front.
func createTrayWindow() uintptr {
	instance, _, _ := procGetModuleHandle.Call(0)
	trayClassOnce.Do(func() {
		wc := wndClassEx{WndProc: windows.NewCallback(trayWindowProc), Instance: instance, ClassName: utf16Ptr(trayWindowClass)}
		wc.Size = uint32(unsafe.Sizeof(wc))
		if r, _, err := procRegisterClassEx.Call(uintptr(unsafe.Pointer(&wc))); r == 0 {
			trayClassErr = err
		}
		r, _, _ := procRegisterWindowMessage.Call(uintptr(unsafe.Pointer(utf16Ptr("TaskbarCreated"))))
		taskbarCreated = uint32(r)
	})
	if trayClassErr != nil {
		return 0
	}
	hwnd, _, _ := procCreateWindowEx.Call(wsExToolWindow, uintptr(unsafe.Pointer(utf16Ptr(trayWindowClass))),
		uintptr(unsafe.Pointer(utf16Ptr("Kvit tray"))), wsPopup, 0, 0, 0, 0, 0, 0, instance, 0)
	return hwnd
}

func trayWindowProc(hwnd, msg, wparam, lparam uintptr) uintptr {
	if v, ok := trayWindows.Load(hwnd); ok && v.(*winTray).handle(uint32(msg), wparam, lparam) {
		return 0
	}
	r, _, _ := procDefWindowProc.Call(hwnd, msg, wparam, lparam)
	return r
}

// handle is the tray window's messages, on the tray's thread, and reports
// whether it used the message.
func (t *winTray) handle(msg uint32, wparam, lparam uintptr) bool {
	switch {
	case msg == wmTrayRun:
		t.mu.Lock()
		queue := t.queue
		t.queue = nil
		t.mu.Unlock()
		for _, f := range queue {
			f()
		}
	case msg == wmTrayIcon:
		// With NOTIFYICON_VERSION_4 the event is in lParam's low word and
		// the point it happened at in wParam.
		switch lparam & 0xFFFF {
		case ninSelect, ninKeySelect:
			t.clickedAt = time.Now()
			t.events.clicked()
		case wmContextMenu:
			t.showMenu(int32(int16(wparam&0xFFFF)), int32(int16(wparam>>16&0xFFFF)))
		case ninBalloonUserClick:
			t.pending = t.balloon
			procSetTimer.Call(t.hwnd, balloonTimer, uintptr(balloonWait/time.Millisecond), 0)
		}
	case msg == wmTimer && wparam == balloonTimer:
		procKillTimer.Call(t.hwnd, balloonTimer)
		if time.Since(t.clickedAt) > 2*balloonWait {
			t.events.notificationClicked(t.pending)
		}
	case msg == taskbarCreated && taskbarCreated != 0:
		// Explorer started again and has forgotten every icon.
		if t.added {
			t.added = false
			t.apply()
		}
	case msg == wmDestroy:
		procPostQuitMessage.Call(0)
	default:
		return false
	}
	return true
}

// do runs f on the tray's thread.
func (t *winTray) do(f func()) {
	t.mu.Lock()
	t.queue = append(t.queue, f)
	t.mu.Unlock()
	procPostMessage.Call(t.hwnd, wmTrayRun, 0, 0)
}

func (t *winTray) available() bool { return true }

// authorization is Authorized: a notification area on Windows always takes
// balloons, and the person turns them off in Windows' own settings, which
// an app is not told about.
func (t *winTray) authorization() Authorization { return Authorized }

func (t *winTray) requestAuthorization() {}

func (t *winTray) data() notifyIconData {
	d := notifyIconData{Wnd: t.hwnd, ID: trayIconID}
	d.Size = uint32(unsafe.Sizeof(d))
	return d
}

func shellNotifyIcon(message uintptr, d *notifyIconData) bool {
	r, _, _ := procShellNotifyIcon.Call(message, uintptr(unsafe.Pointer(d)))
	return r != 0
}

// copyUTF16 writes s into a fixed field, cut to fit with its terminating
// zero.
func copyUTF16(dst []uint16, s string) {
	u := utf16.Encode([]rune(s))
	if len(u) > len(dst)-1 {
		u = u[:len(dst)-1]
	}
	n := copy(dst, u)
	dst[n] = 0
}

func (t *winTray) show(s trayState) {
	t.do(func() {
		t.state = s
		t.apply()
	})
}

// apply adds the icon, or changes it to what the state says. Tray thread.
func (t *winTray) apply() {
	d := t.data()
	d.Flags = nifMessage | nifTip | nifShowTip
	d.CallbackMessage = wmTrayIcon
	copyUTF16(d.Tip[:], t.state.tooltip)
	old := t.icon
	t.icon = makeIcon(t.state.icons)
	if t.icon != 0 {
		d.Flags |= nifIcon
		d.Icon = t.icon
	}
	if t.added {
		shellNotifyIcon(nimModify, &d)
	} else if shellNotifyIcon(nimAdd, &d) {
		t.added = true
		d.TimeoutOrVersion = notifyIconVersion4
		shellNotifyIcon(nimSetVersion, &d)
	}
	if old != 0 {
		procDestroyIcon.Call(old)
	}
}

func (t *winTray) hide() { t.do(t.remove) }

// remove takes the icon out of the notification area. Tray thread.
func (t *winTray) remove() {
	if !t.added {
		return
	}
	d := t.data()
	shellNotifyIcon(nimDelete, &d)
	t.added = false
}

// notify shows the notification as the icon's balloon. A balloon without
// text is not shown, so a notification with a title alone shows the title
// as its text.
func (t *winTray) notify(n Notification) {
	t.do(func() {
		if !t.added {
			return
		}
		title, text := n.Title, n.Message
		if text == "" {
			title, text = "", title
		}
		d := t.data()
		d.Flags = nifInfo
		d.InfoFlags = niifInfo
		copyUTF16(d.InfoTitle[:], title)
		copyUTF16(d.Info[:], text)
		t.balloon = n.ID
		shellNotifyIcon(nimModify, &d)
	})
}

// close takes the icon away and ends the tray's thread, waiting for it so
// the icon is gone before the app exits.
func (t *winTray) close() {
	t.do(func() {
		t.remove()
		if t.icon != 0 {
			procDestroyIcon.Call(t.icon)
			t.icon = 0
		}
		procDestroyWindow.Call(t.hwnd)
	})
	select {
	case <-t.done:
	case <-time.After(2 * time.Second):
	}
}

// showMenu shows the menu at a point on the screen, and runs the line
// chosen. Tray thread. The window has to be in front for the menu to close
// when the person clicks elsewhere, and the empty message after it lets the
// menu open properly the next time (TrackPopupMenu's documentation).
func (t *winTray) showMenu(x, y int32) {
	if len(t.state.menu) == 0 {
		return
	}
	menu, _, _ := procCreatePopupMenu.Call()
	if menu == 0 {
		return
	}
	defer procDestroyMenu.Call(menu)
	for i, e := range t.state.menu {
		if e.separator {
			procAppendMenu.Call(menu, mfSeparator, 0, 0)
			continue
		}
		flags := uintptr(mfString)
		if e.disabled {
			flags |= mfGrayed
		}
		// An ampersand marks an access key in a Windows menu; doubled, it
		// shows as itself.
		procAppendMenu.Call(menu, flags, uintptr(i+1), uintptr(unsafe.Pointer(utf16Ptr(strings.ReplaceAll(e.text, "&", "&&")))))
	}
	if x == 0 && y == 0 {
		var p winPoint
		procGetCursorPos.Call(uintptr(unsafe.Pointer(&p)))
		x, y = p.X, p.Y
	}
	flags := uintptr(tpmReturnCmd | tpmNoNotify | tpmRightButton | tpmBottomAlign)
	if r, _, _ := procGetSystemMetrics.Call(smMenuDropAlignment); r != 0 {
		flags |= tpmRightAlign
	}
	procSetForegroundWindow.Call(t.hwnd)
	chosen, _, _ := procTrackPopupMenuEx.Call(menu, flags, uintptr(x), uintptr(y), t.hwnd, 0)
	procPostMessage.Call(t.hwnd, wmNull, 0, 0)
	if chosen > 0 {
		t.events.chose(int(chosen) - 1)
	}
}

// makeIcon makes an icon of the picture at the size Windows draws small
// icons, from its pixels: a 32-bit colour bitmap whose fourth byte is
// alpha, and a mask of zeros, which a 32-bit icon's alpha takes the place
// of. It is 0 without a picture.
func makeIcon(icons []image.Image) uintptr {
	side, _, _ := procGetSystemMetrics.Call(smCxSmIcon)
	if side == 0 {
		side = 16
	}
	n := int(side)
	img := iconAt(icons, n)
	if img == nil {
		return 0
	}
	h := bitmapV5Header{Width: int32(n), Height: -int32(n), Planes: 1, BitCount: 32, Compression: biBitfields,
		RedMask: 0x00FF0000, GreenMask: 0x0000FF00, BlueMask: 0x000000FF, AlphaMask: 0xFF000000}
	h.Size = uint32(unsafe.Sizeof(h))
	dc, _, _ := procGetDC.Call(0)
	var bits unsafe.Pointer
	color, _, _ := procCreateDIBSection.Call(dc, uintptr(unsafe.Pointer(&h)), dibRGBColors, uintptr(unsafe.Pointer(&bits)), 0, 0)
	procReleaseDC.Call(0, dc)
	if color == 0 || bits == nil {
		return 0
	}
	defer procDeleteObject.Call(color)
	pix := unsafe.Slice((*byte)(bits), n*n*4)
	for i := 0; i < len(pix); i += 4 {
		// Top-down rows of blue, green, red and alpha, with alpha not
		// multiplied in, as in an icon file.
		pix[i], pix[i+1], pix[i+2], pix[i+3] = img.Pix[i+2], img.Pix[i+1], img.Pix[i], img.Pix[i+3]
	}
	stride := (n + 15) / 16 * 2 // a mask's rows are whole 16-bit words
	zeros := make([]byte, stride*n)
	mask, _, _ := procCreateBitmap.Call(uintptr(n), uintptr(n), 1, 1, uintptr(unsafe.Pointer(&zeros[0])))
	if mask == 0 {
		return 0
	}
	defer procDeleteObject.Call(mask)
	info := iconInfo{Icon: 1, Mask: mask, Color: color}
	icon, _, _ := procCreateIconIndirect.Call(uintptr(unsafe.Pointer(&info)))
	return icon
}
