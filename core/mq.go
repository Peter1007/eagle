package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"reflect"
	"time"
)

const (
	FRONTEND              = "frontend"
	SERVICE               = "service"
	STATISTIC             = "statistic"
	ERR_FORMAT            = `{"code":%d,"msg":"%s","data":[]}`
	REQLOG_ROUTE          = "reqlog/add"
	ERRNO_ARGS_ERR        = 10000
	ERRNO_ROUTE_NOT_FOUND = 10001
	ERRNO_INNER           = 10002
	CORRID_LEN            = 32
	MQ_CONFIG_FILE        = "mq.yml"
	CONTENT_TYPE          = "text/plain"
	ERR_STEP              = 999
	ERR_REQNO             = "error req no"
)

type StatisticItem struct {
	ExchangeName string
	EXchangeType string
	Key          string
	Keys         []string
}

type MqConfig struct {
	Host   string
	Type   int
	Queue  string
	Queues map[string]string
	Items  map[string]StatisticItem
}

type ReqMeta struct {
	ReqNo    string `json:"req_no"`
	Step     int    `json:"step"`
	Source   string `json:"source"`
	RemoteIp int64  `json:"remote_ip"`
	ReqTime  int64  `json:"req_time"`
}

type ReqData struct {
	Meta  ReqMeta                `json:"meta"`
	Route string                 `json:"route"`
	Data  map[string]interface{} `json:"data"`
}

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type CallParam struct {
	Module string                 `json:"module"`
	Route  string                 `json:"route"`
	Data   map[string]interface{} `json:"data"`
}

type Mq struct {
	module string
}

var mqConfigs map[string]MqConfig
var conns map[string]*amqp.Connection
var channels map[string]*amqp.Channel
var mq *Mq
var reqMetaMap map[string]*ReqMeta
var isExit chan bool

func GetMqInstance() *Mq {
	if mq == nil {
		mq = &Mq{}

		reqMetaMap = map[string]*ReqMeta{}

		isExit = make(chan bool)
	}

	return mq
}

func (this *Mq) SetModule(module string) {
	this.module = module
}

func (this *Mq) initConfig() {
	mqConfigs = map[string]MqConfig{}
	SetConfig(MQ_CONFIG_FILE, &mqConfigs)
}

func (this *Mq) initConnsAndChannels() {
	var err error
	conns = map[string]*amqp.Connection{}
	channels = map[string]*amqp.Channel{}

	for key, mqConfig := range mqConfigs {
		conns[key], err = amqp.Dial(mqConfig.Host)
		if err != nil {
			logger.LogPanic("Failed to connect to RabbitMQ.", mqConfig.Host, err)
		}

		channels[key], err = conns[key].Channel()
		if err != nil {
			logger.LogPanic("Faile to open a channel.", err)
		}

		//设置RPC channel的qos
		if mqConfig.Type == 1 {
			err = channels[key].Qos(1, 0, false)
			if err != nil {
				logger.LogPanic("Failed to set Qos.", err)
			}
		}
	}
}

func (this *Mq) Run() {
	switch mqConfigs[FRONTEND].Type {
	case 1:
		this.runRpc(channels[FRONTEND])
	case 2:
		this.runStatistic(channels[FRONTEND])
	default:
		logger.LogWarm(fmt.Sprintf("Frontend type %d is error", mqConfigs[FRONTEND].Type))
	}
}

func (this *Mq) runRpc(ch *amqp.Channel) {
	queue, err := ch.QueueDeclare(mqConfigs[FRONTEND].Queue, false, false, false, false, nil)
	if err != nil {
		logger.LogPanic("Failed to declare a queue.", err)
	}

	reqs, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		logger.LogPanic("Failed to register a consumer.", err)
	}

	logger.LogInfo("[*] Awaiting RPC requests")

	for req := range reqs {
		go this.handleReq(req, 0)
	}
}

func (this *Mq) runStatistic(ch *amqp.Channel) {
	for _, item := range mqConfigs[FRONTEND].Items {
		go this.runStatisticItems(ch, &item)
	}

	<-isExit
}

func (this *Mq) runStatisticItems(ch *amqp.Channel, item *StatisticItem) {
	err := ch.ExchangeDeclare(item.ExchangeName, item.EXchangeType, true, false, false, false, nil)
	if err != nil {
		logger.LogPanic("Failed to declear an exchange.", err)
	}

	queue, err := channels[FRONTEND].QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		logger.LogPanic("Failed to declear a queue.", err)
	}

	for _, key := range item.Keys {
		logger.LogInfo(fmt.Sprintf("Binding queue %s to exchange %s with routing key %s", queue.Name, item.ExchangeName, key))
		err = ch.QueueBind(queue.Name, key, item.ExchangeName, false, nil)
		if err != nil {
			logger.LogPanic("Failed to bind a queue.", err)
		}
	}

	reqs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		logger.LogPanic("Failed to register a consumer.", err)
	}

	logger.LogInfo("[*] Starting", item.ExchangeName)

	for req := range reqs {
		go this.handleReq(req, 1)
	}
}

