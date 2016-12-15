package events

const (
	CONNECTED = iota //Initial status when a connection to a server has been been completed
	READY //When all bootstrapping and data updates for a server have been received
	JOINEDSERVER //A member, or this bot, have joined a server
	EXITEDSERVER //A member, or this bot, have exited the server
	JOINEDCHANNEL //A member, or this bot, have entered a channel
	EXITEDCHANNEL //A member, or this bot, have exited a channel
	ANNOUNCEMENT //A server/guild wide announcement was received
	MESSAGE //A generic message, public/private is unknown
	BOTMESSAGE //A message the bot is supposed to reply to, from any context, with or without parameters
	PUBLICMESSAGE //A general message seen in a public context
	PRIVATEMESSAGE //A message seen in a private context
	ACTION //An action seen in any context
	FILE //A file being sent directly to bot
	ATTACHMENT //An attachment seen in any context
	STATUSUPDATE //The status of any member in any context has changed
	CONNECTIONUPDATE //Any updates pertaining to connectivity to the server
	GETMEMBERSUPDATE //When information on a specific member is updated/received
	BANNED //Any member, including this bot, is banned from the server
	KICKED //Any member, including this bot, is kicked from the server
	SHUTDOWN //Called just prior to a full shutdown of the server
)