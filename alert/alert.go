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
	once    bool
	options []beeep.Option
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
		}

		trouble := a.action()
		if trouble == nil {
			continue
		}

		trouble.options = append(trouble.options, beeep.AppOption(trouble.title))
		trouble.options = append(trouble.options, beeep.MessageOption(trouble.message))
		if trouble != nil {
			logger.Warnf("send notification [title=%s] [message=%s]", trouble.title, trouble.message)
			if err := beeep.Notify(trouble.options...); err != nil {
				logger.Errorf("failed to send notification", err.Error())
			}
		}

		if trouble.once {
			break
		}
	}
}

func (a *Alert) Done() { a.tick.Stop() }
func (a *Alert) Wait() { a.wg.Wait() }

type WatchEndpointParams struct {
	Ctx         context.Context
	HealthRoute string
	Message     string
	CheckDur    time.Duration
	Options     []beeep.Option
	Once        bool
}

func WatchEndpoint(params WatchEndpointParams) *Alert {
	defer func() {
		logger.Infof("alert on endpoint created for [route=%s]", params.HealthRoute)
	}()

	title := fmt.Sprintf("Err With [Endpoint=%s]", params.HealthRoute)

	return NewAlert(params.Ctx, params.CheckDur, func() *Trouble {
		rsp, err := http.Get(params.HealthRoute)

		if err != nil {
			return &Trouble{
				title:   title,
				message: fmt.Sprintf("[Message=%s] [Status=%s]", params.Message, "nil"),
				options: params.Options,
				once:    params.Once,
			}
		}

		if rsp != nil && (rsp.StatusCode < 200 || rsp.StatusCode > 299) {
			return &Trouble{
				title:   title,
				message: fmt.Sprintf("[Message=%s] [Status=%s]", params.Message, rsp.Status),
				options: params.Options,
				once:    params.Once,
			}
		}

		return nil
	})
}
