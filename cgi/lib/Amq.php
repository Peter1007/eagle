<?php
class Amq
{
    private static $_instance;
    private static $_mqConn;
    private static $_response;
    private static $_corr_id;
    
    public static function getInstance()
    {
        if (!(self::$_instance instanceof self)) {
            self::$_instance = new self;
            
            $amqConfig = Yaf_Application::app()->getConfig()->amq;
            self::$_mqConn['queue'] = $amqConfig->queue;

            $connection = new PhpAmqpLib\Connection\AMQPStreamConnection($amqConfig->host, $amqConfig->port, $amqConfig->user, $amqConfig->password, $amqConfig->vhost);
            self::$_mqConn['channel'] = $connection->channel();
        }
        
        return self::$_instance;
    }
    
    public function call()
    {
        
    }
    
    public function callRpc($queue, $param)
    {
        list($queueName, ,) = self::$_mqConn['channel']->queue_declare('', false, false, true, false);
		self::$_corr_id = uniqid();
        
		$data = new PhpAmqpLib\Message\AMQPMessage(json_encode($param), ['correlation_id' => self::$_corr_id, 'reply_to' => $queueName]);
		self::$_mqConn['channel']->basic_publish($data, '', self::$_mqConn['queue'][$queue]);
        
        self::$_mqConn['channel']->basic_consume($queueName, '', false, false, false, false, [$this, 'callRpcResponse']);
        
		while(!self::$_response) {
			self::$_mqConn['channel']->wait();
		}
        
		return json_decode(self::$_response, TRUE);
    }
    
    public function callRpcResponse($response)
    {
        if($response->get('correlation_id') == self::$_corr_id) {
			self::$_response = $response->body;
		}
    }


    public function callMultiRpc()
    {
        
    }
    
    public static function callMultiRpcResponse()
    {
        
    }
}
?>
