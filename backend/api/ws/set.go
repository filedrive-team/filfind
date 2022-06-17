package ws

type ClientSet struct {
	m map[*Client]struct{}
}

func NewClientSet(clients ...*Client) *ClientSet {
	s := &ClientSet{
		m: make(map[*Client]struct{}),
	}
	s.Add(clients...)
	return s
}

func (s *ClientSet) Add(clients ...*Client) {
	for _, client := range clients {
		s.m[client] = struct{}{}
	}
}

func (s *ClientSet) Remove(clients ...*Client) {
	for _, client := range clients {
		delete(s.m, client)
	}
}

func (s *ClientSet) Contains(client *Client) bool {
	_, ok := s.m[client]
	return ok
}

func (s *ClientSet) Map(f func(cli *Client)) {
	for client := range s.m {
		f(client)
	}
}

func (s *ClientSet) Clear() {
	s.m = make(map[*Client]struct{})
}

func (s *ClientSet) Size() int {
	return len(s.m)
}

type ChannelSet struct {
	m map[Channel]struct{}
}

func NewChannelSet(channels ...Channel) *ChannelSet {
	s := &ChannelSet{
		m: make(map[Channel]struct{}),
	}
	s.Add(channels...)
	return s
}

func (s *ChannelSet) Add(channels ...Channel) {
	for _, channel := range channels {
		s.m[channel] = struct{}{}
	}
}

func (s *ChannelSet) Remove(channels ...Channel) {
	for _, channel := range channels {
		delete(s.m, channel)
	}
}

func (s *ChannelSet) Contains(channel Channel) bool {
	_, ok := s.m[channel]
	return ok
}

func (s *ChannelSet) Map(f func(channel Channel)) {
	for channel := range s.m {
		f(channel)
	}
}

func (s *ChannelSet) Clear() {
	s.m = make(map[Channel]struct{})
}

func (s *ChannelSet) Size() int {
	return len(s.m)
}
