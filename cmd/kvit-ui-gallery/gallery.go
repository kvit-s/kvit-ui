package main

import (
	"github.com/richardwilkes/toolbox/v2/geom"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/paintstyle"
)

// newGalleryWindow opens the gallery's window. firstFrame, when not nil, is
// called once, after the window has drawn for the first time.
func newGalleryWindow(firstFrame func()) (*unison.Window, error) {
	wnd, err := unison.NewWindow("kvit-ui gallery")
	if err != nil {
		return nil, err
	}
	content := wnd.Content()
	content.SetLayout(&unison.FlexLayout{Columns: 1, HAlign: align.Middle, VAlign: align.Middle})
	drawn := false
	content.DrawCallback = func(gc *unison.Canvas, rect geom.Rect) {
		gc.DrawRect(rect, unison.ThemeSurface.Paint(gc, rect, paintstyle.Fill))
		if !drawn {
			drawn = true
			if firstFrame != nil {
				unison.InvokeTask(firstFrame)
			}
		}
	}
	label := unison.NewLabel()
	label.SetTitle("kvit-ui gallery: each component will have a page here")
	label.SetLayoutData(&unison.FlexLayoutData{HAlign: align.Middle, VAlign: align.Middle, HGrab: true, VGrab: true})
	content.AddChild(label)
	wnd.SetContentRect(geom.NewRect(100, 100, 900, 600))
	wnd.ToFront()
	return wnd, nil
}
