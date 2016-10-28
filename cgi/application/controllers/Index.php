<?php

class IndexController extends BaseController
{
   public function indexAction()
    {
       $connection = new PhpAmqpLib\Connection\AMQPStreamConnection('192.168.1.129', 5672, 'test_user', '123456', 'test');
        $channel = $connection->channel();
        
        $channel->queue_declare('task_queue', false, true, false, false);
        
        $data = [
            'meta' => [
                'req_id' => REQ_ID,
                'step' => 0
            ],
        ];
        $msg = new PhpAmqpLib\Message\AMQPMessage($data,
                                array('delivery_mode' => 2) # make message persistent
                              );
        $channel->basic_publish($msg, '', 'task_queue');

        echo " [x] Sent ", $data, "\n";

        $channel->close();
        $connection->close();

       $ret = [
           'code' => 0,
           'msg' => 'success',
           'data' => []
       ];
       $this->response($ret);
    }
}