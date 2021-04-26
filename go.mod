module github.com/sardinasystems/sensu-go-chrony-check

go 1.15

replace go.etcd.io/etcd => go.etcd.io/etcd v0.0.0-20210226220824-aa7126864d82

require (
	github.com/gocarina/gocsv v0.0.0-20210408192840-02d7211d929d
	github.com/pelletier/go-toml v1.9.0 // indirect
	github.com/sardinasystems/sensu-go-haproxy-check v0.0.0-20210326111919-f2cd32992d01
	github.com/sensu/sensu-go/api/core/v2 v2.8.0 // indirect
	github.com/sensu/sensu-go/types v0.6.0
	github.com/sensu/sensu-plugin-sdk v0.13.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/multierr v1.6.0
	golang.org/x/net v0.0.0-20210423184538-5f58ad60dda6 // indirect
	golang.org/x/sys v0.0.0-20210423185535-09eb48e85fd7 // indirect
	google.golang.org/genproto v0.0.0-20210423144448-3a41ef94ed2b // indirect
	google.golang.org/grpc v1.37.0 // indirect
)
