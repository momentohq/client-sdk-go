package main

import (
	"context"
	"fmt"
	"github.com/loov/hrtime"
	"github.com/momentohq/client-sdk-go/config/logger"
	"strconv"
	"strings"
	"time"
)

type TimingMiddleware struct {
	Log       logger.MomentoLogger
	timerChan chan string
}

func timer(timerChan chan string, log logger.MomentoLogger) {
	startTimes := make(map[uint64]int64)
	for {
		select {
		case timingMsg := <-timerChan:
			res := strings.Split(timingMsg, ":")
			operation := res[0]
			requestId, _ := strconv.ParseUint(res[1], 10, 64)
			timePoint, _ := strconv.ParseInt(res[2], 10, 64)
			if operation == "start" {
				startTimes[requestId] = timePoint
				continue
			}
			// we got an "end" message
			elapsed := timePoint - startTimes[requestId]
			log.Info(
				fmt.Sprintf(
					"Request %d took %dms", requestId, time.Duration(elapsed).Milliseconds(),
				),
			)
		}
	}
}

func (mw *TimingMiddleware) OnRequest(requestId uint64, theRequest interface{}, metadata context.Context) {
	mw.timerChan <- fmt.Sprintf("start:%d:%d", requestId, hrtime.Now())
}

func (mw *TimingMiddleware) OnResponse(requestId uint64, theResponse map[string]string) {
	mw.timerChan <- fmt.Sprintf("end:%d:%d", requestId, hrtime.Now())
}

func NewTimingMiddleware(log logger.MomentoLogger) *TimingMiddleware {
	mw := &TimingMiddleware{
		Log:       log,
		timerChan: make(chan string),
	}
	go func() {
		timer(mw.timerChan, mw.Log)
	}()
	return mw
}
