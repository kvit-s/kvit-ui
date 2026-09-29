package platform

// The notification area on macOS: an NSStatusItem in the menu bar with an
// NSMenu, which opens on a click as menu bar items do, and
// UNUserNotificationCenter for notifications, as the app's
// systemtray_mac.mm has them. It calls Objective-C through purego's runtime
// bindings without cgo, the way unison's internal/cocoa does.
//
// This file is compiled for darwin/arm64 and darwin/amd64 on Linux and has
// not yet run on a Mac, as none was available when it was written.
//
// AppKit is used from the main thread only, which is unison's UI thread on
// macOS, where the Tray's methods are called. The notification center
// answers on its own threads; those answers reach the app through the
// Tray, which hands them to the UI thread.
//
// UNUserNotificationCenter raises an exception in a program that is not in
// an application bundle, so a program started on its own, without
// Info.plist and a bundle identifier, reports notifications as Unsupported.

import (
	"bytes"
	"image/png"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

const (
	nsUTF8StringEncoding     = 4
	nsVariableStatusItemSize = -1.0
	// statusIconPoints is the height of a menu bar icon; it is drawn from
	// a picture twice as large for Retina screens.
	statusIconPoints = 18

	unAuthorizationOptionAlert    = 1 << 2
	unPresentationOptionList      = 1 << 3
	unPresentationOptionBanner    = 1 << 4
	unDefaultActionIdentifier     = "com.apple.UNNotificationDefaultActionIdentifier"
	activationKey                 = "kvit-activation-id"
	unAuthorizationDenied         = 1
	unAuthorizationAuthorized     = 2
	unAuthorizationProvisional    = 3
	unAuthorizationEphemeral      = 4
	appKitPath                    = "/System/Library/Frameworks/AppKit.framework/AppKit"
	userNotificationsPath         = "/System/Library/Frameworks/UserNotifications.framework/UserNotifications"
	objcLibraryPath               = "/usr/lib/libobjc.A.dylib"
	trayTargetClassName           = "KvitTrayTarget"
	notificationDelegateClassName = "KvitNotificationDelegate"
)

// macSize is NSSize.
type macSize struct{ Width, Height float64 }

var (
	macClassesOnce       sync.Once
	macTarget            objc.ID // the menu items' target
	macNotificationClass objc.Class
	macActive            atomic.Pointer[macTray] // the tray the classes report to

	poolOnce sync.Once
	poolPush func() uintptr
	poolPop  func(uintptr)
)

// macTray is the notification area on macOS.
type macTray struct {
	events   trayEvents
	item     objc.ID // the NSStatusItem while shown
	menu     objc.ID // its NSMenu
	center   objc.ID // the UNUserNotificationCenter; 0 without a bundle
	delegate objc.ID
	auth     Authorization
}

func newTrayBackend(events trayEvents) trayBackend {
	if _, err := purego.Dlopen(appKitPath, purego.RTLD_LAZY|purego.RTLD_GLOBAL); err != nil {
		return noTray{}
	}
	t := &macTray{events: events, auth: Unsupported}
	macActive.Store(t)
	registerMacClasses()
	withPool(t.startNotifications)
	return t
}

func sel(name string) objc.SEL { return objc.RegisterName(name) }

func class(name string) objc.ID { return objc.ID(objc.GetClass(name)) }

// withPool runs f inside an autorelease pool of its own, on one OS thread.
func withPool(f func()) {
	poolOnce.Do(func() {
		lib, err := purego.Dlopen(objcLibraryPath, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
		if err != nil {
			return
		}
		purego.RegisterLibFunc(&poolPush, lib, "objc_autoreleasePoolPush")
		purego.RegisterLibFunc(&poolPop, lib, "objc_autoreleasePoolPop")
	})
	if poolPush == nil {
		f()
		return
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	pool := poolPush()
	defer poolPop(pool)
	f()
}

// nsString returns an autoreleased NSString; call it inside withPool.
func nsString(s string) objc.ID {
	str := class("NSString").Send(sel("alloc")).Send(sel("initWithBytes:length:encoding:"), s, uint64(len(s)), uint64(nsUTF8StringEncoding))
	return str.Send(sel("autorelease"))
}

func goString(str objc.ID) string {
	if str == 0 {
		return ""
	}
	n := objc.Send[uint64](str, sel("lengthOfBytesUsingEncoding:"), uint64(nsUTF8StringEncoding))
	p := objc.Send[*byte](str, sel("UTF8String"))
	if n == 0 || p == nil {
		return ""
	}
	return string(unsafe.Slice(p, n))
}

// callBlock calls a block the system handed over, through the function
// pointer at offset 16 of every block, after its class and flags.
func callBlock(b objc.Block, args ...uintptr) {
	if b == 0 {
		return
	}
	// The block's value as a pointer, taken without a uintptr conversion.
	p := *(*unsafe.Pointer)(unsafe.Pointer(&b))
	invoke := *(*uintptr)(unsafe.Add(p, 16))
	purego.SyscallN(invoke, append([]uintptr{uintptr(b)}, args...)...)
}

// registerMacClasses registers the menu items' target, whose action runs
// the line with the item's tag, and the notification center's delegate.
func registerMacClasses() {
	macClassesOnce.Do(func() {
		target, err := objc.RegisterClass(trayTargetClassName, objc.GetClass("NSObject"), nil, nil, []objc.MethodDef{{
			Cmd: sel("kvitTrayChoose:"),
			Fn: func(_ objc.ID, _ objc.SEL, sender objc.ID) {
				if t := macActive.Load(); t != nil {
					t.events.chose(int(objc.Send[int64](sender, sel("tag"))))
				}
			},
		}})
		if err == nil {
			macTarget = objc.ID(target).Send(sel("new"))
		}
		var protocols []*objc.Protocol
		if p := objc.GetProtocol("UNUserNotificationCenterDelegate"); p != nil {
			protocols = append(protocols, p)
		}
		delegate, err := objc.RegisterClass(notificationDelegateClassName, objc.GetClass("NSObject"), protocols, nil, []objc.MethodDef{
			{
				// A notification arriving while the app is in front is
				// shown all the same, as a banner.
				Cmd: sel("userNotificationCenter:willPresentNotification:withCompletionHandler:"),
				Fn: func(_ objc.ID, _ objc.SEL, _, _ objc.ID, completion objc.Block) {
					callBlock(completion, unPresentationOptionBanner|unPresentationOptionList)
				},
			},
			{
				Cmd: sel("userNotificationCenter:didReceiveNotificationResponse:withCompletionHandler:"),
				Fn: func(_ objc.ID, _ objc.SEL, _, response objc.ID, completion objc.Block) {
					withPool(func() { reportResponse(response) })
					callBlock(completion)
				},
			},
		})
		if err == nil {
			macNotificationClass = delegate
		}
	})
}

// reportResponse passes on a click on a notification's body, with the ID
// the notification was posted with.
func reportResponse(response objc.ID) {
	t := macActive.Load()
	if t == nil || response == 0 {
		return
	}
	action := response.Send(sel("actionIdentifier"))
	if !objc.Send[bool](action, sel("isEqualToString:"), nsString(unDefaultActionIdentifier)) {
		return
	}
	info := response.Send(sel("notification")).Send(sel("request")).Send(sel("content")).Send(sel("userInfo"))
	if info == 0 {
		return
	}
	value := info.Send(sel("objectForKey:"), nsString(activationKey))
	if value != 0 && objc.Send[bool](value, sel("isKindOfClass:"), class("NSString")) {
		t.events.notificationClicked(goString(value))
	}
}

// startNotifications finds the notification center, where the program is
// in a bundle, and asks it for the stored answer.
func (t *macTray) startNotifications() {
	if _, err := purego.Dlopen(userNotificationsPath, purego.RTLD_LAZY|purego.RTLD_GLOBAL); err != nil {
		return
	}
	bundle := class("NSBundle").Send(sel("mainBundle"))
	if bundle == 0 || bundle.Send(sel("bundleIdentifier")) == 0 || class("UNUserNotificationCenter") == 0 || macNotificationClass == 0 {
		return
	}
	t.center = class("UNUserNotificationCenter").Send(sel("currentNotificationCenter")).Send(sel("retain"))
	if t.center == 0 {
		return
	}
	t.delegate = objc.ID(macNotificationClass).Send(sel("new"))
	t.center.Send(sel("setDelegate:"), t.delegate)
	t.auth = Unknown
	t.refresh()
}

// refresh reads the stored answer, which arrives on the center's thread.
func (t *macTray) refresh() {
	block := objc.NewBlock(func(_ objc.Block, settings objc.ID) {
		t.events.authorized(authorizationFor(objc.Send[int64](settings, sel("authorizationStatus"))))
	})
	t.center.Send(sel("getNotificationSettingsWithCompletionHandler:"), block)
	block.Release()
}

func authorizationFor(status int64) Authorization {
	switch status {
	case unAuthorizationDenied:
		return Denied
	case unAuthorizationAuthorized, unAuthorizationProvisional, unAuthorizationEphemeral:
		return Authorized
	}
	return Unknown
}

func (t *macTray) available() bool { return true }

func (t *macTray) authorization() Authorization { return t.auth }

// requestAuthorization puts up macOS's question. A false answer can mean
// either denied or not yet decided, so the stored answer is read again
// afterwards, as the app does.
func (t *macTray) requestAuthorization() {
	if t.center == 0 {
		return
	}
	withPool(func() {
		block := objc.NewBlock(func(_ objc.Block, _ bool, failure objc.ID) {
			if failure != 0 {
				t.events.authorized(Unknown)
				return
			}
			t.refresh()
		})
		t.center.Send(sel("requestAuthorizationWithOptions:completionHandler:"), uint64(unAuthorizationOptionAlert), block)
		block.Release()
	})
}

func (t *macTray) show(s trayState) {
	withPool(func() {
		if t.item == 0 {
			bar := class("NSStatusBar").Send(sel("systemStatusBar"))
			t.item = bar.Send(sel("statusItemWithLength:"), float64(nsVariableStatusItemSize)).Send(sel("retain"))
		}
		button := t.item.Send(sel("button"))
		if img := statusImage(s); img != 0 {
			button.Send(sel("setImage:"), img)
			button.Send(sel("setTitle:"), nsString(""))
		} else {
			button.Send(sel("setTitle:"), nsString(s.tooltip))
		}
		button.Send(sel("setToolTip:"), nsString(s.tooltip))
		old := t.menu
		t.menu = statusMenu(s)
		t.item.Send(sel("setMenu:"), t.menu)
		if old != 0 {
			old.Send(sel("release"))
		}
	})
}

// statusImage is the icon as an autoreleased NSImage, from the picture as
// PNG, or 0 without one.
func statusImage(s trayState) objc.ID {
	img := iconAt(s.icons, 2*statusIconPoints)
	if img == nil {
		return 0
	}
	var buf bytes.Buffer
	if png.Encode(&buf, img) != nil {
		return 0
	}
	b := buf.Bytes()
	data := class("NSData").Send(sel("dataWithBytes:length:"), unsafe.Pointer(&b[0]), uint64(len(b)))
	image := class("NSImage").Send(sel("alloc")).Send(sel("initWithData:"), data)
	if image == 0 {
		return 0
	}
	image.Send(sel("setSize:"), macSize{statusIconPoints, statusIconPoints})
	return image.Send(sel("autorelease"))
}

// statusMenu builds the menu, owned by the caller; each item's tag is its
// line's index.
func statusMenu(s trayState) objc.ID {
	menu := class("NSMenu").Send(sel("alloc")).Send(sel("initWithTitle:"), nsString(""))
	menu.Send(sel("setAutoenablesItems:"), false)
	for i, e := range s.menu {
		if e.separator {
			menu.Send(sel("addItem:"), class("NSMenuItem").Send(sel("separatorItem")))
			continue
		}
		item := class("NSMenuItem").Send(sel("alloc")).Send(sel("initWithTitle:action:keyEquivalent:"),
			nsString(e.text), sel("kvitTrayChoose:"), nsString(""))
		item.Send(sel("setTag:"), int64(i))
		item.Send(sel("setTarget:"), macTarget)
		item.Send(sel("setEnabled:"), !e.disabled)
		menu.Send(sel("addItem:"), item)
		item.Send(sel("release"))
	}
	return menu
}

func (t *macTray) hide() {
	if t.item == 0 {
		return
	}
	withPool(func() {
		class("NSStatusBar").Send(sel("systemStatusBar")).Send(sel("removeStatusItem:"), t.item)
		t.item.Send(sel("release"))
		t.item = 0
		if t.menu != 0 {
			t.menu.Send(sel("release"))
			t.menu = 0
		}
	})
}

// notify posts the notification with its ID in userInfo, so a click on it
// can say which it was.
func (t *macTray) notify(n Notification) {
	if t.center == 0 {
		return
	}
	withPool(func() {
		content := class("UNMutableNotificationContent").Send(sel("alloc")).Send(sel("init"))
		content.Send(sel("setTitle:"), nsString(n.Title))
		content.Send(sel("setBody:"), nsString(n.Message))
		content.Send(sel("setUserInfo:"), class("NSDictionary").Send(sel("dictionaryWithObject:forKey:"),
			nsString(n.ID), nsString(activationKey)))
		id := class("NSUUID").Send(sel("UUID")).Send(sel("UUIDString"))
		request := class("UNNotificationRequest").Send(sel("requestWithIdentifier:content:trigger:"), id, content, objc.ID(0))
		t.center.Send(sel("addNotificationRequest:withCompletionHandler:"), request, objc.ID(0))
		content.Send(sel("release"))
	})
}

func (t *macTray) close() {
	t.hide()
	if t.center != 0 {
		if t.center.Send(sel("delegate")) == t.delegate {
			t.center.Send(sel("setDelegate:"), objc.ID(0))
		}
		t.center.Send(sel("release"))
		t.center = 0
		t.delegate.Send(sel("release"))
		t.delegate = 0
	}
	macActive.CompareAndSwap(t, nil)
}
