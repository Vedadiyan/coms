package common

type Events string

const (
	JOIN_GROUP  Events = Events("join:group")
	LEAVE_GROUP Events = Events("leave:group")
	EMIT_GROUP  Events = Events("emit:group")
	EMIT_SOCKET Events = Events("emit:socket")
	REPLY       Events = Events("reply")
	WHO_AM_I    Events = Events("whoami")
)
