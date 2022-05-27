package alert

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/zerodoctor/zdcli/logger"
)

type Trouble struct {
	title   string
	message string
}

func NewTrouble(title, message string) *Trouble {
	return &Trouble{
		title:   title,
		message: message,
	}
}

type Alert struct {
	action func() *Trouble
	ctx    context.Context

	checkDur time.Duration
	tick     *time.Ticker
}

func NewAlert(ctx context.Context, checkDur time.Duration, action func() *Trouble) *Alert {
	a := &Alert{
		action:   action,
		checkDur: checkDur,
		ctx:      ctx,
	}

	go a.Listener()

	return a
}

func (a *Alert) Listener() {
	a.tick = time.NewTicker(a.checkDur)

	for {
		select {
		case <-a.ctx.Done():
			logger.Info("stopping alert...")
			a.tick.Stop()
			return
		case <-a.tick.C:
			trouble := a.action()
			if trouble != nil {
				if err := beeep.Notify(trouble.title, trouble.message, ""); err != nil {
					logger.Errorf("failed to send notification", err.Error())
				}
			}
		}
	}
}

func (a *Alert) Done() { a.tick.Stop() }

func WatchEndpoint(ctx context.Context, healthRoute string, message string) *Alert {
	title := fmt.Sprintf("Err With [Endpoint=%s]", healthRoute)

	return NewAlert(ctx, 5*time.Second, func() *Trouble {
		rsp, err := http.Get(healthRoute)
		if err != nil {
			return NewTrouble(title, fmt.Sprintf("[Message=%s] [Error=%s]", message, err.Error()))
		}

		if rsp != nil && (rsp.StatusCode < 200 || rsp.StatusCode > 299) {
			return NewTrouble(title, fmt.Sprintf("[Message=%s] [Status=%s]", message, rsp.Status))
		}

		return nil
	})
}
