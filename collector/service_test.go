package collector_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/kristaxox/unilog-server/collector"
	"github.com/kristaxox/unilog-server/pb"
	"google.golang.org/grpc"
)

type mockStream struct {
	Stream chan *pb.Log
	grpc.ServerStream
}

func (m *mockStream) Recv() (*pb.Log, error) {
	v, ok := <-m.Stream
	if v != nil {
		return v, nil
	}
	if v == nil && !ok {
		return nil, io.EOF
	}

	return nil, io.EOF
}

func (m *mockStream) SendAndClose(*empty.Empty) error {
	return nil
}

func TestLogRecord(t *testing.T) {
	mock := mockStream{
		Stream: make(chan *pb.Log),
	}

	go func() {
		for i := 0; i < 5; i++ {
			now := time.Now()
			mock.Stream <- &pb.Log{
				Id: fmt.Sprintf("id%d", i),
				CreatedAt: &timestamp.Timestamp{
					Nanos: int32(now.UnixNano()),
				},
				Body: fmt.Sprintf("body_%d", i),
			}
		}
		close(mock.Stream)
	}()

	c := collector.NewServer()
	c.Record(&mock)
}
