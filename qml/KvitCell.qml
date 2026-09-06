// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Delegates in this file read ids from the enclosing component. Bound is what
// makes that legal rather than accidental: without it a delegate resolves an
// outer id at run time through the object hierarchy, which works until the
// delegate is reused for a different row and quietly reads the wrong one.
pragma ComponentBehavior: Bound

import QtQuick
import Kvit.Ui

// One cell of a table, drawn according to what kind of value it holds.
//
// The kind comes from the column rather than from the value, which is what
// makes a column readable: every cell in it is aligned the same way, set in
// the same face, and says the same thing about a missing value. A cell that
// decided for itself would give a column where some rows are right-aligned and
// some are not, depending on whether that row's value happened to parse as a
// number.
//
// The kinds are TableModelBase::CellKind, so the model says what a column is
// once and every view of it agrees.
//
// A cell formats nothing. It is given the string a column holds and it draws
// that string, which is the same contract KvitBeforeAfter documents for its
// `before` and `after`. Money is where that matters most: an amount is a
// signed 64-bit count of the minor unit of its own currency, and both halves
// of that sentence are lost the moment a cell converts one to a JavaScript
// number and prints two decimal places after it.
Item {
    id: root

    // "Text" | "Figure" | "Chip" | "Slug" | "Date" | "Money" | "Marks" | "Check"
    property string kind: "Text"
    property var value: undefined
    property bool measured: true
    // For a Chip cell: which tone. The model supplies it through MarkRole.
    property string mark: "neutral"
    property string unit: ""
    // The whole value where `value` is the shortened one a narrow column has
    // room for, which a model supplies through TableModelBase's FullTextRole
    // under the role name `fullText`. Empty says the cell already shows
    // everything there is; where it does not, the pointer and the keyboard
    // both disclose the rest.
    property string fullValue: ""
    property bool selected: false
    // Whether the keyboard cursor is on this cell. A cell is not a control
    // and takes no focus of its own, so the view that owns the cursor says
    // where it is, and the full-value disclosure below opens for a keyboard
    // reader as well as for a pointer — which is the pair `hovered ||
    // visualFocus` gives everywhere a real control does this.
    property bool current: false
    // For a Marks cell: every state this row is in, as { tone, shape, label }
    // objects, one dot each. An entry may also set `hollow`, which is what
    // lets more states than there are shapes still differ by something other
    // than their colour.
    property var marks: []
    // For a Check cell: whether this row is picked.
    property bool checked: false

    // What a Check cell asks for. The cell does not change `checked` itself:
    // whatever owns the selection does, and the cell redraws from it, so a
    // request that is refused leaves no tick behind.
    signal toggled(bool checked)

    implicitHeight: Interface.rowHeightSlim
    implicitWidth: content.implicitWidth

    readonly property bool numeric: kind === "Figure" || kind === "Money"

    // The value as the string it is drawn as, with nothing in place of a
    // value that is not there.
    readonly property string valueText:
        root.value === undefined || root.value === null ? "" : String(root.value)

    // Whether a money value is money leaving, read from the sign already
    // written into it. The cell has a formatted string rather than a number,
    // and the character that marks a negative belongs to the reader's locale:
    // some write the hyphen and some the typographic minus.
    readonly property bool negative: {
        if (root.valueText === "")
            return false
        const sign = root.valueText.charAt(0)
        return sign === "-" || sign === "−"
            || sign === Qt.locale().negativeSign
    }

    Accessible.role: Accessible.Cell
    Accessible.name: {
        switch (root.kind) {
        // The marks name themselves, one accessible object each, so a reader
        // is told where one state ends and the next begins. Joining them into
        // a single string here is the collision the Marks kind exists to
        // prevent, and it would arrive in a copy of the row as well.
        case "Marks": {
            // A model that has no marks for this row answers with nothing at
            // all rather than with an empty list, which is a different thing
            // from a row that is in no state and has to read the same.
            const count = root.marks ? root.marks.length : 0
            if (count === 0)
                return ""
            // Two source strings, chosen here by the count, rather than one
            // with `%n` in it. With no translator installed Qt substitutes
            // the number into the source and does not choose a plural form,
            // so a row in exactly one state announces itself as "1 marks".
            // The count is still handed to qsTr, so which form to use stays
            // Qt's choice and a translator into a language with three plural
            // forms gets all three.
            return count === 1 ? qsTr("%n mark", "", count)
                               : qsTr("%n marks", "", count)
        }
        // A checkbox is a control: it publishes its own name and its own
        // checked state, and a name on the cell around it is said twice.
        case "Check":
            return ""
        default:
            if (!root.measured)
                return qsTr("not measured")
            // The whole value rather than the shortened one. A screen reader
            // that says a truncated payee says a different payee.
            return root.fullValue !== "" ? root.fullValue : root.valueText
        }
    }

    Loader {
        id: content
        anchors.fill: parent
        anchors.leftMargin: Interface.spaceNear
        anchors.rightMargin: Interface.spaceNear
        sourceComponent: {
            switch (root.kind) {
            case "Figure": return figureCell
            case "Money":  return moneyCell
            case "Chip":   return chipCell
            case "Slug":   return slugCell
            case "Date":   return dateCell
            case "Marks":  return marksCell
            case "Check":  return checkCell
            default:       return textCell
            }
        }
    }

    Component {
        id: textCell
        // The text, and a route to the rest of it.
        //
        // A column narrow enough to elide is the ordinary case in a table the
        // reader resizes, and a name that is only ever seen cut short is a
        // name they cannot check. Two things are disclosed: the full value
        // where the model answered one and it differs from what is drawn, and
        // the drawn value itself where the label ran out of room. Both open
        // under the pointer and under the keyboard cursor, because a reader
        // who never touches the pointer has the same column to read.
        Item {
            id: textHolder

            // The width the text wants. Without it a Text cell reports zero,
            // because the label fills this item and so this item has no
            // implicit size of its own — KvitTable never asks, since it sizes
            // columns from the model, but a caller laying cells out in a plain
            // Row does, and gets cells of no width at all.
            implicitWidth: shown.implicitWidth

            readonly property string disclosed: {
                if (root.fullValue !== "" && root.fullValue !== root.valueText)
                    return root.fullValue
                return shown.truncated ? root.valueText : ""
            }

            KvitLabel {
                id: shown
                anchors.fill: parent
                text: root.valueText
                role: "body"
                color: root.selected ? Theme.textPrimary : Theme.textSecondary
                verticalAlignment: Text.AlignVCenter
            }

            HoverHandler { id: textHover }

            // Built when the pointer or the keyboard arrives rather than
            // with the row, for the reason the marks give at length below: a
            // table draws several hundred cells at a time and a tooltip is a
            // popup with a window behind it.
            Loader {
                anchors.fill: parent
                active: (textHover.hovered || root.current)
                        && textHolder.disclosed !== ""
                sourceComponent: KvitTooltip {
                    text: textHolder.disclosed
                    visible: true
                }
            }
        }
    }

    Component {
        id: figureCell
        // Right-aligned, because a column of numbers is compared down its
        // units digit and a left-aligned column of numbers cannot be.
        Item {
            KvitFigure {
                anchors.right: parent.right
                anchors.verticalCenter: parent.verticalCenter
                value: root.valueText
                unit: root.unit
                measured: root.measured
                role: "body"
            }
        }
    }

    Component {
        id: moneyCell
        // An amount that has already been formatted, exactly as
        // KvitBeforeAfter takes its `before` and `after`. Two facts about
        // money make that necessary rather than tidy.
        //
        // An amount is a signed 64-bit count of minor units and a JavaScript
        // number carries 53 bits of integer, so a cell that converted the
        // value to a number would round large amounts and never say that it
        // had. And how many minor digits an amount has belongs to its
        // currency — none for the yen, three for the dinar — so a cell that
        // printed two of them would give a ¥1,000 fare as ¥10.00.
        //
        // What is left here is the drawing. Right-aligned with the tabular
        // numerals KvitFigure sets, so a column of amounts in one currency
        // lines up on its decimal separator; the em dash where the amount was
        // never measured, which KvitFigure draws and nothing else in the
        // estate draws for itself; and the sign written rather than coloured,
        // with the colour agreeing with a minus that is already there. That
        // last part is what keeps the column readable in a printed export and
        // for a reader who cannot separate the hues.
        Item {
            KvitFigure {
                anchors.right: parent.right
                anchors.verticalCenter: parent.verticalCenter
                value: root.valueText
                unit: root.unit
                measured: root.measured
                role: "body"
                color: root.negative ? Theme.danger : Theme.textPrimary
            }
        }
    }

    Component {
        id: chipCell
        Item {
            KvitChip {
                anchors.left: parent.left
                anchors.verticalCenter: parent.verticalCenter
                text: root.valueText
                tone: root.mark
            }
        }
    }

    Component {
        id: slugCell
        KvitSlug {
            anchors.verticalCenter: parent.verticalCenter
            text: root.valueText
            // No ground inside a table: a box on every row of a column reads
            // as a column of boxes rather than as a column of identifiers.
            ground: false
        }
    }

    Component {
        id: dateCell
        KvitLabel {
            text: {
                const v = root.value
                if (v === undefined || v === null)
                    return ""
                return v instanceof Date
                    ? Qt.formatDate(v, Locale.ShortFormat) : String(v)
            }
            role: "body"
            color: Theme.textSecondary
            tabular: true
            verticalAlignment: Text.AlignVCenter
        }
    }

    Component {
        id: marksCell
        // Every state a row is in at once, as one dot each.
        //
        // Separate objects rather than a list of words, which is the whole
        // reason this kind exists. Three states written into one string run
        // together for the eye, come out of a copy as a sentence, and reach a
        // screen reader as one run of text with no boundary in it — and where
        // a state shares its name with a payee or a category, the reader
        // cannot tell the row's own value from the mark on it.
        //
        // Each entry says three things. `tone` is the meaning, in the same
        // words a chip and a toast use. `shape` is the second channel beside
        // the hue, because about one man in twelve cannot separate red from
        // green and every screenshot is read in grayscale by somebody
        // eventually; a caller that names only a tone gets the shape
        // KvitTimeline gives that tone. And `label` is the one string that is
        // both the word on hover and the name a screen reader says, so the
        // two cannot drift apart.
        Item {
            implicitWidth: markRow.implicitWidth

            Row {
                id: markRow
                anchors.left: parent.left
                anchors.verticalCenter: parent.verticalCenter
                spacing: Interface.spaceNear

                Repeater {
                    model: root.marks

                    delegate: Item {
                        id: markItem

                        required property var modelData

                        // Guarded once, here, rather than in each of the
                        // three bindings below: a binding that throws inside
                        // a delegate leaves the cell blank and says nothing
                        // about why.
                        readonly property var entry:
                            markItem.modelData === undefined
                            || markItem.modelData === null
                                ? ({}) : markItem.modelData

                        readonly property string markLabel:
                            markItem.entry.label === undefined
                                ? "" : String(markItem.entry.label)
                        readonly property string markTone:
                            markItem.entry.tone === undefined
                                ? "neutral" : String(markItem.entry.tone)
                        readonly property string markShape: {
                            if (markItem.entry.shape !== undefined)
                                return String(markItem.entry.shape)
                            return markItem.markTone === "danger" ? "diamond"
                                 : markItem.markTone === "warning" ? "square"
                                 : "circle"
                        }
                        // Optional, and the reason it is here: KvitDot has
                        // three shapes, a row can be in more than three
                        // states at once, and two states that differ only in
                        // colour differ in nothing at all for the reader the
                        // shape was put there for. Outlined against filled
                        // doubles what the shapes can say.
                        readonly property bool markHollow:
                            markItem.entry.hollow === true

                        readonly property color markColor: {
                            switch (markItem.markTone) {
                            case "accent":  return Theme.accent
                            case "success": return Theme.success
                            case "warning": return Theme.warning
                            case "danger":  return Theme.danger
                            case "info":    return Theme.link
                            default:        return Theme.textMuted
                            }
                        }

                        // Wider than the dot it holds, so the word is
                        // reachable by a pointer that is near a six-pixel
                        // mark rather than exactly on it, and so two marks
                        // keep a gap the eye reads as a gap.
                        implicitWidth: Interface.spaceLoose
                        implicitHeight: Interface.spaceLoose

                        Accessible.role: Accessible.Indicator
                        Accessible.name: markItem.markLabel
                        Accessible.ignored: markItem.markLabel === ""

                        KvitDot {
                            anchors.centerIn: parent
                            color: markItem.markColor
                            shape: markItem.markShape
                            hollow: markItem.markHollow
                            // Named by the item around it, which is what
                            // carries the tooltip as well; a label here would
                            // be said twice.
                            label: ""
                        }

                        HoverHandler { id: markHover }

                        // Built when the pointer arrives rather than with the
                        // row. A tooltip is a popup with a window behind it,
                        // and a table draws several hundred cells at a time:
                        // one popup per mark per row is the single most
                        // expensive thing a cell could hold, and the reader
                        // sees at most one of them.
                        Loader {
                            // Filling the mark, so the tooltip sits where it
                            // sat when it was declared here directly: a popup
                            // positions itself against its parent item.
                            anchors.fill: parent
                            active: markHover.hovered && markItem.markLabel !== ""
                            sourceComponent: KvitTooltip {
                                text: markItem.markLabel
                                visible: true
                            }
                        }
                    }
                }
            }
        }
    }

    Component {
        id: checkCell
        // A real checkbox rather than a drawn tick, so the cell reaches the
        // keyboard and a screen reader is told both that it is a checkbox and
        // whether it is ticked.
        Item {
            implicitWidth: box.implicitWidth

            KvitCheck {
                id: box
                anchors.left: parent.left
                anchors.verticalCenter: parent.verticalCenter
                checked: root.checked
                // The box has no text of its own, and a control whose whole
                // label is a shape has nothing else to say to a screen
                // reader.
                Accessible.name: qsTr("Select this row")
                onToggled: {
                    // CheckBox has already moved its own state, which breaks
                    // the binding to `checked` and leaves the tick saying
                    // something the model has not agreed to. Put the binding
                    // back and let the signal carry the request instead, so
                    // the row is ticked when the selection says it is and at
                    // no other time.
                    const wanted = box.checked
                    box.checked = Qt.binding(() => root.checked)
                    root.toggled(wanted)
                }
            }
        }
    }
}
