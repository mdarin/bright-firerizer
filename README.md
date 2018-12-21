# bright-firerizer
The Firebase Admin Go SDK microservice with libmqtt interface

its just a draft

Admin SDK
https://github.com/firebase/firebase-admin-go

MQTT
https://github.com/goiiot/libmqtt

```bash
#..start the emqttd broker before
./build.sh
./run.sh
```
you should see 
```
$ ./run.sh 
prepareing evn...
starting the server...
2018/12/21 10:15:43 connect
2018/12/21 10:15:43 subscribe
2018/12/21 10:15:43 subscribe to topic [foo] success: <nil>
2018/12/21 10:15:43 subscribe to topic [bar] success: <nil>
2018/12/21 10:15:43 publish
2018/12/21 10:15:43 publish packet to topic [foo] success: <nil>
2018/12/21 10:15:43 publish packet to topic [bar] success: <nil>
2018/12/21 10:15:43 [foo] message: bar
2018/12/21 10:15:43 [bar] message: foo
 [M] Channel -> bar
 [M] Channel -> foo
error: http error status: 400; reason: request contains an invalid argument; code: invalid-argument; details: The registration token is not a valid FCM registration token
Dry run successful: 
error: http error status: 400; reason: request contains an invalid argument; code: invalid-argument; details: The registration token is not a valid FCM registration token
Successfully sent message: 
 (I) Workers done: 1
 (I) Workers done: 2
Successfully sent message: projects/cargo-b2ec7/messages/6890471482712802953
 (I) Workers done: 3
Successfully sent message: projects/cargo-b2ec7/messages/5653827703332019154
 (I) Workers done: 4
 ! Done
 ! Timeout
```

re-build
```bash
./build.sh
```
clean all
```bash
./sweepup.sh
```
