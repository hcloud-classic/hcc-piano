package server

import (
	"context"
	"hcc/piano/action/grpc/pb/rpcpiano"
	"hcc/piano/driver/influxdb"
)

type pianoServer struct {
	rpcpiano.UnimplementedPianoServer
}

func (s *pianoServer) Telegraph(ctx context.Context, in *rpcpiano.ReqMetricInfo) (*rpcpiano.ResMonitoringData, error) {
	series := influxdb.GetInfluxData(in)

	return series, nil
}
