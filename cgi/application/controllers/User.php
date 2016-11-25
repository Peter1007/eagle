<?php

class UserController extends BaseController
{
    private $userModel;
    
    public function init()
    {
        parent::init();
        
        $this->userModel = new UserModel();
    }
    
    public function __destruct() {
        parent::__destruct();
    }

    public function indexAction()
    {
        $statusId = isset($this->inputData['status_id']) ? (int)$this->inputData['status_id'] : 0;
        
        $result = $this->userModel->index($statusId);

        $this->response($result, 'User/index');
    }
    
    public function getCityListAction()
    {
        $statusId = isset($this->inputData['status_id']) ? (int)$this->inputData['status_id'] : 0;
        
        $result = $this->userModel->getCityList($statusId);

        $this->response($result, 'User/getCityList');
    }
    
    public function getMemberListAction()
    {
        $statusId = isset($this->inputData['status_id']) ? (int)$this->inputData['status_id'] : 0;
        
        $result = $this->userModel->getMemberList($statusId);

        $this->response($result, 'User/getMemberList');
    }
    
    public function getMemberTypeListAction()
    {
        $statusId = isset($this->inputData['status_id']) ? (int)$this->inputData['status_id'] : 0;
        
        $result = $this->userModel->getMemberTypeList($statusId);

        $this->response($result, 'User/getMemberTypeList');
    }
}