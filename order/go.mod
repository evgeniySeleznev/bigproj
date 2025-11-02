module github.com/evgeniyseleznev/bigproj/order

go 1.24.0

replace github.com/evgeniyseleznev/bigproj/shared => ../shared

require (
	github.com/evgeniyseleznev/bigproj/shared v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.2.3
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.76.0
)

require (
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
