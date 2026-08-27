# Third-party notices

`kvit-ui` is MPL-2.0 (see `LICENSE`). It ships one third-party asset.

## Phosphor Icons

`resources/fonts/Phosphor.ttf` is the regular weight of the Phosphor icon font,
used by `KvitIcon` and therefore by every component in the library that draws a
symbol. The full licence text is `resources/fonts/Phosphor-LICENSE.txt`.

    Copyright (c) 2020 Phosphor Icons
    Licensed under the MIT License.

The MIT licence places no obligation on a consumer beyond carrying the notice,
so an application linking this library may ship under MPL-2.0 or under a
proprietary licence without the font constraining either.

Only the regular weight is shipped. Phosphor publishes six weights at roughly
480 KB each; a second weight is added when a caller needs one and not before,
so that the size a consumer pays is the size something in the estate draws
with.
