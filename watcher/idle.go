package watcher

import (
	"context"
	"time"

	"github.com/mrzack99s/cocong/model/inmemory_model"
	"github.com/mrzack99s/cocong/session"
	"github.com/mrzack99s/cocong/vars"
)

func NetIdleChecking(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				return
			default:
				now := time.Now().In(vars.TZ)
				diff30 := now.Add(time.Minute * (-time.Duration(vars.Config.SessionIdle)))

				allIdles := []inmemory_model.Session{}
				vars.InMemoryDatabase.Where("last_seen <= ?", diff30).Find(&allIdles)

				for _, p := range allIdles {
					err := session.CutOffSession(p)
					if err != nil {
						vars.SystemLog.Println(err.Error())
					}
				}
				time.Sleep(500 * time.Millisecond)

			}
		}
	}(ctx)
}
