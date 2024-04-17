package main

type Message struct {
	From    string
	Payload MessageData
}

type MessageData struct {
	Key  string
	Data []byte
}
