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
	sdk "gitlab.paradise-soft.com.tw/glob/legion-sdk"

	"gitlab.paradise-soft.com.tw/glob/dispatcher"
	"gitlab.paradise-soft.com.tw/glob/tracer"
)

func dynamicScrape(data []byte) (err error) {
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

	legionResp := legionReq.GetDynamicResult()

	const staticCachePath = `/v1/apis/dynamic/cache`
	cacheKey := fmt.Sprintf("[%s][%s]", legionReq.RespTopic, uuid.New().String())
	queryData := url.Values{}
	queryData.Add("key", cacheKey)

	notice := &sdk.Notice{
		UUID: legionReq.UUID,
	}
	notice.InternalURL = fmt.Sprintf("%s%s?%s",
		glob.Config.WWW.InternalHost,
		staticCachePath,
		queryData.Encode(),
	)

	notice.ExternalURL = fmt.Sprintf("%s%s?%s",
		glob.Config.WWW.ExternalHost,
		staticCachePath,
		queryData.Encode(),
	)
	notice.CreatedAt = time.Now()

	var noticeBytes []byte
	noticeBytes, err = json.Marshal(notice)
	if err != nil {
		// internal error
		tracer.Error("internal", err)
		return
	}

	ok := glob.RespCache.SaveDynamic(cacheKey, legionResp)
	if !ok {
		// internal error
		err = errors.New("key exist")
		tracer.Error("internal", err)
		return
	}

	err = dispatcher.Send(legionReq.RespTopic, noticeBytes)
	if err != nil {
		// internal error
		tracer.Error("internal", err)
		return
	}
	return

}

func DispatcherErrHandler(data []byte, err error) {
	if err != nil {
		tracer.Error("Dispatcher", err, data)
	}
}

func dynamicScrapeAPI(ctx *gin.Context) {
	legionReq := &service.LegionRequest{}
	err := ctx.ShouldBindJSON(legionReq)
	if err != nil {
		responseParamError(ctx, err)
		return
	}

	legionResp := legionReq.GetDynamicResult()
	response(ctx, legionResp, 1, glob.ScrapeSuccess, nil)
}

func getDynamicCache(ctx *gin.Context) {
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

	value, ok := glob.RespCache.GetDynamic(req.Key)
	if !ok {
		err = errors.New("key does not exist")
		responseParamError(ctx, err)
		return
	}

	glob.RespCache.DeleteDynamic(req.Key)
	response(ctx, value, 1, glob.ScrapeSuccess, nil)
}
