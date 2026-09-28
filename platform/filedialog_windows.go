package platform

// The file dialog on Windows is the common item dialog, IFileOpenDialog,
// called through its COM interface without cgo, as unison's own picker
// calls it. COM is already set up on the UI thread, where unison calls
// OleInitialize at start-up; the dialog is an apartment object and is used
// on that thread only.

import (
	"runtime"
	"syscall"
	"unsafe"

	"github.com/richardwilkes/toolbox/v2/errs"
	"github.com/richardwilkes/unison"
	"golang.org/x/sys/windows"
)

var (
	ole32 = windows.NewLazySystemDLL("ole32.dll")

	procCoCreateInstance            = ole32.NewProc("CoCreateInstance")
	procSHCreateItemFromParsingName = shell32.NewProc("SHCreateItemFromParsingName")
	procGetActiveWindow             = user32.NewProc("GetActiveWindow")

	clsidFileOpenDialog = mustGUID("{DC1C5A9C-E88A-4DDE-A5A1-60F82A20AEF7}")
	iidFileOpenDialog   = mustGUID("{D57C7288-D4AD-4768-BE02-9D969532D960}")
	iidShellItem        = mustGUID("{43826D1E-E718-42EE-BC55-A1E261C37BFE}")
)

const (
	clsctxInprocServer = 0x1
	sigdnFileSysPath   = 0x80058000

	fosForceFileSystem  = 0x40
	fosAllowMultiSelect = 0x200
	fosPathMustExist    = 0x800
	fosFileMustExist    = 0x1000

	// The methods' places in the interfaces' tables: IUnknown's three,
	// then IModalWindow's Show, then IFileDialog's and IFileOpenDialog's.
	comRelease           = 2
	dialogShow           = 3
	dialogSetFileTypes   = 4
	dialogSetTypeIndex   = 5
	dialogSetOptions     = 9
	dialogGetOptions     = 10
	dialogSetFolder      = 12
	dialogSetTitle       = 17
	dialogGetResults     = 27
	itemArrayGetCount    = 7
	itemArrayGetItemAt   = 8
	itemGetDisplayName   = 5
	hresultCancelled     = 0x800704C7 // HRESULT_FROM_WIN32(ERROR_CANCELLED)
	hresultFailureBitSet = 0x80000000
)

func mustGUID(s string) windows.GUID {
	g, err := windows.GUIDFromString(s)
	if err != nil {
		panic(err)
	}
	return g
}

// comObject is a COM interface pointer: the object begins with a pointer
// to its table of methods.
type comObject struct{ methods *[32]uintptr }

func (o *comObject) call(method int, args ...uintptr) uintptr {
	r, _, _ := syscall.SyscallN(o.methods[method], append([]uintptr{uintptr(unsafe.Pointer(o))}, args...)...)
	return r
}

func (o *comObject) release() { o.call(comRelease) }

// filterSpec is COMDLG_FILTERSPEC.
type filterSpec struct{ name, spec *uint16 }

func (d FileDialog) openNative() []string {
	var dialog *comObject
	if hr, _, _ := procCoCreateInstance.Call(uintptr(unsafe.Pointer(&clsidFileOpenDialog)), 0, clsctxInprocServer,
		uintptr(unsafe.Pointer(&iidFileOpenDialog)), uintptr(unsafe.Pointer(&dialog))); uint32(hr) != 0 || dialog == nil {
		errs.Log(errs.New("the file dialog could not be created"), "hresult", hr)
		return d.openWithUnison()
	}
	defer dialog.release()

	if d.Title != "" {
		dialog.call(dialogSetTitle, uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(d.Title))))
	}
	types := d.fileTypes()
	specs := make([]filterSpec, len(types))
	for i, t := range types {
		specs[i] = filterSpec{windows.StringToUTF16Ptr(t[0]), windows.StringToUTF16Ptr(t[1])}
	}
	dialog.call(dialogSetFileTypes, uintptr(len(specs)), uintptr(unsafe.Pointer(&specs[0])))
	dialog.call(dialogSetTypeIndex, 1)
	runtime.KeepAlive(specs)

	var options uint32
	dialog.call(dialogGetOptions, uintptr(unsafe.Pointer(&options)))
	options |= fosForceFileSystem | fosPathMustExist | fosFileMustExist
	if d.Multiple {
		options |= fosAllowMultiSelect
	}
	dialog.call(dialogSetOptions, uintptr(options))

	if d.Folder != "" {
		var folder *comObject
		if hr, _, _ := procSHCreateItemFromParsingName.Call(uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(d.Folder))), 0,
			uintptr(unsafe.Pointer(&iidShellItem)), uintptr(unsafe.Pointer(&folder))); uint32(hr) == 0 && folder != nil {
			dialog.call(dialogSetFolder, uintptr(unsafe.Pointer(folder)))
			folder.release()
		}
	}

	// The dialog is owned by the active window, which Windows disables
	// while the dialog runs its own message loop; unison's windows go on
	// drawing through that loop.
	active := unison.ActiveWindow()
	owner, _, _ := procGetActiveWindow.Call()
	// An HRESULT is 32 bits; the register it comes back in is wider.
	hr := uint32(dialog.call(dialogShow, owner))
	if active != nil && active.IsVisible() {
		active.ToFront()
	}
	if hr != 0 {
		if hr != hresultCancelled && hr&hresultFailureBitSet != 0 {
			errs.Log(errs.New("the file dialog failed"), "hresult", hr)
		}
		return nil
	}
	return shownResults(dialog, d.Multiple)
}

// shownResults are the files the reader chose, as file system paths.
func shownResults(dialog *comObject, multiple bool) []string {
	var items *comObject
	if uint32(dialog.call(dialogGetResults, uintptr(unsafe.Pointer(&items)))) != 0 || items == nil {
		return nil
	}
	defer items.release()
	var count uint32
	items.call(itemArrayGetCount, uintptr(unsafe.Pointer(&count)))
	var paths []string
	for i := range count {
		var item *comObject
		if uint32(items.call(itemArrayGetItemAt, uintptr(i), uintptr(unsafe.Pointer(&item)))) != 0 || item == nil {
			continue
		}
		var name *uint16
		if uint32(item.call(itemGetDisplayName, sigdnFileSysPath, uintptr(unsafe.Pointer(&name)))) == 0 && name != nil {
			paths = append(paths, windows.UTF16PtrToString(name))
			windows.CoTaskMemFree(unsafe.Pointer(name))
		}
		item.release()
	}
	if !multiple && len(paths) > 1 {
		paths = paths[:1]
	}
	return paths
}
