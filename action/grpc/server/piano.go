package server

import (
	"context"
	"hcc/piano/driver/influxdb"

	"innogrid.com/hcloud-classic/pb"
)

type pianoServer struct {
	pb.UnimplementedPianoServer
}

func (s *pianoServer) Telegraph(ctx context.Context, in *pb.ReqMetricInfo) (*pb.ResMonitoringData, error) {
	series := influxdb.GetInfluxData(in)

	return series, nil
}
