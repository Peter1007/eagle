<?php

class MemberController extends BaseController
{
    public function getListByCityIdAction()
    {
        $cityId = $_GET['city_id'];
        $statusId = $_GET['status_id'];
        
        $memberModel = new MemberModel();
        $result = $memberModel->getListByCityId($cityId, $statusId);

        $this->response($result);
    }
    
    public function getInfoByIdAction()
    {
        $memberId = $_GET['member_id'];
        
        $memberModel = new MemberModel();
        $result = $memberModel->getInfoById($memberId);

        $this->response($result);
    }
}