<?php
class UserModel extends BaseModel
{
    public function __construct()
    {
        parent::__construct();
    }
    
    public function index($statusId)
    {
        $paramList = [
            [
                'route' => 'city/get_list',
                'data' => [
                    'status_id' => $statusId
                ],
                'queue' => $this->businessQueue->base
            ],
            [
                'route' => 'member/get_list',
                'data' => [
                    'status_id' => $statusId
                ],
                'queue' => $this->businessQueue->user
            ],
            [
                'route' => 'member_type/get_list',
                'data' => [
                    'status_id' => $statusId
                ],
                'queue' => $this->businessQueue->base
            ]
        ];
        
        $responseList = $this->callMultiRpc($paramList);
        foreach ($responseList as $response) {
            if (!isset($response['code']) || $response['code'] > 0) {
                return $response;
            }
        }
        
        return [
            'code' => 0,
            'msg' => 'success',
            'data' => [
                'city_list' => $responseList[0]['data'],
                'member_list' => $responseList[1]['data'],
                'member_type_list' => $responseList[2]['data']
            ]
        ];
    }
    
    public function getCityList($statusId)
    {
        $route = 'city/get_list';
        $data = [
            'status_id' => $statusId
        ];
        
        return $this->callRpc($route, $data, $this->businessQueue->base);
    }
    
    public function getMemberList($statusId)
    {
        $route = 'member/get_list';
        $data = [
            'status_id' => $statusId
        ];
        
        return $this->callRpc($route, $data, $this->businessQueue->user);
    }
    
    public function getMemberTypeList($statusId)
    {
        $route = 'member_type/get_list';
        $data = [
            'status_id' => $statusId
        ];
        
        return $this->callRpc($route, $data, $this->businessQueue->base);
    }
}
?>
