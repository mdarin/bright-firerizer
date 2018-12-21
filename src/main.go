// --- Initialize the SDK ---
package main

import (
// common purupose
	"fmt"
	"time"
// admin sdk
	"context"
	"log"

	firebase "firebase.google.com/go"
	messaging "firebase.google.com/go/messaging"
//	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
// mqtt
	"github.com/goiiot/libmqtt"
//	"examples"
//TODO:rest
)

const(
	TIMEOUT = 3000
)

var(
	msgCounter int = 0
)
//
// main driver
//
func main() {

	//
	// MQTTLIB example
	//
	// If you would like to explore all the options available, please refer to
	// [https://godoc.org/github.com/goiiot/libmqtt#Option]
	res := make(chan []byte)
	defer close(res)
	helloMQTT(res)

	//
	// Admin SDK example
	//
	opt := option.WithCredentialsFile("/home/user/src/gpush/credentials/cargo-b2ec7-firebase-adminsdk-xv6k8-dc31d82f3e.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	done := make(chan bool)
	defer close(done)
	// --- Send to individual devices ---
	go func() {
		defer func() { done<- true }()
		defer func() { recover() }()
		helloAdminSDK_Send_to_individual_devices(app)
	}()

	// --- Send to a topic ---
	go func() {
		defer func() { done<- true }()
		defer func() { recover() }()
		helloAdminSDK_Send_to_a_topic(app)
	}()

	// --- Send to a condition ---
	go func() {
		defer func() { done<- true }()
		defer func() { recover() }()
		helloAdminSDK_Send_to_a_condition(app)
	}()

	// --- Sending in the dry run mode ---
	go func() {
		defer func() { done<- true }()
		defer func() { recover() }()
		helloAdminSDK_Sending_in_the_dry_run_mode(app)
	}()


	doneCounter := 0;
	for stop := false; stop != true; {
		select {
			case msg := <-res:
				fmt.Println(" [M] Channel ->", string(msg))
			case <-done:
				doneCounter++
				if doneCounter >= 4 {
					fmt.Println(" (I) Workers done:", doneCounter)
					fmt.Println(" ! Done")
				} else {
					fmt.Println(" (I) Workers done:", doneCounter)
				}
			//case <-time.After(TIMEOUT * time.Microsecond):
			case <-time.After(TIMEOUT * time.Millisecond):
				fmt.Println(" ! Timeout")
				stop = true
		} // eof select
	} // eof for
} // eof main


func helloMQTT(res chan []byte) {

	// 1.Go get this project
	// 2.Import this package in your project file
	// 3.Create a custom client
	// Create a client and enable auto reconnect when connection lost
	// We primarily use `RegexRouter` for client
	client, err := libmqtt.NewClient(
		//try MQTT 5.0 and fallback to MQTT 3.1.1
		//libmqtt.WithVersion(libmqtt.V5, true),
		// server address(es)
		libmqtt.WithServer("localhost:1883"),
		libmqtt.WithIdentity("user", "public"),
		//libmqtt.WithServer("cargo.tvzavr.ru:1883"),
		// authorize 
		//> t_auth:login(<<"+79615244722">>, <<"123123">>).
		// mqtt_login:Login  mqtt_password: Password
		//libmqtt.WithIdentity("+79615244722", "7292a0b4-0457-11e9-b766-f44d309c2889"),
		// enable keepalive (10s interval) with 20% tolerance
		libmqtt.WithKeepalive(10, 1.2),
		// enable auto reconnect and set backoff strategy
		libmqtt.WithAutoReconnect(true),
		libmqtt.WithBackoffStrategy(time.Second, 5*time.Second, 1.2),
		// use RegexRouter for topic routing if not specified
		// will use TextRouter, which will match full text
		libmqtt.WithRouter(libmqtt.NewRegexRouter()),
	)

	if err != nil {
		// handle client creation error
		panic("create mqtt client failed")
	}

	// register net handler
	{
		client.HandleNet(func(server string, err error) {
			if err != nil {
				log.Printf("error happened to connection to server [%v]: %v", server, err)
			}
		})
		// register persist handler, you don't need this if all your message had QoS 0
		client.HandlePersist(func(err error) {
			if err != nil {
				log.Printf("session persist error: %v", err)
			}
		})
		// register subscribe handler
		client.HandleSub(func(topics []*libmqtt.Topic, err error) {
			if err != nil {
				for _, t := range topics {
					log.Printf("subscribe to topic [%v] failed: %v", t.Name, err)
				}
			} else {
				for _, t := range topics {
					log.Printf("subscribe to topic [%v] success: %v", t.Name, err)
				}

				// run publishing messages
				//go func() {
					log.Printf("publish")
					// publish some packet (just for example)
					client.Publish([]*libmqtt.PublishPacket{
						{TopicName: "foo", Payload: []byte("bar"), Qos: libmqtt.Qos2},//Qos0
						{TopicName: "bar", Payload: []byte("foo"), Qos: libmqtt.Qos2},//Qos1
					}...)
				//}()
			}
		})
		// register unsubscribe handler
		client.HandleUnSub(func(topic []string, err error) {
			if err != nil {
				// handle unsubscribe failure
				for _, t := range topic {
					log.Printf("unsubscribe to topic [%v] failed: %v", t, err)
				}
			} else {
				for _, t := range topic {
					log.Printf("unsubscribe to topic [%v] failed: %v", t, err)
				}
			}
		})
		// register publish handler
		client.HandlePub(func(topic string, err error) {
			if err != nil {
				log.Printf("publish packet to topic [%v] failed: %v", topic, err)
			} else {
				log.Printf("publish packet to topic [%v] success: %v", topic, err)
			}
		})

		// handle every subscribed message (just for example)
		client.Handle(".*", func(topic string, qos libmqtt.QosLevel, msg []byte) {
			log.Printf("[%v] message: %v", topic, string(msg))
			go func() {
				res<- msg
			}()
			msgCounter++
			if msgCounter >= 2 {
				//5.Unsubscribe topic(s)
				client.UnSubscribe("foo", "bar")
				//6.Destroy the client when you would like to
				// use true for a immediate disconnect to server
				// use false to send a DisConn packet to server before disconnect
				client.Destroy(true)
			}
		})

	} // eof handlers

	//4.Register the handlers and Connect, then you are ready to pub/sub with server
	// connect to server
	client.Connect(func(server string, code byte, err error) {
		log.Printf("connect")
		if err != nil {
			// failed
			panic(err)
		}

		if code != libmqtt.CodeSuccess {
			// server rejected or in error
			panic(code)
		}

		// success
		// you are now connected to the `server`
		// (the `server` is one of your provided `servers` when create the client)
		// start your business logic here or send a signal to your logic to start

		// on connected success
		// subscribe some topic(s)
		go func() {
			log.Printf("subscribe")
			client.Subscribe([]*libmqtt.Topic{
				{Name: "foo", Qos: libmqtt.Qos2},
				{Name: "bar", Qos: libmqtt.Qos2},
			}...)
		}()


	})

	// wait for all connection to exit 
	client.Wait()


	//5.Unsubscribe topic(s)
	//client.UnSubscribe("foo", "bar")

	//6.Destroy the client when you would like to
	// use true for a immediate disconnect to server
	// use false to send a DisConn packet to server before disconnect
	//client.Destroy(true)
}


func helloAdminSDK_Send_to_individual_devices(app *firebase.App) {
	// --- Send to individual devices ---

	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.
	registrationToken := "YOUR_REGISTRATION_TOKEN"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Token: registrationToken,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		//log.Fatalln(err)
		fmt.Println("error:", err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
}


func helloAdminSDK_Send_to_a_topic(app *firebase.App) {
	// --- Send to a topic ---

	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.
	//registrationToken := "YOUR_REGISTRATION_TOKEN"

	// The topic name can be optionally prefixed with "/topics/".
	topic := "highScores"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Topic: topic,
	}

	// Send a message to the devices subscribed to the provided topic.
	response, err := client.Send(ctx, message)
	if err != nil {
		//log.Fatalln(err)
		fmt.Println("error:", err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)

}

func helloAdminSDK_Send_to_a_condition(app *firebase.App) {
	// --- Send to a condition ---
	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.
	//registrationToken := "YOUR_REGISTRATION_TOKEN"

	// Define a condition which will send to devices which are subscribed
	// to either the Google stock or the tech industry topics.
	condition := "'stock-GOOG' in topics || 'industry-tech' in topics"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Condition: condition,
	}

	// Send a message to devices subscribed to the combination of topics
	// specified by the provided condition.
	response, err := client.Send(ctx, message)
	if err != nil {
		//log.Fatalln(err)
		fmt.Println("error:", err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)

}

func helloAdminSDK_Sending_in_the_dry_run_mode(app *firebase.App) {

	// --- Sending in the dry run mode ---

	// Firebase Admin SDK supports sending FCM messages in the dry run mode. The SDK and the 
	// FCM service perform all the usual validations on the messages sent in this mode, 
	// but they are not actually delivered to the target devices. Therefore this feature 
	// can be used to check if a certain message will be accepted by the SDK and the FCM
	// service for sending.
	// Send a message in the dry run mode.

	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.
	registrationToken := "YOUR_REGISTRATION_TOKEN"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Token: registrationToken,
	}

	response, err := client.SendDryRun(ctx, message)
	if err != nil {
		//log.Fatalln(err)
		fmt.Println("error:", err)
	}
	// Response is a message ID string.
	fmt.Println("Dry run successful:", response)

}

