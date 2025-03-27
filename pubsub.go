package main

type pubsub []chan string

func (ps pubsub) publish(msg string) {
	for _, ch := range ps {
		ch <- msg
	}
}

func (ps *pubsub) subscribe() <-chan string {
	ch := make(chan string)
	*ps = append(*ps, ch)
	return ch
}
