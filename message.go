package main

type Message struct {
	From    string
	Payload any
}

type MessageStoreFile struct {
	Key  string
	Size int64
	Data []byte
}

type MessageGetFile struct {
	Key string
}
