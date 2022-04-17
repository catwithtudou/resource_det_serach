package service

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"resource_det_search/api"
	v1 "resource_det_search/api/v1"
	"resource_det_search/internal/biz"
	"resource_det_search/internal/utils"
)

type DimensionService struct {
	log *zap.SugaredLogger
	dm  biz.IDimensionUsecase
}

func NewDimensionService(dimension biz.IDimensionUsecase, logger *zap.SugaredLogger) *DimensionService {
	return &DimensionService{
		log: logger,
		dm:  dimension,
	}
}

func (d *DimensionService) AddUserDm(c *gin.Context) {
	uid, _ := c.Get("uid")

	var req v1.DimensionAddUserDmReq
	if err := c.ShouldBind(&req); err != nil {
		d.log.Errorf("[DimensionService-AddUserDm]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	if !utils.CheckUserType(req.Type) || len(req.Name) > 50 {
		d.log.Errorf("[DimensionService-AddUserDm]illegal params")
		c.JSON(http.StatusOK, api.FormIllegalErr)
		return
	}

	err := d.dm.AddUserDm(c, &biz.Dimension{
		Uid:  uid.(uint),
		Type: req.Type,
		Name: req.Name,
	})
	if err != nil {
		d.log.Errorf("[DimensionService-AddUserDm]failed to AddUserDm:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	c.JSON(http.StatusOK, api.Success)
	return
}

func (d *DimensionService) GetUserDm(c *gin.Context) {
	var req v1.DimensionGetUserDmReq
	var uid uint
	if _ = c.ShouldBind(&req); req.Uid > 0 {
		uid = req.Uid
	}

	if uid <= 0 {
		getUid, _ := c.Get("uid")
		uid = getUid.(uint)
	}

	data, err := d.dm.GetUserDm(c, uid)
	if err != nil {
		d.log.Errorf("[DimensionService-GetUserDm]failed to GetUserDm:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	resp := &v1.DimensionGetUserDmResp{
		RespCommon: api.Success,
		Data:       make(map[string][]*v1.DimensionUserDmData),
	}
	for k, v := range data {
		if _, ok := resp.Data[k]; !ok {
			resp.Data[k] = make([]*v1.DimensionUserDmData, 0, len(v))
		}
		for _, vv := range v {
			resp.Data[k] = append(resp.Data[k], &v1.DimensionUserDmData{
				Id:   vv.ID,
				Name: vv.Name,
			})
		}
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (d *DimensionService) UpdateUserDm(c *gin.Context) {
	uid, _ := c.Get("uid")

	var req v1.DimensionUpdateUserDmReq
	if err := c.ShouldBind(&req); err != nil {
		d.log.Errorf("[DimensionService-UpdateUserDm]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	if len(req.Name) > 50 {
		d.log.Errorf("[DimensionService-UpdateUserDm]illegal params")
		c.JSON(http.StatusOK, api.FormIllegalErr)
		return
	}

	err := d.dm.UpdateUserDm(c, req.Did, req.Name, uid.(uint))
	if err != nil {
		d.log.Errorf("[DimensionService-UpdateUserDm]failed to UpdateUserDm:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	c.JSON(http.StatusOK, api.Success)
	return
}

func (d *DimensionService) DeleteUserDm(c *gin.Context) {
	uid, _ := c.Get("uid")

	var req v1.DimensionDeleteUserDmReq
	if err := c.ShouldBind(&req); err != nil {
		d.log.Errorf("[DimensionService-DeleteUserDm]failed to bind:err=[%+v]", err)
		c.JSON(http.StatusOK, api.FormEmptyErr)
		return
	}

	err := d.dm.DeleteUserDm(c, req.Did, uid.(uint))
	if err != nil {
		d.log.Errorf("[DimensionService-DeleteUserDm]failed to UpdateUserDm:err=[%+v]", err)
		c.JSON(http.StatusOK, api.DefaultErr)
		return
	}

	c.JSON(http.StatusOK, api.Success)
	return
}
