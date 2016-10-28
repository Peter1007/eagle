<?php
class CityModel extends BaseModel
{
    private $_business_queue = 'user';
    private $_source = 'web';
    
    public function getList()
    {
        $route = 'city/get_list';
        $data = [];
        
        return $this->callRpc($this->_source, $route, $data, $this->_business_queue);
    }
    
    public function getInfo($cityId)
    {
        $route = 'city/get_info';
        $data = ['city_id' => $cityId];
        
        return $this->callRpc($this->_source, $route, $data, $this->_business_queue);
    }
}
?>
