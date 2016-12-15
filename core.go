package mangokit

import (
	"github.com/morganhein/mangokit/log"
	"github.com/morganhein/mangokit/plugins"
	"github.com/morganhein/stacknqueue"
	"sync"
	"os"
	"github.com/morganhein/mangokit/events"
)

type core struct {
	input            *stacknqueue.StackNQueue
	output           *stacknqueue.StackNQueue
	subscribedEvents map[int][]plugins.SkillPlugineers
}

const (
	QUIT = iota
	PAUSE
	RESUME
)

var Core *core

func init() {
	Core = &core{
		// give the core a little memory
		input: stacknqueue.NewStackNQueue(true),
		output: stacknqueue.NewStackNQueue(true),
		// create the event types to skill association map
		subscribedEvents: make(map[int][]plugins.SkillPlugineers),
	}
	plugins.Core = Core
}

func (c *core) Loop() {
	// create the communication channels
	// anyone know of a better way to do this? it feels clunky
	controlListen := make(chan int)
	//controlThink := make(chan int)
	controlThought := make(chan int)
	controlSpeak := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(3)
	go c.listen(&wg, controlListen)
	//go b.think(&wg, controlThink)
	go c.thought(&wg, controlThought)
	go c.speak(&wg, controlSpeak)
	log.Debug("Waiting...")
	wg.Wait()
}

func (c *core) AddEventTriggers(e []int, skill plugins.SkillPlugineers) {
	for _, v := range e {
		if pls, exists := c.subscribedEvents[v]; exists {
			c.subscribedEvents[v] = append(pls, skill)
		} else {
			c.subscribedEvents[v] = []plugins.SkillPlugineers{skill}
		}
	}
}

// listen takes new events from networks and stores them in short term memory to be processed by a skill
func (b *core) listen(wg *sync.WaitGroup, control chan int) {
	defer wg.Done()

	log.Debug("Starting up listener.")

	Loop:
	for {
		for c, _ := range plugins.NetworkPlugins {
			select {
			// check if we have any new control messages
			case con := <-control:
				switch con {
				case QUIT:
					break Loop
				}
			// check for any new events
			case event := <-c.FromPlugin:
				log.Debug("Received network event: " + event.Raw)
			// add the connection this event came from onto the event
				event.Connection = c
			// figure out if this is a botcommand
				_ = plugins.PopulateCmd(&event)
			// Process this in case it's a low-level control command not handled by a plugin
			//	b.preProcess(&event)
				b.think(&event)
			}
		}
	}
}

// think figures out which skills plugins want this event
func (c *core) think(e *plugins.Event) {
	for _, p := range c.subscribedEvents[e.Type] {
		go p.NewEvent(*e)
	}
	//todo: a skill could potentially receive events twice+ if subbed to ALL+1
	for _, p := range c.subscribedEvents[events.ALL] {
		go p.NewEvent(*e)
	}

}

// thought takes new responses from skills and adds them to short term memory to be sent to the networks
func (c *core) thought(wg *sync.WaitGroup, control chan int) {
	defer wg.Done()

	log.Debug("Thoughts incoming.")

	Loop:
	for {
		for p, _ := range plugins.SkillPlugins {
			select {
			// check if we have any new control messages
			case con := <-control:
				switch con {
				case QUIT:
					break Loop
				}
			// check for any new events
			case event := <-p.FromPlugin:
				log.Debug("Received skill response: " + event.Message)
			// add the event to be processed by the skills
				c.output.Push(event)
			}
		}
	}
	log.Debug("Thoughts lost.")
}

// speak takes thoughts and figures out where they need to be sent
func (c *core) speak(wg *sync.WaitGroup, control chan int) {
	defer wg.Done()

	log.Debug("Mouthpiece warming up.")

	Loop:
	for {
		if next := c.output.Pop(); next != nil {
			response := next.(plugins.Event)
			// send the response back up the connection through the ToPlugin chan
			response.Connection.ToPlugin <- response
		}
		select {
		case con := <-control:
			switch con {
			case QUIT:
				break Loop
			}
		}
	}
	log.Debug("Shutting up.")
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


