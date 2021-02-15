package helper

import "time"

type TimerFunc func(params interface{})

func SetInterval(fn TimerFunc, d time.Duration, params ...interface{}) {
	go func() {
		t := time.NewTimer(d)
		defer t.Stop()

		for {
			<-t.C
			fn(params)

			t.Reset(d)
		}
	}()
}
