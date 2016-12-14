package mangokit

import (
	"time"
	"github.com/morganhein/mangokit/log"
	"github.com/morganhein/mangokit/plugins"
)

type Brain struct{}

func (b *Brain) Loop() {
	for {
		for c, _ := range plugins.NetworkPlugins {
			select {
			case event := <-c.FromPlugin:
				log.Info("Received " + event.Data)
			}
		}

		time.Sleep(2 * time.Second)
	}
}


