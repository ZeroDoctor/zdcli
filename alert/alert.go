package alert

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/zerodoctor/beeep"
	"github.com/zerodoctor/zdcli/logger"
)

type Trouble struct {
	title   string
	message string
	options []beeep.Option
}

func NewTrouble(title, message string, options ...beeep.Option) *Trouble {
	return &Trouble{
		title:   title,
		message: message,
		options: options,
	}
}

type Alert struct {
	action func() *Trouble
	ctx    context.Context
	wg     sync.WaitGroup

	checkDur time.Duration
	tick     *time.Ticker
}

func NewAlert(ctx context.Context, checkDur time.Duration, action func() *Trouble) *Alert {
	a := &Alert{
		action:   action,
		checkDur: checkDur,
		ctx:      ctx,
	}

	a.wg.Add(1)
	go a.Listener()

	return a
}

func (a *Alert) Listener() {
	a.tick = time.NewTicker(a.checkDur)
	defer a.tick.Stop()
	defer a.wg.Done()

	for {
		select {
		case <-a.ctx.Done():
			logger.Info("stopping alert...")
			return
		case <-a.tick.C:
			trouble := a.action()
			if trouble == nil {
				continue
			}

			trouble.options = append(trouble.options, beeep.AppOption(trouble.title))
			trouble.options = append(trouble.options, beeep.MessageOption(trouble.message))
			if trouble != nil {
				if err := beeep.Notify(trouble.options...); err != nil {
					logger.Errorf("failed to send notification", err.Error())
				}
			}
		}
	}
}

func (a *Alert) Done() { a.tick.Stop() }
func (a *Alert) Wait() { a.wg.Wait() }

func WatchEndpoint(ctx context.Context, healthRoute string, message string, checkDur time.Duration, options ...beeep.Option) *Alert {
	title := fmt.Sprintf("Err With [Endpoint=%s]", healthRoute)

	return NewAlert(ctx, checkDur, func() *Trouble {
		rsp, err := http.Get(healthRoute)
		if err != nil {
			return NewTrouble(title, fmt.Sprintf("[Message=%s] [Error=%s]", message, err.Error()), options...)
		}

		if rsp != nil && (rsp.StatusCode < 200 || rsp.StatusCode > 299) {
			return NewTrouble(title, fmt.Sprintf("[Message=%s] [Status=%s]", message, rsp.Status), options...)
		}

		return nil
	})
}
