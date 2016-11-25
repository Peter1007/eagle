<?php

class BaseController extends Yaf_Controller_Abstract
{
    protected $inputData;
    
    protected function init()
    {
        $this->inputData = $_GET;
        $source = empty($this->inputData['source']) ? DEFAULT_SOURCE : $this->inputData['source'];
        define('SOURCE', $source);
    }
    
    public function __destruct() {
        Amq::getInstance()->close();
    }

    protected function response($ret, $route)
    {
        $resp = '';
        
        if ($ret['code'] == CODE_SUCCESS) {
            $resp = $this->responseOk($ret['data']);
        } else {
            $msg = !empty($ret['msg']) ? $ret['msg'] : '';
            $resp = $this->responseFailed($ret['code'], $msg);
            $ret['data'] = new stdClass();
        }
        
        $this->reqLog($route, $ret);
        
        echo $resp;
    }

    private function responseOk($data)
    {
        $response = [
            'code' => CODE_SUCCESS,
            'msg' => MSG_SUCCESS,
            'data' => $data
        ];

        return json_encode($response,JSON_UNESCAPED_UNICODE);
    }
    
    private function responseFailed($code, $msg)
    {
        $response = [
            'code' => (int)$code,
            'msg' => $msg,
            'data' => new stdClass()
        ];
        
        return json_encode($response,JSON_UNESCAPED_UNICODE);
    }
    
    private function reqLog($route, $ret)
    {
        $reqData = json_encode($this->inputData);
        
        $resp = json_encode($ret['data'], JSON_UNESCAPED_UNICODE);
        
        $param = [
            'meta' => [
                'req_no' => REQ_NO,
                'step' => 0,
                'source' => SOURCE,
                'remote_ip' => ip2long($_SERVER['REMOTE_ADDR']),
                'req_time' => START_TIME
            ],
            'route' => REQLOG_ROUTE,
            'data' => [
                'module' => MODULE,
                'route' => $route,
                'req_data' => $reqData,
                'req_size' => strlen($reqData),
                'cost_time' => $this->getCostTime(),
                'code' => $ret['code'],
                'msg' => $ret['msg'],
                'resp_data' => $resp,
                'resp_size' => strlen($resp)
            ]
        ];
        
        return Amq::getInstance()->callReqLog($param);
    }
    
    private function getCostTime()
    {
        $costTime = (time()-START_TIME+getMicrotime()-MICRO_TIME)*1E+6;
        
        return (int)$costTime;
    }
}
