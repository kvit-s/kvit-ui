module github.com/kvit-s/kvit-ui

go 1.27.0

toolchain go1.27.1

require (
	github.com/go-text/typesetting v0.3.5
	github.com/richardwilkes/canvas v0.3.1
	github.com/richardwilkes/toolbox/v2 v2.20.0
	github.com/richardwilkes/unison v0.108.0
	golang.org/x/image v0.46.0
	golang.org/x/sys v0.48.0
)

require (
	github.com/HugoSmits86/nativewebp v1.3.0 // indirect
	github.com/ebitengine/purego v0.11.0 // indirect
	github.com/klauspost/cpuid/v2 v2.4.0 // indirect
	github.com/yuin/goldmark v1.8.6 // indirect
	github.com/zeebo/xxh3 v1.1.0 // indirect
	golang.org/x/net v0.59.0 // indirect
	golang.org/x/text v0.42.0 // indirect
)

replace github.com/go-text/typesetting => ./third_party/typesetting
