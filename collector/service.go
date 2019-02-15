package collector

import (
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/kristaxox/unilog-server/pb"
	"github.com/sirupsen/logrus"
)

type collectorServer struct{}

// NewServer returns a collectorServer
// nolint
func NewServer() *collectorServer {
	return &collectorServer{}
}

func (c *collectorServer) Record(stream pb.LogCollector_RecordServer) error {
	for {
		log, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&empty.Empty{})
		}
		if err != nil {
			return err
		}
		logrus.Println(log)
	}
}
