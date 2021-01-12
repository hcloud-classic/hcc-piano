module hcc/piano

go 1.13

require (
	github.com/Scalingo/go-utils v5.6.1+incompatible
	github.com/Terry-Mao/goconf v0.0.0-20161115082538-13cb73d70c44
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/hcloud-classic/hcc_errors v1.1.3
	github.com/hcloud-classic/pb v0.0.0
	github.com/influxdata/influxdb v1.8.3 // indirect
	github.com/influxdata/influxdb1-client v0.0.0-20190809212627-fc22c7df067e
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b // indirect
	golang.org/x/sys v0.0.0-20210110051926-789bb1bd4061 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210111234610-22ae2b108f89 // indirect
	google.golang.org/grpc v1.34.1
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v0.0.0-20200812184716-7d8921505e1b // indirect
	google.golang.org/protobuf v1.25.0 // indirect
)

replace github.com/hcloud-classic/pb => ../pb
