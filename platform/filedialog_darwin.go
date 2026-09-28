package platform

// The file dialog on macOS is NSOpenPanel, called through purego's runtime
// bindings as the tray is (tray_darwin.go, whose helpers it uses). A panel
// shows no title bar, so the title is also its message, the line of text
// above the files. The filters are a pop-up button in the panel's accessory
// view, below the file list, as the Qt dialog put them; choosing one sets
// the file types the panel lets through.
//
// This file is compiled for darwin/arm64 and darwin/amd64 on Linux and has
// not yet run on a Mac, as none was available when it was written.

import (
	"strings"
	"sync"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
	"github.com/richardwilkes/unison"
)

const (
	nsModalResponseOK       = 1
	filterTargetClassName   = "KvitFileDialogFilterTarget"
	filterTargetChooseEvent = "kvitFileDialogFilterChosen:"
)

var (
	filterTargetOnce sync.Once
	filterTarget     objc.ID
	// shownPanel is the panel that is open and its filters; one runs at a
	// time, since each is modal.
	shownPanel struct {
		panel   objc.ID
		filters []nameFilter
	}
)

func (d FileDialog) openNative() []string {
	if _, err := purego.Dlopen(appKitPath, purego.RTLD_LAZY|purego.RTLD_GLOBAL); err != nil {
		return d.openWithUnison()
	}
	registerFilterTarget()
	var paths []string
	active := unison.ActiveWindow()
	withPool(func() {
		panel := class("NSOpenPanel").Send(sel("openPanel"))
		if d.Title != "" {
			panel.Send(sel("setTitle:"), nsString(d.Title))
			panel.Send(sel("setMessage:"), nsString(d.Title))
		}
		panel.Send(sel("setCanChooseFiles:"), true)
		panel.Send(sel("setCanChooseDirectories:"), false)
		panel.Send(sel("setAllowsMultipleSelection:"), d.Multiple)
		if d.Folder != "" {
			panel.Send(sel("setDirectoryURL:"), class("NSURL").Send(sel("fileURLWithPath:"), nsString(d.Folder)))
		}
		filters := d.filters()
		if len(filters) > 1 && filterTarget != 0 {
			popup := class("NSPopUpButton").Send(sel("alloc")).Send(sel("init")).Send(sel("autorelease"))
			for _, f := range filters {
				popup.Send(sel("addItemWithTitle:"), nsString(f.label))
			}
			popup.Send(sel("setTarget:"), filterTarget)
			popup.Send(sel("setAction:"), sel(filterTargetChooseEvent))
			popup.Send(sel("sizeToFit"))
			panel.Send(sel("setAccessoryView:"), popup)
			panel.Send(sel("setAccessoryViewDisclosed:"), true)
		}
		shownPanel.panel, shownPanel.filters = panel, filters
		defer func() { shownPanel.panel, shownPanel.filters = 0, nil }()
		allowFilter(0)
		if objc.Send[int64](panel, sel("runModal")) != nsModalResponseOK {
			return
		}
		urls := panel.Send(sel("URLs"))
		for i := range objc.Send[uint64](urls, sel("count")) {
			if path := goString(urls.Send(sel("objectAtIndex:"), i).Send(sel("path"))); path != "" {
				paths = append(paths, path)
			}
		}
	})
	if active != nil && active.IsVisible() {
		active.ToFront()
	}
	return paths
}

// registerFilterTarget registers the class the filters' pop-up reports to,
// whose action lets through the files of the filter chosen.
func registerFilterTarget() {
	filterTargetOnce.Do(func() {
		target, err := objc.RegisterClass(filterTargetClassName, objc.GetClass("NSObject"), nil, nil, []objc.MethodDef{{
			Cmd: sel(filterTargetChooseEvent),
			Fn: func(_ objc.ID, _ objc.SEL, sender objc.ID) {
				allowFilter(int(objc.Send[int64](sender, sel("indexOfSelectedItem"))))
			},
		}})
		if err == nil {
			filterTarget = objc.ID(target).Send(sel("new"))
		}
	})
}

// allowFilter sets the open panel's file types to one filter's extensions.
// A filter that matches every file, or has a pattern that is not an
// extension, lets every file through.
func allowFilter(index int) {
	panel, filters := shownPanel.panel, shownPanel.filters
	if panel == 0 || index < 0 || index >= len(filters) {
		return
	}
	f := filters[index]
	types := class("NSMutableArray").Send(sel("array"))
	for _, p := range f.patterns {
		ext, ok := strings.CutPrefix(p, "*.")
		if !ok || ext == "*" || strings.ContainsAny(ext, "*?[") {
			panel.Send(sel("setAllowedFileTypes:"), objc.ID(0))
			return
		}
		types.Send(sel("addObject:"), nsString(ext))
	}
	panel.Send(sel("setAllowedFileTypes:"), types)
}
