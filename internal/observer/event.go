package observer

type EventType uint8
type EventLevel uint8

const (
	RequestLog EventType = iota
	SystemLog
)

const (
	Info EventLevel = iota
	Warning
	Error
	Critical
)

type Event struct {
	Type    EventType
	Lvl     EventLevel
	From    string
	Message string
	Err     error
	Ctx     any
}

func NewRequestLogEvent(lvl EventLevel, from, message string, ctx any, err error) Event {
	return Event{
		Type:    RequestLog,
		Lvl:     lvl,
		From:    from,
		Message: message,
		Err:     err,
		Ctx:     ctx,
	}
}

func NewSystemLogEvent(lvl EventLevel, from, message string, ctx any, err error) Event {
	return Event{
		Type:    SystemLog,
		Lvl:     lvl,
		From:    from,
		Message: message,
		Err:     err,
		Ctx:     ctx,
	}
}
