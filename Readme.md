# Seemer Messaging Server

Instant messaging server. Backend in pure [Go](http://golang.org) (license [GPL 3.0](http://www.gnu.org/licenses/gpl-3.0.en.html)), client-side binding in Flutter as well as [gRPC](https://grpc.io/) client support for C++, C#, Go, Java, Node, PHP, Python, Ruby, Objective-C, etc. Persistent storage is any one of [YugaByteDB](https://yugabyte.com/). 

### Supported

* Multiple platforms:
  * [Android](https://github.com/tinode/tindroid/)
  * [iOS](https://github.com/tinode/ios)
  * [Web](https://github.com/tinode/webapp/)
* One-on-one and group messaging.
* Channels with an unlimited number of read-only subscribers.
* Server-generated presence notifications for people, group chats.
* User search/discovery.
* Inline images, file attachments.
* Message status notifications: message delivery to server; received and read notifications; typing notifications.
* Most recent message preview in contact list.
* Ability to block unwanted communication server-side.
* Anonymous users (important for use cases related to tech support over chat).
* Storage and out of band transfer of large objects like video files using local file system or Amazon S3.
* End to end encryption with [OTR](https://en.wikipedia.org/wiki/Off-the-Record_Messaging) for one-on-one messaging and undecided method for group messaging.
* Replying and forwarding messages.
* Voice and video messages, location sharing.
* Previews of attached videos, documents, links.