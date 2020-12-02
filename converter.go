package mimir

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

func ToPbTimestamp(golangTime time.Time) *timestamp.Timestamp {
	timestamps, _ := ptypes.TimestampProto(golangTime)
	return timestamps
}

func PbTimestampToTime(golangTimestamp *timestamp.Timestamp) time.Time {
	times, _ := ptypes.Timestamp(golangTimestamp)
	return times
}
