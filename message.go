package main

type Message struct {
	From    string
	Payload any
}

type MessageStoreFile struct {
	ID   string
	Key  string
	Size int64
}

type MessageGetFile struct {
	ID  string
	Key string
}

type MessageDeleteFile struct {
	ID  string
	Key string
}
