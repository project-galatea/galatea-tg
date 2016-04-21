package main

import (
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

	"gopkg.in/telegram-bot-api.v4"
)

type Chat struct {
	chatId int
}

type Follower struct {
	Chats []Chat
	Conn  net.Conn
}

var followerList []*Follower

var chatsToFollower = make(map[int64]*Follower)
var followersLock = &sync.Mutex{}

func connectNewChat(chatid int64) {
	followersLock.Lock()
	rand.Seed(time.Now().Unix())
	setFollower := followerList[rand.Intn(len(followerList))]
	chatsToFollower[chatid] = setFollower
	followersLock.Unlock()
}

func ConnectNewFollower(ip string) {
	go func() {
		conn, err := net.Dial("tcp", ip)
		if err == nil {
			followersLock.Lock()
			newFollower := &Follower{Conn: conn}
			followerList = append(followerList, newFollower)
			followersLock.Unlock()
			handleConn(newFollower)
		}
	}()
}

func handleConn(follower *Follower) {

	defer follower.Conn.Close()
	defer onConnClose(follower)

	for {
		inMsg := &Message{}
		data := make([]byte, 4096)
		n, err := follower.Conn.Read(data)
		if err != nil {
			return
		}
		err = proto.Unmarshal(data[0:n], inMsg)
		if err != nil {
			return
		}
		SendMessage(*inMsg.Text, *inMsg.ChatId)
	}
}

func onConnClose(follower *Follower) {
}

func GotNewMessage(msg *tgbotapi.Message) {
	_, ok := chatsToFollower[msg.Chat.ID]
	if !ok {
		connectNewChat(msg.Chat.ID)
	}
	sendNewMessageToAI(msg.Text, int64(msg.From.ID), msg.Chat.ID)
}

func sendNewMessageToAI(text string, userid int64, chatid int64) {
	msg := &Message{
		Text:   proto.String(text),
		UserId: proto.Int64(userid),
		ChatId: proto.Int64(chatid),
	}
	marshaledMsg, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	follower := chatsToFollower[chatid]
	follower.Conn.Write(marshaledMsg)
}
