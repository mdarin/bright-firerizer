# Cloud Messaging Backend

## Abstract
FCM(Firebase Cloud Messaging) основан на модели публикации / подписки, что позволяет реализовать отправку сообщений в каналы, которые, в свою очередь, позволяют гибко группировать подписчиков по темам на которые они подписаны, за счёт использования условной маршрутизации.
Предлагается реализовать службу push-уведомлений для сервиса cago на языке go с использованием Firebase Admin SDK работающую через брокер MQTT и  дубль через REST API на всякий случай. 

список библиотек mqtt для go
https://golanglibs.com/category/mqtt?sort=top

библиотека для работы по mqtt для go её для IoT пилят и выглядит она годной
https://github.com/goiiot/libmqtt

SDK для реализации бакенда
https://firebase.google.com/docs/admin/setup?authuser=0

отправка сообщений администрирование (бакенд) нотификаций
https://firebase.google.com/docs/cloud-messaging/admin/send-messages?authuser=0

описание библиотеки SDK для go
https://godoc.org/firebase.google.com/go/messaging

есть и REST встроенный сервер с POST запросами
https://firebase.google.com/docs/cloud-messaging/send-message?authuser=0

