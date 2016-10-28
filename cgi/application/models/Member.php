<?php
class MemberModel extends BaseModel
{
    private $_business_queue = 'user';
    private $_source = 'web';
    
    public function getListByCityId($cityId, $statusId)
    {
        $route = 'member/get_list_by_city_id';
        $data = [
            'city_id' => $cityId,
            'status_id' => $statusId
        ];
        
        return $this->callRpc($this->_source, $route, $data, $this->_business_queue);
    }
    
    public function getInfoById($memberId)
    {
        $route = 'member/get_info_by_id';
        $data = [
            'member_id' => $memberId
        ];
        
        return $this->callRpc($this->_source, $route, $data, $this->_business_queue);
    }
}
?>
