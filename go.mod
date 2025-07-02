module github.com/volodymyrzuyev/goCsInspect

go 1.24.4

require (
	github.com/Philipp15b/go-steam/v3 v3.0.0
	github.com/bbqtd/go-steam-authenticator v0.0.0-20160724194112-c5890fde0935
	github.com/davecgh/go-spew v1.1.1
)

require (
	golang.org/x/net v0.9.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

replace github.com/Philipp15b/go-steam/v3 v3.0.0 => github.com/csfloat/go-steam/v3 v3.0.12
