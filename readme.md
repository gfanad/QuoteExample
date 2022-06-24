1. 首先创建客户端设置，主要有两个参数，服务器地址及用户id（推送信息可以从由推送端主动从用户中心拉取，客户端也可以主动订阅或者取消订阅）。
![img.png](imgs/img.png)
2. 首先需要获取读取channel读取推送信息。
![img_1.png](imgs/img_1.png)
3. 用户可以主动调用sub、ubsub来订阅和取消订阅。
![img.png](imgs/img_3.png)