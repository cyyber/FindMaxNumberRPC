package client

import (
	"context"
	"fmt"
	"github.com/cyyber/FindMaxNumberRPC/pkg/generated"
	"github.com/cyyber/FindMaxNumberRPC/pkg/mock_generated"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"testing"
	"time"
)

func TestRunClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stream := mock_generated.NewMockFindMaxNumber_GetMaxNumberClient(ctrl)

	stream.EXPECT().Send(
		gomock.Any(),
	).Return(nil)
	msg := &generated.Number{Number:10}
	stream.EXPECT().Recv().Return(msg, nil)
	stream.EXPECT().CloseSend().Return(nil)

	client := mock_generated.NewMockFindMaxNumberClient(ctrl)
	client.EXPECT().GetMaxNumber(
		gomock.Any(),
	).Return(stream, nil)
	if err := testClient(client); err != nil {
		t.Fatalf("Test Failed: %v", err)
	}
}

func testClient(client generated.FindMaxNumberClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
	defer cancel()

	stream, err := client.GetMaxNumber(ctx)
	if err != nil {
		return err
	}
	msg := &generated.Number{Number:10}
	if err := stream.Send(msg); err != nil {
		return err
	}
	if err := stream.CloseSend(); err != nil {
		return err
	}
	got, err := stream.Recv()
	if err != nil {
		return err
	}
	if !proto.Equal(got, msg) {
		return fmt.Errorf("stream.Recv() = %v, want %v", got, msg)
	}
	return nil
}
