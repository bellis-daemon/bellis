package redistream

type MessageHandler func(message *Message) error
