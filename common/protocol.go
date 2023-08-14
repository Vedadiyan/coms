package common

type Events string

const (
	ROOM_JOIN   Events = Events("room:join")
	ROOM_LEAVE  Events = Events("room:leave")
	EMIT_ROOM   Events = Events("emit:room")
	EMIT_SOCKET Events = Events("emit:socket")
)
