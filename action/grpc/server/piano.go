package server

import (
	"context"
	"github.com/hcloud-classic/pb"
	"hcc/piano/driver/influxdb"
)

type pianoServer struct {
	pb.UnimplementedPianoServer
}

func (s *pianoServer) Telegraph(ctx context.Context, in *pb.ReqMetricInfo) (*pb.ResMonitoringData, error) {
	series := influxdb.GetInfluxData(in)

	return series, nil
}
