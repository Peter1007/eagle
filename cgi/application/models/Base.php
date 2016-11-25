<?php

class BaseModel
{
    private static $setp = 0;
    protected $businessQueue;
    
    public function __construct()
    {
        $this->businessQueue = Yaf_Application::app()->getConfig()->rabbitmq->business->queue;
    }

    protected function callRpc($route, $data, $queue)
    {
        $rpcParam = [
            'meta' => [
                'req_no' => REQ_NO,
                'step' => $this->nowStep(),
                'source' => SOURCE,
                'remote_ip' => ip2long($_SERVER['REMOTE_ADDR']),
                'req_time' => START_TIME
            ],
            'route' => $route,
            'data' => $data
            
        ];
        
        return Amq::getInstance()->callRpc($queue, $rpcParam);
    }
    
    protected function callMultiRpc($paramList)
    {
        $mqParamList = [];
        
        foreach ($paramList as $param) {
            $mqParamList[] = [
                'queue' => $param['queue'],
                'body' => [
                    'meta' => [
                        'req_no' => REQ_NO,
                        'step' => $this->nowStep(),
                        'source' => SOURCE,
                        'remote_ip' => ip2long($_SERVER['REMOTE_ADDR']),
                        'req_time' => START_TIME
                    ],
                    'route' => $param['route'],
                    'data' => $param['data']
                ]
            ];
        }
        
        if (!empty($mqParamList)) {
            return Amq::getInstance()->callMultiRpc($mqParamList);
        } else {
            return [];
        }
    }
    
    private function nowStep()
    {
        if (self::$setp == 0) {
            self::$setp = 1000;
        } else {
            self::$setp += 100;
        }
        
        return self::$setp;
    }
}
