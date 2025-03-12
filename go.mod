module github.com/volodymyrzuyev/goCsInspect

go 1.24.0

require (
	github.com/Philipp15b/go-steam/v3 v3.0.0
	github.com/bbqtd/go-steam-authenticator v0.0.0-20160724194112-c5890fde0935
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/Philipp15b/go-steam/v3 v3.0.0 => github.com/csfloat/go-steam/v3 v3.0.11
