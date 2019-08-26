package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.paradise-soft.com.tw/dwh/legion/glob"
	"gitlab.paradise-soft.com.tw/dwh/legion/service"
	"gitlab.paradise-soft.com.tw/glob/dispatcher"
	"gitlab.paradise-soft.com.tw/glob/tracer"
)

func staticScrape(data []byte) (err error) {
	legionReq := &service.LegionRequest{}
	if err = json.Unmarshal(data, legionReq); err != nil {
		return
	}

	now := time.Now()
	if legionReq.SentAt.IsZero() || legionReq.SentAt.After(now) || legionReq.SentAt.Add(service.ExpiredTime).Before(now) {
		err = fmt.Errorf("task expired sent at %v", legionReq.SentAt)
		return
	}

	err = legionReq.CheckKafka()
	if err != nil {
		return
	}

	legionResp := legionReq.GetStaticResponse()
	// var legionRespBytes []byte
	// legionRespBytes, err = json.Marshal(legionResp)
	// if err != nil {
	// 	// internal error
	// 	tracer.Error("internal", err)
	// 	return
	// }

	const staticCachePath = `/v1/apis/static/cache`
	cacheKey := fmt.Sprintf("[%s][%s]", legionReq.RespTopic, uuid.New().String())
	queryData := url.Values{}
	queryData.Add("key", cacheKey)

	legionKafkaResp := &service.Notice{}
	legionKafkaResp.InternalURL = fmt.Sprintf("%s%s%s?%s",
		glob.Config.WWW.InternalHost,
		glob.Config.WWW.Addr,
		staticCachePath,
		queryData.Encode(),
	)

	legionKafkaResp.ExternalURL = fmt.Sprintf("%s%s%s?%s",
		glob.Config.WWW.ExternalHost,
		glob.Config.WWW.Addr,
		staticCachePath,
		queryData.Encode(),
	)

	var legionKafkaRespBytes []byte
	legionKafkaRespBytes, err = json.Marshal(legionKafkaResp)
	if err != nil {
		// internal error
		tracer.Error("internal", err)
		return
	}

	ok := glob.RespCache.SaveStatic(cacheKey, legionResp)
	if !ok {
		// internal error
		err = errors.New("key exist")
		tracer.Error("internal", err)
		return
	}

	err = dispatcher.Send(
		legionReq.RespTopic,
		legionKafkaRespBytes,
		dispatcher.ProducerAddErrHandler(func(value []byte, err error) {
			tracer.Error("sdk", err)
		}),
	)

	if err != nil {
		// internal error
		tracer.Error("internal", err)
		return
	}
	return
}

func staticScrapeAPI(ctx *gin.Context) {
	legionReq := &service.LegionRequest{}
	err := ctx.ShouldBindJSON(legionReq)
	if err != nil {
		responseParamError(ctx, err)
		return
	}

	legionResp := legionReq.GetStaticResponse()
	response(ctx, legionResp, 1, glob.ScrapeSuccess, nil)
}

func getStaticCache(ctx *gin.Context) {
	var err error
	req := &service.CacheRequest{}

	err = ctx.ShouldBindQuery(req)
	if err != nil {
		responseParamError(ctx, err)
		return
	}

	if req.Key == "" {
		err = errors.New("key is empty")
		responseParamError(ctx, err)
		return
	}

	value, ok := glob.RespCache.GetStatic(req.Key)
	if !ok {
		err = errors.New("key does not exist")
		responseParamError(ctx, err)
		return
	}

	// glob.RespCache.DeleteStatic(req.Key)
	response(ctx, value, 1, glob.ScrapeSuccess, nil)
}
