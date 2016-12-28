package mangokit

import (
	"os"
	"sync"

	"github.com/morganhein/mangokit/events"
	"github.com/morganhein/mangokit/plugins"
	"github.com/morganhein/stacknqueue"
)

type core struct {
	output           *stacknqueue.StackNQueue
	subscribedEvents map[int][]plugins.Plugineer
}

const (
	QUIT = iota
	PAUSE
	RESUME
)

var Core *core

func init() {
	Core = &core{
		output: stacknqueue.NewStackNQueue(true),
		// create the event types to skill association map
		subscribedEvents: make(map[int][]plugins.Plugineer),
	}
	plugins.Core = Core
}

func (c *core) Loop() {
	// create the communication channels
	// anyone know of a better way to do this? it feels clunky
	ctrlReceive := make(chan int)
	//controlThink := make(chan int)
	ctrlBufferResponse := make(chan int)
	ctrlRespond := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(3)
	go c.receive(&wg, ctrlReceive)
	// go b.think(&wg, controlThink)
	go c.bufferResponse(&wg, ctrlBufferResponse)
	go c.respond(&wg, ctrlRespond)
	log.Debug("Core is ready...")
	wg.Wait()
}

func (c *core) AddEventTriggers(p plugins.Plugineer) {
	log.Debug("Adding event triggers.")
	for _, e := range p.Events() {
		pls, exists := c.subscribedEvents[e]
		if !exists {
			c.subscribedEvents[e] = []plugins.Plugineer{p}
		}
		c.subscribedEvents[e] = append(pls, p)
	}
}

// receive takes new events from networks and stores them in short term memory to be processed by a skill
func (c *core) receive(wg *sync.WaitGroup, control chan int) {
	defer wg.Done()
	log.Debug("Starting up receiver.")
Loop:
	for {
		for _, pl := range plugins.NetworkPlugins {
			select {
			// check if we have any new control messages
			case con := <-control:
				switch con {
				case QUIT:
					break Loop
				}
			// check for any new events
			case event := <-pl.FromPlugin():
				log.Debug("Received network event: " + event.Raw)
				// figure out if this is a botcommand
				_ = plugins.PopulateCmd(&event)
				// Process this in case it's a low-level control command not handled by a plugin
				//	b.preProcess(&event)
				go c.distribute(&event)
			}
		}
	}
}

// distribute figures out which skills plugins want this event
func (c *core) distribute(e *plugins.Event) {
	log.Debug("Distributing event: " + e.Raw)
	//todo: a skill could potentially receive events twice+ if subbed to ALL+1
	for _, p := range c.subscribedEvents[events.ALL] {
		if notFull(p.ToPlugin()) {
			log.Debug("Sending event to: " + p.Name())
			p.ToPlugin() <- *e
		}
	}
	for _, p := range c.subscribedEvents[e.Type] {
		if notFull(p.ToPlugin()) {
			log.Debug("Sending event to: " + p.Name())
			p.ToPlugin() <- *e
		}
	}

}

// bufferResponse takes new responses from skills and adds them to short term memory to be sent to the networks
func (c *core) bufferResponse(wg *sync.WaitGroup, control chan int) {
	defer wg.Done()

	log.Debug("Starting up responder queue.")

Loop:
	for {
		for _, pl := range plugins.SkillPlugins {
			select {
			// check if we have any new control messages
			case con := <-control:
				switch con {
				case QUIT:
					break Loop
				}
			// check for any new events
			case event := <-pl.FromPlugin():
				log.Debug("Received skill response: " + event.Message)
				// add the event to be processed by the skills
				c.output.Push(event)
			}
		}
	}
	log.Debug("Responder queue shutting down.")
}

// respond sends out the responses from the response buffer
func (c *core) respond(wg *sync.WaitGroup, control chan int) {
	defer wg.Done()

	log.Debug("Responder starting up.")

Loop:
	for {
		if next := c.output.Pop(); next != nil {
			event := next.(plugins.Event)
			// send the event back up the connection through the ToPlugin chan
			event.Source.ToPlugin() <- event
		}
		select {
		case con := <-control:
			switch con {
			case QUIT:
				break Loop
			}
		}
	}
	log.Debug("Going silent. Responder shut down.")
}

func (c *core) Leave(pc plugins.Contexter) {
	//todo: implement this stub
}

func (c *core) Quit() {
	// send disconnect commands to all networks
	for _, p := range plugins.NetworkPlugins {
		p.Disconnect()
	}
	os.Exit(0)
}

func notFull(c chan plugins.Event) bool {
	return len(c) < 10
}
