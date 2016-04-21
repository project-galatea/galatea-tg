package main

var (
	BotToken  string
	Version   string
	BuildTime string
)

func main() {
	ConnectNewFollower("127.0.0.1:24833")
	startBot(BotToken)
}
