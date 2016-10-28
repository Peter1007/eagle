<?php
header('Content-type: text/html; charset=utf-8');

//设置中国时区
date_default_timezone_set('PRC');

//指向public的上一级
define("APP_PATH",  realpath(dirname(__FILE__) . '/../'));
define("SERVICE_PATH",  APP_PATH.'/lib/service');
define('CODE_SUCCESS', 0);
define('CODE_ERROR', -1);
define('MSG_SUCCESS', 'success');
define('ENV', 'dev');
define('REQ_ID', md5(uniqid().mt_rand(0, 1000000)));

ini_set('yaf.library', APP_PATH.'/lib');

header("Content-type:text/json; charset=utf-8");

Yaf_Dispatcher::getInstance()->disableView();

//xhprof_enable();

$app  = new Yaf_Application(APP_PATH.'/conf/application.ini', ENV);
$app->run();

//$xhprof_data = xhprof_disable();
//$XHPROF_ROOT = "D:\xhj\xhprof";
//include_once $XHPROF_ROOT . "/xhprof_lib/utils/xhprof_lib.php";  
//include_once $XHPROF_ROOT . "/xhprof_lib/utils/xhprof_runs.php";  
//
//$xhprof_runs = new XHProfRuns_Default();
//$source = 'eagle_cgi';
//$run_id = $xhprof_runs->save_run($xhprof_data, $source);
//
//$href = 'http://xhprof.panwang.com/xhprof_html/index.php?run='.$run_id.'&source='.$source;
//echo "<a href='$href'>profile</a>";