func (this *Mq) handleReq(req amqp.Delivery, reqType int) {
	fmt.Println(string(req.Body))
	startTimeNano := time.Now().UnixNano()

	reqData := ReqData{}
	err := json.Unmarshal(req.Body, &reqData)
	if err != nil {
		logger.LogWarm("Failed to nnmarshal req body.", err)
		this.sendErrorResponse(&req, err, startTimeNano, nil)
		return
	}

	err = param.CheckMeta(&reqData.Meta, reqType)
	if err != nil {
		logger.LogWarm(err)
		this.sendErrorResponse(&req, err, startTimeNano, &reqData)
		return
	}

	switch reqType {
	case 0:
		this.handleRpc(&req, &reqData, startTimeNano)
	case 1:
		this.handleStatistic(&req, reqData)
	}
}

func (this *Mq) handleRpc(req *amqp.Delivery, reqData *ReqData, startTimeNano int64) {
	step := reqData.Meta.Step

	reqMetaMap[reqData.Meta.ReqNo] = &reqData.Meta
	defer delete(reqMetaMap, reqData.Meta.ReqNo)

	respBody := this.getResponseBody(reqData)

	err := this.sendResponse(req, respBody)
	if err != nil {
		logger.LogWarm("Failed to publish a message.", err)
	}

	this.reqLog(startTimeNano, step, reqData, respBody)
}

func (this *Mq) handleStatistic(req *amqp.Delivery, reqData ReqData) {
	hander, err := routeCore.GetHander(reqData.Route)
	if err != nil {
		logger.LogWarm(fmt.Sprintf(ERR_FORMAT, ERRNO_ROUTE_NOT_FOUND, err.Error()))
		return
	}

	reflectParam := []reflect.Value{reflect.ValueOf(reqData.Meta), reflect.ValueOf(reqData.Data)}
	hander.Call(reflectParam)
}

func (this *Mq) getResponseBody(reqData *ReqData) ResponseBody {
	hander, err := routeCore.GetHander(reqData.Route)
	if err != nil {
		return ResponseFail(ERRNO_ROUTE_NOT_FOUND, err)
	}

	reflectParam, err := param.GenReflectParam(reqData.Route, reqData.Data)
	if err != nil {
		return ResponseFail(ERRNO_ARGS_ERR, err)
	}

	//business层，加参数reqId和setp
	if reqData.Meta.Step%100 == 0 {
		reflectParam = append([]reflect.Value{reflect.ValueOf(reqData.Meta.ReqNo)}, reflectParam...)
	}

	ret := hander.Call(reflectParam)

	if ret[1].Interface() != nil {
		return ResponseFail(ERRNO_INNER, ret[1].Interface().(error))
	} else {
		return ResponseOk(ret[0].Interface())
	}
}

func (this *Mq) sendErrorResponse(req *amqp.Delivery, err error, startTimeNano int64, reqData *ReqData) {
	respBody := ResponseFail(ERRNO_ARGS_ERR, err)

	err = this.sendResponse(req, respBody)
	if err != nil {
		logger.LogWarm("Failed to publish a message.", err)
	}

	if reqData == nil {
		reqData = &ReqData{
			Meta: ReqMeta{
				ReqTime: time.Now().Unix(),
				ReqNo:   ERR_REQNO,
			},
			Data: map[string]interface{}{},
		}
	}

	this.reqLog(startTimeNano, ERR_STEP, reqData, respBody)
}

func (this *Mq) sendResponse(req *amqp.Delivery, respBody ResponseBody) error {
	body, _ := json.Marshal(respBody)
	response := amqp.Publishing{
		ContentType:   "text/plain",
		CorrelationId: req.CorrelationId,
		Body:          body,
	}
	err := channels[FRONTEND].Publish("", req.ReplyTo, false, false, response)
	if err != nil {
		return err
	}

	req.Ack(false)

	return nil
}

func (this *Mq) CallService(reqNo, serviceQueue, route string, data map[string]interface{}) (interface{}, error) {
	queue, err := channels[SERVICE].QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return nil, err
	}

	msgs, err := channels[SERVICE].Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	reqMetaMap[reqNo].Step += 1

	corrId := randomString(CORRID_LEN)

	req := this.getPublishingReq(*reqMetaMap[reqNo], queue.Name, route, corrId, data)

	err = channels[SERVICE].Publish("", mqConfigs[SERVICE].Queues[serviceQueue], false, false, *req)
	if err != nil {
		return nil, err
	}

	return this.getServiceResponse(msgs, corrId)
}

