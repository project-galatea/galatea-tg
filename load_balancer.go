package main

import (
	"errors"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

	"gopkg.in/telegram-bot-api.v4"
)

type Chat struct {
	chatId int64
}

type Follower struct {
	Chats []Chat
	Conn  net.Conn
}

var followerList []*Follower

var chatsToFollower = make(map[int64]*Follower)
var followersLock = &sync.Mutex{}

// TODO: Make this more of an algorithm and less of a guess
func connectNewChat(chatid int64) error {

	err, fl := selectBestFollower()
	if err != nil {
		return err
	}

	followersLock.Lock()

	chatsToFollower[chatid] = fl

	fl.Chats = append(fl.Chats, Chat{chatId: chatid})

	followersLock.Unlock()

	return nil
}

func selectBestFollower() (error, *Follower) {
	followersLock.Lock()
	defer followersLock.Unlock()

	if len(followerList) == 0 {
		return errors.New("No followers to connect to"), nil
	}

	rand.Seed(time.Now().Unix())
	setFollower := followerList[rand.Intn(len(followerList))]
	return nil, setFollower
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
	for i, fl := range followerList {
		if fl.Conn.LocalAddr() == follower.Conn.LocalAddr() {
			followerList = append(followerList[:i], followerList[i+1:]...)
		}
	}
	for _, chat := range follower.Chats {
		err := connectNewChat(chat.chatId)
		if err != nil {
			return
		}
	}
}

func GotNewMessage(msg *tgbotapi.Message) error {
	_, ok := chatsToFollower[msg.Chat.ID]
	if !ok {
		err := connectNewChat(msg.Chat.ID)
		if err != nil {
			return err
		}
	}

	err := sendNewMessageToAI(msg.Text, int64(msg.From.ID), msg.Chat.ID)

	return err
}

func sendNewMessageToAI(text string, userid int64, chatid int64) error {
	msg := &Message{
		Text:   proto.String(text),
		UserId: proto.Int64(userid),
		ChatId: proto.Int64(chatid),
	}
	marshaledMsg, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	follower := chatsToFollower[chatid]
	follower.Conn.Write(marshaledMsg)

	return nil
}
