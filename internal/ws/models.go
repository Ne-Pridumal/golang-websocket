package ws

type subscriber struct {
	msgs      chan []byte
	closeSlow func()
}