func (this *Mq) getPublishingReq(reqMeta ReqMeta, queueName, route, corrId string, data map[string]interface{}) *amqp.Publishing {
	reqData := ReqData{
		Meta:  reqMeta,
		Route: route,
		Data:  data,
	}
	body, _ := json.Marshal(reqData)

	return &amqp.Publishing{
		ContentType:   CONTENT_TYPE,
		CorrelationId: corrId,
		ReplyTo:       queueName,
		Body:          body,
	}
}

func (this *Mq) getServiceResponse(msgs <-chan amqp.Delivery, corrId string) (interface{}, error) {
	var respBody ResponseBody

	for msg := range msgs {
		if corrId == msg.CorrelationId {
			err := json.Unmarshal(msg.Body, &respBody)
			if err != nil {
				return nil, err
			}
			if respBody.Code != CODE_SUCCESS {
				return nil, errors.New(respBody.Msg)
			}

			break
		}
	}

	return respBody.Data, nil
}

func (this *Mq) CallMultiService(reqNo string, params []CallParam) ([]string, error) {
	if len(params) == 0 {
		return []string{}, nil
	}

	queue, err := channels[SERVICE].QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Call Multi Service error: %s", err.Error()))
	}

	msgs, err := channels[SERVICE].Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Call Multi Service error: %s", err.Error()))
	}

	corrIdList := map[string]int{}
	for index, param := range params {
		reqMetaMap[reqNo].Step += 1
		reqData := ReqData{
			Meta:  *reqMetaMap[reqNo],
			Route: param.Route,
			Data:  param.Data,
		}
		body, _ := json.Marshal(reqData)

		corrId := randomString(CORRID_LEN)
		req := amqp.Publishing{
			ContentType:   CONTENT_TYPE,
			CorrelationId: corrId,
			ReplyTo:       queue.Name,
			Body:          body,
		}
		err = channels[SERVICE].Publish("", mqConfigs[SERVICE].Queues[param.Module], false, false, req)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Call Multi Service error: %s", err.Error()))
		}

		corrIdList[corrId] = index
	}

	respList := make([]string, len(corrIdList))
	for msg := range msgs {
		index, ok := corrIdList[msg.CorrelationId]
		if ok {
			respList[index] = string(msg.Body)

			delete(corrIdList, msg.CorrelationId)
		}

		if len(corrIdList) == 0 {
			break
		}
	}

	return respList, nil
}

func (this *Mq) genReqLogData(startTimeNano int64, step int, reqData *ReqData, respBody ResponseBody, reqLogData *ReqData) {
	reqLogData.Route = reqData.Route
	reqBody, err := json.Marshal(reqData.Data)
	if err != nil {
		logger.LogWarm(err)
		return
	}

	endTimeNano := time.Now().UnixNano()

	reqLogData.Meta = ReqMeta{
		ReqNo:    reqData.Meta.ReqNo,
		Step:     step,
		Source:   reqData.Meta.Source,
		RemoteIp: reqData.Meta.RemoteIp,
		ReqTime:  reqData.Meta.ReqTime,
	}

	reqLogData.Route = REQLOG_ROUTE

	respData, _ := json.Marshal(respBody.Data)
	reqLogData.Data = map[string]interface{}{
		"module":    this.module,
		"route":     reqData.Route,
		"req_data":  string(reqBody),
		"req_size":  len(reqBody),
		"cost_time": (endTimeNano - startTimeNano) / 1000,
		"code":      respBody.Code,
		"msg":       respBody.Msg,
		"resp_data": string(respData),
		"resp_size": len(respData),
	}
}

func (this *Mq) reqLog(startTimeNano int64, step int, reqData *ReqData, respBody ResponseBody) {
	reqLogData := ReqData{}
	this.genReqLogData(startTimeNano, step, reqData, respBody, &reqLogData)

	mqBody, err := json.Marshal(reqLogData)
	req := amqp.Publishing{
		ContentType: CONTENT_TYPE,
		Body:        mqBody,
	}

	exchangeName := mqConfigs[STATISTIC].Items["reqlog"].ExchangeName
	key := mqConfigs[STATISTIC].Items["reqlog"].Key

	err = channels[STATISTIC].Publish(exchangeName, key, false, false, req)
	if err != nil {
		logger.LogWarm(fmt.Sprintf("Reqlog error: %s", err.Error()))
	}
}

func (this *Mq) Close() {
	for _, channel := range channels {
		channel.Close()
	}

	for _, conn := range conns {
		conn.Close()
	}
}
