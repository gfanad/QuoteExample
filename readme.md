1. 首先创建客户端设置，主要有两个参数，服务器地址及用户id（推送信息可以从由推送端主动从用户中心拉取，客户端也可以主动订阅或者取消订阅）。
![img.png](imgs/img.png)
2. 首先需要获取读取channel读取推送信息，信息为以json编码。
![img_2.png](imgs/img_2.png)
3. 在获取读取channel后，用户可以开始主动调用sub、ubsub来订阅和取消订阅。
![img_3.png](imgs/img_3.png)
![img_4.png](imgs/img_4.png)
