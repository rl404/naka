package entity

// Color.
const (
	ColorGreyLight = 12370112
	ColorBlue      = 3447003
)

// Message template.
const (
	HelpCmd    = "{{prefix}}help"
	HelpDesc   = "Naka is a bot for playing song from youtube."
	PlayCmd    = "{{prefix}}play"
	PlayDesc   = "```{{prefix}}play <song name|youtube url>``` ```{{prefix}}play ina unravel\n{{prefix}}play https://www.youtube.com/watch?v=dQw4w9WgXcQ```"
	JoinCmd    = "{{prefix}}join"
	JoinDesc   = "join voice channel"
	LeaveCmd   = "{{prefix}}leave"
	LeaveDesc  = "leave voice channel"
	PauseCmd   = "{{prefix}}pause"
	PauseDesc  = "pause song"
	ResumeCmd  = "{{prefix}}resume"
	ResumeDesc = "resume song"
	NextCmd    = "{{prefix}}next"
	NextDesc   = "next song"
	PrevCmd    = "{{prefix}}prev"
	PrevDesc   = "previous song"
	StopCmd    = "{{prefix}}stop"
	StopDesc   = "stop song"
	QueueCmd   = "{{prefix}}queue"
	QueueDesc  = "add or see queue list"
	RemoveCmd  = "{{prefix}}remove"
	RemoveDesc = "remove song from queue"
	PurgeCmd   = "{{prefix}}purge"
	PurgeDesc  = "purge queue list"

	EmptyQueue         = "Queue is empty. Go add some songs."
	EndQueue           = "End of queue."
	NotInVC            = "You are not in voice channel."
	InvalidYoutubeURL  = "Invalid youtube url."
	InvalidSearchQuery = "Invalid search query."
	PromptExpired      = "Prompt has expired."
	InvalidPrompt      = "Invalid response."
	Paused             = "Paused."
	Resumed            = "Resumed."
	Stopped            = "Stopped."
	Next               = "Next."
	Previous           = "Previous."
	InvalidQueue       = "Invalid queue."
	Purged             = "Queue has been purged."
)
