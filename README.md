# Pub/Sub System

### Overview
This project is a simple implementation of a Publish/Subscribe system, also known as Pub/Sub.

### Features  
**Publisher-Subscriber Model:** Implements the classic Pub/Sub pattern where publishers and subscribers are decoupled.  
**Topic-Based Messaging:** Messages are categorized into topics, and subscribers can subscribe to specific topics of interest.  
**Flexible Subscription:** Subscribers can dynamically subscribe and unsubscribe to topics at runtime. Topics will be auto created on-to-fly if not present.  
**Asynchronous Communication:** Supports asynchronous communication between publishers and subscribers.  
**Multiple API:** Supports gRPC and Websocket for real-time communication between publishers and subscribers. 

   
### Installation
You will need Go 1.22 or higher to install this.  
Run `git clone https://github.com/sgkandale/go_pub_sub` in your directory.  
`cd go_pub_sub` to change the directory.  
`go build ./cmd/pubsub`  to build the executable binary.  
Sample config file is included in repo. Customize it as per your needs.  
You can specify the location of config file using `--config` flag, if it is located in different directory.

### Contributing
Contributions are welcome! If you have any ideas, enhancements, or bug fixes, please feel free to open an issue or submit a pull request.

### License
This project is licensed under the MIT License - see the LICENSE file for details.