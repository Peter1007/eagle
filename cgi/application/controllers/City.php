<?php

class CityController extends BaseController
{
    public function getListAction()
    {
        $cityModel = new CityModel();
        $result = $cityModel->getList();

        $this->response($result);
    }
    
    public function getInfoAction()
    {
        $cityId = $_GET['city_id'];
        $cityModel = new CityModel();
        $result = $cityModel->getInfo($cityId);

        $this->response($result);
    }
}