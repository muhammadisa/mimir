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

func TimeServerFormat(times time.Time) string {
	return times.UTC().Format("2006-01-02T15:04:05-0700")
}

func TimeSimpleFormat(times time.Time) string {
	return times.UTC().Format("2006-01-02 15:04:05")
}