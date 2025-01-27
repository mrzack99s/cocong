package watcher

import (
	"context"
	"time"

	"github.com/mrzack99s/cocong/session"
)

func NetIdleChecking(ctx context.Context) {

	// pubsub := vars.RedisCache.PSubscribe(ctx, "__keyevent@0__:expired")
	go func(ctx context.Context) {
		// for msg := range pubsub.Channel() {

		// 	key := msg.Payload
		// 	split := strings.Split(key, "|")
		// 	s := inmemory_model.Session{
		// 		ID:        fmt.Sprintf("%s|%s", split[1], split[2]),
		// 		IPAddress: split[2],
		// 		User:      split[1],
		// 	}

		// 	session.CutOffSession(s)
		// 	time.Sleep(100 * time.Millisecond)
		// }

		for {
			select {
			case <-time.After(500 * time.Millisecond):
			case <-ctx.Done():
				return
			default:
				allSessions := session.Instance.GetAllSession()

				for _, s := range allSessions {
					if expired, err := session.Instance.IsExpired(s.ID); err == nil && expired {
						session.Instance.Delete(s.IPAddress)
					}
				}
				time.Sleep(500 * time.Millisecond)

			}
		}
	}(ctx)
}
