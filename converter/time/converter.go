package time

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type CustomTime struct {
	// TimeZone constant
	TimeZone string
	// Adjustment this UTC hour addition
	Adjustment time.Duration
}

func (ct *CustomTime) CurrentTime() time.Time {
	location, _ := time.LoadLocation(ct.TimeZone)
	return time.Now().In(location).Add(ct.Adjustment)
}

func (ct *CustomTime) Time(t time.Time) time.Time {
	location, _ := time.LoadLocation(ct.TimeZone)
	return t.In(location).Add(ct.Adjustment)
}

func ToTimestampPb(golangTime time.Time) *timestamppb.Timestamp {
	return timestamppb.New(golangTime)
}

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
