# go-text/typesetting v0.3.5, with one change

This is go-text/typesetting v0.3.5 (`github.com/go-text/typesetting`),
without its tests and test data, used by kvit-ui-go through a `replace`
line in `go.mod`. Every app that uses kvit-ui-go needs the same line:

    replace github.com/go-text/typesetting => ../kvit-ui-go/third_party/typesetting

**The change.** `fontscan.FontMap` gains `SetFaceLoader`, which lets a
caller decide how a system font is loaded once the map has chosen it. The
diff is `../typesetting-fontscan-face-loader.patch`; nothing else differs
from v0.3.5.

**Why.** go-text parses a font file in full, glyph outlines and colour
tables included, and keeps it all for the life of the face. Kvit's text
package only shapes with these faces and draws glyphs through canvas, which
parses its own copy. Without the change, every font in use was parsed in
full twice. On Windows the kvit-ui gallery held 130 MB of parsed fonts
(2026-09-26). With it, the text package gives the map faces parsed from
copies without outline, bitmap and colour tables, which cost 80–97% less
and shape identically.

**Upstream.** The change is written to be offered to go-text as it stands.
It applies unchanged to go-text's main branch as of 2026-09-22 (commit
`e8da365`), where `fontscan/fontmap.go` is the same as in v0.3.5. When a
go-text release has it (or something equivalent), delete this directory and
the `replace` lines, and require that release instead.

**Test.** `fontscan/faceloader_test.go` is the one test file kept here, and
is part of the patch: the loader is used once per font chosen, and removing
it restores loading from disk.

**Draft of the offer**, for a pull request or an issue at
github.com/go-text/typesetting:

> **fontscan: let the caller load the faces a FontMap chooses**
>
> `FontMap` always loads a chosen system font by reading and parsing the
> whole file (`Footprint.loadFromDisk`). An application that uses the map
> only to pick fonts for shaping, and draws glyphs with a separate
> rasterizer holding its own parse, pays for every font twice. go-text
> parses glyph outlines and colour tables eagerly and keeps them for the
> life of the face, so this is large. In our case, a desktop toolkit
> drawing through a Skia-like rasterizer, a single window on Windows held
> 130 MB of parsed fonts, most of it Segoe UI Emoji's COLRv1 paint tables
> and glyph outlines that shaping never reads.
>
> This adds `FontMap.SetFaceLoader(func(Location) (*font.Face, error))`.
> When set, it is called instead of loading from disk once the map has
> chosen a font; faces are still cached per location as before, and
> passing nil restores the default. With a loader that parses a copy of
> the font without `glyf`/`loca`/`CFF `/`COLR`/`CPAL`/`CBDT`/`CBLC`/`sbix`/`SVG `,
> the cost per face fell by 80–97% (DejaVu Sans 2.6 MB → 0.5 MB, Noto
> Color Emoji 9.6 MB → 0.3 MB) with identical shaping output.
>
> One thing a caller doing this has to know, which could be documented
> with the option: `LineWrapper` decides that a glyph at the end of a line
> is a space by its ink width (`Glyph.Width == 0`). Faces without outlines
> report zero ink width for every glyph, so such a caller has to set a
> width on non-space glyphs after shaping, or the last glyph of a line
> loses its advance.

**Updating go-text.** Copy the new release over this directory, reapply the
patch, and run `tools/check-all.sh`.
