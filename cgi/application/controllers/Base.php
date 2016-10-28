<?php

class BaseController extends Yaf_Controller_Abstract
{
    protected function response($ret)
    {
        $resp = '';
        
        if ($ret['code'] == CODE_SUCCESS) {
            $resp = $this->responseOk($ret['data']);
        } else {
            $msg = !empty($ret['msg']) ? $ret['msg'] : '';
            $resp = $this->responseFailed($ret['code'], $msg);
        }
        
        echo $resp;
    }

    private function responseOk($data)
    {
        $response = [
            'code' => CODE_SUCCESS,
            'msg' => MSG_SUCCESS,
            'body' => $data
        ];

        return json_encode($response,JSON_UNESCAPED_UNICODE);
    }
    
    private function responseFailed($code, $msg)
    {
        $response = [
            'code' => (int)$code,
            'msg' => $msg,
            'body' => []
        ];
        
        return json_encode($response,JSON_UNESCAPED_UNICODE);
    }
}
