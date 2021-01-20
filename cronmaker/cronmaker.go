package cronmaker

import (
	"github.com/jasonlvhit/gocron"
)

type CronJob struct {
	Task map[string]func(taskName string)
}

type ICronJob interface {
	StartCronJobs()
}

func (cj *CronJob) StartCronJobs() {
	for k, v := range cj.Task {
		v(k)
	}
	<-gocron.Start()
}
