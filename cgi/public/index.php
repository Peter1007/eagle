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
define('REQ_NO', md5(uniqid().mt_rand(0, 1000000)));
define('START_TIME', time());
define('MICRO_TIME', getMicrotime());
define('DEFAULT_SOURCE', 'web');
define('MODULE', 'cgi_web');
define('REQLOG_ROUTE', 'reqlog/add');

ini_set('yaf.library', APP_PATH.'/lib');

header("Content-type:text/json; charset=utf-8");

Yaf_Dispatcher::getInstance()->disableView();

//xhprof_enable();

$app  = new Yaf_Application(APP_PATH.'/conf/application.ini', ENV);
$app->run();

function getMicrotime()
{
    $mtime = explode(" ", microtime());
    return $mtime[0];
}

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