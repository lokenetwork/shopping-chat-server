# 买买提客服系统

客服IM系统，采用GO+websocket+protobuf技术实现。第一版已经开发完成，在单机16核32G上压测，可支持40万并发连接，每秒处理1万条消息。第二版打算采用rabbitmq做历史消息写入的中间层，同时做多机器版本，可以支持更多的并发连接，因为没有时间，我暂时没去开发第二版，第二版的具体架构面试详谈。

 ![image](https://github.com/lokenetwork/hybird_shopping_app/blob/master/demo-picture/consult-1.png)
 ![image](https://github.com/lokenetwork/hybird_shopping_app/blob/master/demo-picture/consult-2.png)




# 买买提平台关项目代码

客户端APP：https://github.com/lokenetwork/hybird_shopping_app

PHP接口服务器：https://github.com/lokenetwork/shopping-php-server

GO聊天服务器：https://github.com/lokenetwork/shopping-chat-server

通行证服务器（单点登录）：https://github.com/lokenetwork/passport-shopping-system

商家后台：https://github.com/lokenetwork/shop-manage

总平台后台：https://github.com/lokenetwork/admin-shopping


