package arp

type arpStats struct {
	macaddr      string
	asked        int
	answered     int
	askedList    []string
	answeredList []string
}
