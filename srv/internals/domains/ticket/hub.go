package ticketdomain

type Hub struct {
	In  chan<- string
	Out <-chan string
}

func NewHub() *Hub {
	in := make(chan string)
	out := make(chan string)
	return &Hub{
		In:  in,
		Out: out,
	}
}
