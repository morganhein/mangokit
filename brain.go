package mangokit

import (
	"github.com/morganhein/mangokit/log"
	"github.com/morganhein/mangokit/plugins"
	"github.com/morganhein/stacknqueue"
	"sync"
)

type MangoBrain struct {
	input            *stacknqueue.StackNQueue
	output           *stacknqueue.StackNQueue
	subscribedEvents map[int][]plugins.SkillPlugineers
}

const (
	QUIT = iota
	PAUSE
	RESUME
)

var Brain *MangoBrain

func init() {
	Brain = &MangoBrain{
		// give the brain a little memory
		input: stacknqueue.NewStackNQueue(true),
		output: stacknqueue.NewStackNQueue(true),
		// create the event types to skill association map
		subscribedEvents: make(map[int][]plugins.SkillPlugineers),
	}
}

func (b *MangoBrain) Loop() {
	// create the communication channels
	// anyone know of a better way to do this? it feels clunky
	controlListen := make(chan int)
	//controlThink := make(chan int)
	controlThought := make(chan int)
	controlSpeak := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(3)
	go b.listen(&wg, controlListen)
	//go b.think(&wg, controlThink)
	go b.thought(&wg, controlThought)
	go b.speak(&wg, controlSpeak)
	log.Debug("Waiting...")
	wg.Wait()
}

func (b *MangoBrain) AddEventTriggers(e []int, skill plugins.SkillPlugineers) {
	for _, v := range e {
		if pls, exists := b.subscribedEvents[v]; exists {
			b.subscribedEvents[v] = append(pls, skill)
		} else {
			b.subscribedEvents[v] = []plugins.SkillPlugineers{skill}
		}
	}
}

// listen takes new events from networks and stores them in short term memory to be processed by a skill
func (b *MangoBrain) listen(wg *sync.WaitGroup, control chan int) {
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
			// check for any new events                                                                                       e
			case event := <-c.FromPlugin:
				log.Debug("Received network event: " + event.Data)
			// add the connection this event came from onto the event
				event.Connection = c
				b.think(&event)
			}
		}
	}
}

// think figures out which skills plugins want this event
func (b *MangoBrain) think(e *plugins.Event) {
	for _, p := range plugins.SkillPlugins {
		go p.NewEvent(*e)
	}
}

// thought takes new responses from skills and adds them to short term memory to be sent to the networks
func (b *MangoBrain) thought(wg *sync.WaitGroup, control chan int) {
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
				log.Debug("Received skill response: " + event.Data)
			// add the event to be processed by the skills
				b.output.Push(event)
			}
		}
	}
	log.Debug("Thoughts lost.")
}

// speak takes thoughts and figures out where they need to be sent
func (b *MangoBrain) speak(wg *sync.WaitGroup, control chan int) {
	defer wg.Done()

	log.Debug("Mouthpiece warming up.")

	Loop:
	for {
		if next := b.output.Pop(); next != nil {
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

func notFull(c chan plugins.Event) bool {
	return len(c) < 10
}


