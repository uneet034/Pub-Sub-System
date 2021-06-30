package pubsub

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

var (
	PUBLISH     = "publish"
	SUBSCRIBE   = "subscribe"
	UNSUBSCRIBE = "unsubscribe"

)

type PubSub struct {
	Users  []User
	Subscriptions []Subscription
}

type User struct {
	Id  string
	Connection *websocket.Conn
}

type Message struct {
	Action  string          `json:"action"`
	Topic   string          `json:"topic"`
	Message json.RawMessage `json:"message"`
}

type Subscription struct {
	Topic  string
	User *User
}

func (ps *PubSub)AddUser(user User) (*PubSub) {

	ps.Users = append(ps.Users, user)

	//fmt.Println("Adding new user to the list", user.Id, len(ps.Users))
	payload := []byte("Hello User ID" +
		user.Id)
	user.Connection.WriteMessage(1, payload)

	return ps

}

func (ps *PubSub) RemoveUser(user User) (*PubSub){

	// first remove all subscriptions by this user

	for index, sub := range ps.Subscriptions {

		if user.Id == sub.User.Id {
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}
	// remove user from the list

	for index, c := range ps.Users{

		if c.Id == user.Id {
			ps.Users = append(ps.Users[:index], ps.Users[index+1:]...)
		}

	}
	return ps
}

func (ps *PubSub) GetSubscriptions(topic string, user *User) ([]Subscription){

	var subscriptionList []Subscription

	for _, subscription := range ps.Subscriptions{

		if user != nil {

			if subscription.User.Id == user.Id && subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)

			}
		} else {

			if subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)
			}
		}


	}
	return subscriptionList
}

func (ps *PubSub) Subscribe(user *User, topic string) (*PubSub) {

	userSubs := ps.GetSubscriptions(topic, user)

	if len(userSubs) > 0 {

		// user is subscribed to this topic earlier

		return ps
	}

	newSubscription := Subscription{

		Topic: topic,
		User: user,
	}

	ps.Subscriptions = append(ps.Subscriptions, newSubscription)
	return ps
}


func (ps *PubSub) Publish(topic string, message []byte, excludeUser *User){

	subscriptions := ps.GetSubscriptions(topic, nil)

	for _, sub := range subscriptions {

		fmt.Printf("Sending to user id %s message is %s \n", sub.User.Id, message)
		//sub.User.Connection.WriteMessage(1, message)

		sub.User.Send(message)

	}


}
func (user *User) Send(message [] byte) (error) {

	return user.Connection.WriteMessage(1, message)

}


func (ps *PubSub) Unsubscribe(user *User, topic string) (*PubSub) {

	//userSubscriptions := ps.GetSubscriptions(topic, user)
	for index, sub := range ps.Subscriptions {

		if sub.User.Id == user.Id && sub.Topic == topic {
			// found this subscription from user and we do need to remove it
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}

	return ps

}


func (ps *PubSub) Ack(user User, messageType int, payload []byte) (*PubSub)  {


	m := Message{}

	err := json.Unmarshal(payload, &m)

	if err != nil {
		fmt.Println("This is not correct message payload")
		return ps
	}
	switch m.Action {

	case PUBLISH:

		fmt.Println("This is publish new message")

		ps.Publish(m.Topic, m.Message, nil)

		break

	case SUBSCRIBE:

		ps.Subscribe(&user, m.Topic)

		fmt.Println("New Subscriber to topic", m.Topic, len(ps.Subscriptions), user.Id)

		break

	case UNSUBSCRIBE:

		ps.Unsubscribe(&user, m.Topic)
		fmt.Println("User wants to unsubscribe the topic", m.Topic, user.Id)

		break



	default:
		break
	}

	return ps

}