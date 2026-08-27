// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.
#ifndef TYPOGRAPHY_H
#define TYPOGRAPHY_H

#include <QObject>
#include <QString>
#include <QStringList>

class SettingsStore;

// Typography settings (features.md §10.2):
// editor font family, base size, line-height multiplier, paragraph
// spacing, maximum content width, and the monospace family for code.
// Heading and code sizes derive from the base by fixed ratios — H1
// 32/15, H2 24/15, H3 20/15, H4 17/15, code 13/15 — so one setting
// scales the whole document coherently. The ratios are written against
// 15 because that is the base they were measured at; the base itself is
// a default like any other and setting it to 15 reproduces the original
// pixel sizes exactly.
//
// All properties notify through one typographyChanged signal, and all
// setters clamp, so a hand-edited settings file cannot produce an
// unusable layout. State lives in the settings store.
class Typography : public QObject
{
    Q_OBJECT

    // "" = the application default font.
    Q_PROPERTY(QString fontFamily READ fontFamily WRITE setFontFamily
                   NOTIFY typographyChanged)
    Q_PROPERTY(int baseSize READ baseSize WRITE setBaseSize
                   NOTIFY typographyChanged)
    // Proportional line height: 1.0 = the font's natural height.
    Q_PROPERTY(qreal lineHeight READ lineHeight WRITE setLineHeight
                   NOTIFY typographyChanged)
    // Pixels between blocks (the block list's spacing).
    Q_PROPERTY(int paragraphSpacing READ paragraphSpacing
                   WRITE setParagraphSpacing NOTIFY typographyChanged)
    // 0 = fill the editor pane; otherwise the block column is capped
    // at this many pixels and centered.
    Q_PROPERTY(int maxContentWidth READ maxContentWidth
                   WRITE setMaxContentWidth NOTIFY typographyChanged)
    Q_PROPERTY(QString monoFamily READ monoFamily WRITE setMonoFamily
                   NOTIFY typographyChanged)
    // The two sizes chrome asks for by name rather than per block: ordinary
    // body text, and the monospace size code is set at. A diagram's labels, a
    // query block's table and a math block's fallback text all want one of
    // these, and they used to ask by passing Block::Paragraph or
    // Block::CodeBlock to the size lookup as a stand-in — which read as
    // "size this as a paragraph would be" and made those call sites look like
    // block-type decisions when they are not.
    Q_PROPERTY(int bodySize READ bodySize NOTIFY typographyChanged)
    Q_PROPERTY(int monoSize READ monoSize NOTIFY typographyChanged)

public:
    // Which entry of the frozen type scale a run of text renders at.
    // sizeForRole() turns a role into pixels; the ratios stay where they are.
    //
    // In the editor this enumeration lives on the block header, because there
    // a block's kind is what answers which role it is. Nothing outside the
    // editor has blocks, so it belongs to the type scale itself here, and the
    // editor's copy is the one that goes away in Wave 2. The order and the
    // values match that copy exactly, because both are handed across as an
    // int through a model role and a mismatch would silently resize text.
    enum class FontRole {
        Body,
        Heading1,
        Heading2,
        Heading3,
        Heading4,
        // A code fence, and anything else drawn in the configured monospace
        // family at the code size.
        Mono,
    };
    Q_ENUM(FontRole)

    // Clamps: wide enough for every reasonable preference, tight
    // enough that a corrupt settings value cannot break the layout.
    static constexpr int MinBaseSize = 10, MaxBaseSize = 28;
    static constexpr qreal MinLineHeight = 1.0, MaxLineHeight = 2.0;
    // What a reader who has never opened the settings dialog gets. Named
    // rather than written twice, because the member initialisers below and
    // resetToDefaults() must agree: they did not have to before, and a
    // "Reset typography" that restored something other than the starting
    // state is the failure that costs the most to notice.
    static constexpr int DefaultBaseSize = 14;
    static constexpr qreal DefaultLineHeight = 1.3;
    static constexpr int DefaultParagraphSpacing = 4;
    static constexpr int MinParagraphSpacing = 0, MaxParagraphSpacing = 40;
    static constexpr int MinContentWidth = 300;  // when non-zero

    explicit Typography(QObject *parent = nullptr);

    void setSettings(SettingsStore *settings);

    QString fontFamily() const { return m_fontFamily; }
    void setFontFamily(const QString &family);
    int baseSize() const { return m_baseSize; }
    void setBaseSize(int size);
    qreal lineHeight() const { return m_lineHeight; }
    void setLineHeight(qreal height);
    int paragraphSpacing() const { return m_paragraphSpacing; }
    void setParagraphSpacing(int spacing);
    int maxContentWidth() const { return m_maxContentWidth; }
    void setMaxContentWidth(int width);
    QString monoFamily() const { return m_monoFamily; }
    void setMonoFamily(const QString &family);

    // The type-scale: the pixel size a font role renders at, derived from the
    // base size by the frozen ratios. The role is a FontRole, handed across as
    // an int because that is what a model role carries into QML; a block's
    // kind answers which one it is.
    //
    // This used to switch over Block::BlockType, which made it one more place
    // that had to learn about a new kind of block, and it learned by silently
    // giving anything it did not recognise the body size.
    Q_INVOKABLE int sizeForRole(int fontRole) const;

    int bodySize() const;
    int monoSize() const;

    // Fixed-pitch families for the mono picker (QFontDatabase).
    Q_INVOKABLE QStringList monospaceFamilies() const;

    // Reset every value to the built-in defaults (the dialog's
    // "Reset typography" button).
    Q_INVOKABLE void resetToDefaults();

signals:
    void typographyChanged();

private:
    void save(const QString &key, const QVariant &value);

    SettingsStore *m_settings = nullptr;
    bool m_loading = false;
    QString m_fontFamily;
    int m_baseSize = DefaultBaseSize;
    qreal m_lineHeight = DefaultLineHeight;
    int m_paragraphSpacing = DefaultParagraphSpacing;
    int m_maxContentWidth = 0;
    QString m_monoFamily = QStringLiteral("monospace");
};

#endif // TYPOGRAPHY_H
