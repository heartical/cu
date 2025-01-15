module cu/server

go 1.23.1

require (
	cu/common v0.0.0-00010101000000-000000000000
	github.com/dgraph-io/badger/v4 v4.5.0
	github.com/google/uuid v1.1.2
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgraph-io/ristretto/v2 v2.0.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/ebitengine/purego v0.3.0 // indirect
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20221017161538-93cebf72946b // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/flatbuffers v24.3.25+incompatible // indirect
	github.com/hajimehoshi/ebiten/v2 v2.5.0 // indirect
	github.com/jezek/xgb v1.1.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/tinne26/etxt v0.0.8 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/exp v0.0.0-20190731235908-ec7cb31e5a56 // indirect
	golang.org/x/image v0.9.0 // indirect
	golang.org/x/mobile v0.0.0-20230301163155-e0f57694e12c // indirect
	golang.org/x/net v0.31.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

replace cu/common => ../common
