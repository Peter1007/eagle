<?php
class Amq
{
    private static $_instance;
    private static $_connections;
    private static $_mqConn;
    private static $_response = [];
    private static $_corrId = '';
    private static $_responseList = [];
    private static $_corrIdList = [];
    
    private static $_vhost_business = 'business';
    private static $_vhost_statistic = 'statistic';
    
    public static function getInstance()
    {
        if (!(self::$_instance instanceof self)) {
            self::$_instance = new self;
            
            $rabbitmqConfig = Yaf_Application::app()->getConfig()->rabbitmq;
            
            foreach ($rabbitmqConfig as $key => $config) {
                self::$_connections[$key] = new PhpAmqpLib\Connection\AMQPStreamConnection($config->host, $config->port, $config->user, $config->password, $config->vhost);
                self::$_mqConn[$key]['channel'] = self::$_connections[$key]->channel();
                
                if ($key == 'statistic') {
                    self::$_mqConn[$key]['items'] = $config->items;
                } else {
                    self::$_mqConn[$key]['queue'] = $config->queue;
                }
            }
        }
        
        return self::$_instance;
    }
    
    public function callRpc($queue, $body)
    {
        $channel = self::$_mqConn[self::$_vhost_business]['channel'];
        
        list($queueName, ,) = $channel->queue_declare('', false, false, true, false);
        
		self::$_corrId = self::genCorrId();
        
		$data = new PhpAmqpLib\Message\AMQPMessage(json_encode($body), ['correlation_id' => self::$_corrId, 'reply_to' => $queueName]);
		$channel->basic_publish($data, '', $queue);
        
        $channel->basic_consume($queueName, '', false, false, false, false, [$this, 'callRpcResponse']);
        
        self::$_response = [];
		while(empty(self::$_response)) {
            try {
                //超时5秒
                $channel->wait(NULL, FALSE, 5);
            } catch (Exception $e) {
                self::$_response = [
                    'code' => 1000,
                    'msg' => $e->getMessage(),
                    'data' => new stdClass()
                ];
            }
        }
        
        return self::$_response;
    }
    
    public function callRpcResponse($response)
    {
        if($response->get('correlation_id') == self::$_corrId) {
			self::$_response = json_decode($response->body, TRUE);
		}
    }


    public function callMultiRpc($mqParamList)
    {
        $channel = self::$_mqConn[self::$_vhost_business]['channel'];
        
        list($queueName, ,) = $channel->queue_declare('', false, false, true, false);
        
        self::$_corrIdList[] = [];
        foreach ($mqParamList as $key => $param) {
            $corrId = self::genCorrId();
            self::$_corrIdList[$corrId] = $key;
            
            $data = new PhpAmqpLib\Message\AMQPMessage(json_encode($param['body']), ['correlation_id' => $corrId, 'reply_to' => $queueName]);
            $channel->basic_publish($data, '', $param['queue']);

            $channel->basic_consume($queueName, '', false, false, false, false, [$this, 'callMultiRpcResponse']);
        }
		
        $callTimes = count($mqParamList);
        self::$_responseList = [];
        while(count(self::$_responseList) < $callTimes) {
            try {
                //超时5秒
                $channel->wait(NULL, FALSE, 5);
            } catch (Exception $e) {
                self::$_responseList[] = [
                    'code' => 1000,
                    'msg' => $e->getMessage(),
                    'data' => new stdClass()
                ];
            }
        }
        
        return self::$_responseList;
    }
    
    public static function callMultiRpcResponse($response)
    {
        $corrId = $response->get('correlation_id');
        if(isset(self::$_corrIdList[$corrId])) {
            self::$_responseList[self::$_corrIdList[$corrId]] = json_decode($response->body, TRUE);
		}
    }
    
    public function callReqLog($param)
    {
        $channel = self::$_mqConn[self::$_vhost_statistic]['channel'];
        $exchange = self::$_mqConn[self::$_vhost_statistic]['items']->reqlog->exchange;
        $key = self::$_mqConn[self::$_vhost_statistic]['items']->reqlog->key;
        
        $msg = new PhpAmqpLib\Message\AMQPMessage(json_encode($param, JSON_UNESCAPED_UNICODE));
        $channel->basic_publish($msg, $exchange, $key);
    }
    
    private static function genCorrId()
    {
        return uniqid().mt_rand(0, 1000000);
    }


    public function close()
    {
        foreach (self::$_mqConn as $conn) {
            $conn['channel']->close();
        }
        
        foreach (self::$_connections as $connection) {
            $connection->close();
        }
    }
}
?>
