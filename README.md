# 即时通讯IMxit
## 项目描述
这个项目这项目主要是在终端环境下运行的，可以完成日常的在线聊天，私聊，公聊，查询用户，超时踢人和即时通讯。服务器能够支持相对数量的客户端并发并及时响应。
## 主要工作
1. 利用Socket实现不同主机之间的通信
2. 使用Map记录在线的用户信息，并把用户上线信息通知给在线用户
3. 完成查询在线用户和用户之间的私聊功能
4. 利用Goruntine机制服务，增加并行服务数量
## 个人收获
该项目全面覆盖的go的基础语法和语法特性,加深对知识点的掌握，并对于goruntine,channel 有了进一步的认识。
## 使用方式【发送消息的格式】
1. 查询用户 who
2. 修改用户名 rename|张三
3. 私聊 to|张三|消息内容 
