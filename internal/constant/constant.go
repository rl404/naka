package constant

// Color.
const (
	ColorGreyLight = 12370112
	ColorBlue      = 3447003
)

// Message template.
const (
	MsgHelpCmd       = "{{prefix}}help"
	MsgHelpContent   = "Naka is a bot for playing song from youtube."
	MsgPlayCmd       = "Play Song"
	MsgPlayContent   = "```{{prefix}}play <song name|youtube url>``` ```{{prefix}}play ina unravel\n{{prefix}}play https://www.youtube.com/watch?v=dQw4w9WgXcQ```"
	MsgOtherCmd      = "Other Commands"
	MsgJoinCmd       = "{{prefix}}join"
	MsgJoinContent   = "join voice channel"
	MsgLeaveCmd      = "{{prefix}}leave"
	MsgLeaveContent  = "leave voice channel"
	MsgPauseCmd      = "{{prefix}}pause"
	MsgPauseContent  = "pause song"
	MsgResumeCmd     = "{{prefix}}resume"
	MsgResumeContent = "resume song"
	MsgNextCmd       = "{{prefix}}next"
	MsgNextContent   = "next song"
	MsgPrevCmd       = "{{prefix}}prev"
	MsgPrevContent   = "previous song"
	MsgStopCmd       = "{{prefix}}stop"
	MsgStopContent   = "stop song"
	MsgQueueCmd      = "{{prefix}}queue"
	MsgQueueContent  = "see queue list"
	MsgPurgeCmd      = "{{prefix}}pruge"
	MsgPurgeContent  = "purge queue list"

	MsgInvalid        = "Invalid command. See **{{prefix}}help** for more information."
	MsgNotInVC        = "You are not in a voice channel."
	MsgEmptyQueue     = "Queue is empty. Go add some songs."
	MsgEndQueue       = "End of queue."
	MsgInvalidYoutube = "Invalid youtube url."
)
