package entity

// Color.
const (
	ColorGreyLight = 12370112
	ColorBlue      = 3447003
)

// Message template.
const (
	HelpCmd     = "{{prefix}}help"
	HelpDesc    = "Naka is a bot for playing song from youtube."
	PlaySample  = "```{{prefix}}play <song name|youtube url>``` ```{{prefix}}play ina unravel\n{{prefix}}play https://www.youtube.com/watch?v=dQw4w9WgXcQ```"
	QueueSample = "```{{prefix}}queue <song name|youtube url>``` ```{{prefix}}queue ina unravel\n{{prefix}}queue https://www.youtube.com/watch?v=dQw4w9WgXcQ```"
	PlayCmd     = "{{prefix}}play"
	PlayDesc    = "play queued songs"
	QueueCmd    = "{{prefix}}queue"
	QueueDesc   = "see queue list"
	JoinCmd     = "{{prefix}}join"
	JoinDesc    = "join voice channel"
	LeaveCmd    = "{{prefix}}leave"
	LeaveDesc   = "leave voice channel"
	PauseCmd    = "{{prefix}}pause"
	PauseDesc   = "pause song"
	ResumeCmd   = "{{prefix}}resume"
	ResumeDesc  = "resume song"
	NextCmd     = "{{prefix}}next"
	NextDesc    = "go to next song"
	PrevCmd     = "{{prefix}}prev"
	PrevDesc    = "go to previous song"
	StopCmd     = "{{prefix}}stop"
	StopDesc    = "stop song"
	SkipCmd     = "{{prefix}}skip 2"
	SkipDesc    = "skip to song number 2 in queue"
	RemoveCmd   = "{{prefix}}remove 1 2"
	RemoveDesc  = "remove song number 1 & 2 from queue"
	PurgeCmd    = "{{prefix}}purge"
	PurgeDesc   = "purge queue list"

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
