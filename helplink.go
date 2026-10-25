package kvitui

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/role"
)

// HelpLink opens the help topic for the place it stands, where the window
// cannot say what the topic says (HLP-6, decided Q-34). On the Windows form,
// which Linux takes, it is a link worded as the question the topic answers,
// "How is the balance worked out?", at the bottom of the content or under
// the group it concerns; on the macOS form it is the round question-mark
// button, one per window. A topic is an HTML page shipped with the app,
// opened in the system's browser at its anchor: "budgets#rollover" is
// budgets.html in UI.HelpRoot, at #rollover.
type HelpLink struct {
	unison.Panel
	ui *UI
	// Topic is the page and its anchor under UI.HelpRoot.
	Topic string
	// Question is what the topic answers, a whole question ending in a
	// question mark; the link's words on the Windows form.
	Question string

	link   *Link
	button *IconButton
}

// NewHelpLink returns the help for a topic.
func NewHelpLink(ui *UI, topic, question string) *HelpLink {
	h := &HelpLink{ui: ui, Topic: topic, Question: question}
	h.Self = h
	h.link = NewLink(ui, "")
	h.link.OnActivate = h.Open
	h.button = NewIconButton(ui, "question", "Help")
	h.button.Form = Ordinary
	h.button.OnClick = h.Open
	h.SetLayout(syncing{Layout: &unison.FlexLayout{Columns: 1}, sync: func() {
		h.link.Text = h.Question
		h.button.Explanation = h.Question
		if ui.form == FormMac {
			showOnly(h.AsPanel(), h.button)
		} else {
			showOnly(h.AsPanel(), h.link)
		}
	}})
	return h
}

// TopicURL is where the topic is: a file URL to its page in UI.HelpRoot,
// with its anchor.
func (u *UI) TopicURL(topic string) string {
	page, anchor, _ := strings.Cut(topic, "#")
	path := filepath.Join(u.HelpRoot, page+".html")
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}
	link := url.URL{Scheme: "file", Path: filepath.ToSlash(path), Fragment: anchor}
	return link.String()
}

// Open opens the topic in the system's browser.
func (h *HelpLink) Open() {
	_ = unison.OpenBrowser(h.ui.TopicURL(h.Topic))
}

// ProvideAccessibility leaves the holder out; the link or button inside it
// is what a reader meets.
func (h *HelpLink) ProvideAccessibility(b *unison.AccessibilityBuilder) {
	b.Node().Role = role.None
}
