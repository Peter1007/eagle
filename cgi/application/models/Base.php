<?php

class BaseModel
{
    protected function callRpc($source, $route, $data, $queue)
    {
        $rpcParam = [
            'meta' => [
                'req_id' => REQ_ID,
                'step' => 0,
                'source' => $source,
                'route' => $route
            ],
            'data' => $data
            
        ];
        
        return Amq::getInstance()->callRpc($queue, $rpcParam);
    }
}
