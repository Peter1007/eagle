# eagle
基本rabbitmq的rpc全框架
分为4个部分
cgi：php，基于yaf，客户端请求入口；调用一个或多个business接口
business：golang，接受cgi请求，调用一个或多个service接口
service：golang，接受business请求，直接和数据交互
statistic：golang，数据统计模块，包括全链路请求日志和其他的业务日志