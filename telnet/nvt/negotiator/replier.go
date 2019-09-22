package nego

import (
	"container/list"
	"context"
	"log"
	"sync"

	"github.com/wizcas/mudever.svc/telnet/packet"
)

// replier buffer the negotiator's replies and send
// them to the output channel, so negotiator won't be
// blocked for taking inputs in case of a jammed output channel.
type replier struct {
	sync.Mutex
	chOutput chan<- packet.Packet
	queue    *list.List
}

func newReplier(chOutput chan<- packet.Packet) *replier {
	return &replier{
		chOutput: chOutput,
		queue:    list.New(),
	}
}

func (r *replier) dispose() {
	r.queue.Init()
}

func (r *replier) enqueue(reply packet.Packet) {
	r.Lock()
	defer r.Unlock()
	r.queue.PushBack(reply)
}

func (r *replier) dequeue() packet.Packet {
	r.Lock()
	defer r.Unlock()
	e := r.queue.Front()
	if e == nil {
		return nil
	}
	r.queue.Remove(e)
	reply, ok := e.Value.(packet.Packet)
	if !ok {
		log.Printf("[NEGO RPLY ERR]: %v", e.Value)
		return nil
	}
	log.Printf("\x1b[36m<NEGO RPLY>\x1b[0m %s\n", reply)
	return reply
}

func (r *replier) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			r.dispose()
			return
		default:
			p := r.dequeue()
			if p != nil {
				r.chOutput <- p
			}
		}
	}
}
