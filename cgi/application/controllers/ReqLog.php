<?php

class ReqLogController extends BaseController
{
    public function indexAction()
    {
        $connection = new PhpAmqpLib\Connection\AMQPStreamConnection('192.168.1.129', 5672, 'statisticuser', 'statisticpasswd', 'statistic');
        $channel = $connection->channel();

        $reqData = json_encode(['status' => 1]);
        $respData = json_encode(new stdClass());
        $param = [
            'meta' => new stdClass(),
            'route' => 'reqlog/add',
            'data' => [
                'req_no' => REQ_NO,
                'step' => (string)0,
                'source' => 'web',
                'req_time' => (string)time(),
                'module' => 'web_cgi',
                'route' => 'City/getList',
                'req_data' => $reqData,
                'req_size' => (string)strlen($reqData),
                'cost_time' => (string)2000000,
                'code' => (string)0,
                'msg' => 'success',
                'resp_data' => $respData,
                'resp_size' => (string)strlen($respData)
            ]
        ];
        $msg = new PhpAmqpLib\Message\AMQPMessage(json_encode($param));

        $channel->basic_publish($msg, 'statistic_reqlog', 'cgi_web');

        echo "[x] Sent reqlog: ",json_encode($param)," \n";

        $channel->close();
        $connection->close();
    }
}