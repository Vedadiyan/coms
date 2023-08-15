package common

type Events string

const (
	GROUP_JOIN  Events = Events("group:join")
	GROUP_LEAVE Events = Events("group:leave")
	EMIT_GROUP  Events = Events("emit:group")
	EMIT_SOCKET Events = Events("emit:socket")
	REPLY       Events = Events("reply")
)
