package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lwmacct/251117-go-ddd-template/internal/adapters/http/response"
	"github.com/lwmacct/251117-go-ddd-template/internal/application/product"
	productDomain "github.com/lwmacct/251117-go-ddd-template/internal/domain/product"
)

// ListProductsQuery 产品列表查询参数
type ListProductsQuery struct {
	response.PaginationQueryDTO
}

// ToQuery 转换为 Application 层 Query 对象
func (q *ListProductsQuery) ToQuery() product.ListProductsQuery {
	return product.ListProductsQuery{
		Offset: q.GetOffset(),
		Limit:  q.GetLimit(),
	}
}

// ProductHandler 产品管理 Handler
type ProductHandler struct {
	createHandler *product.CreateHandler
	updateHandler *product.UpdateHandler
	deleteHandler *product.DeleteHandler
	getHandler    *product.GetHandler
	listHandler   *product.ListHandler
}

// NewProductHandler 创建产品管理 Handler
func NewProductHandler(
	createHandler *product.CreateHandler,
	updateHandler *product.UpdateHandler,
	deleteHandler *product.DeleteHandler,
	getHandler *product.GetHandler,
	listHandler *product.ListHandler,
) *ProductHandler {
	return &ProductHandler{
		createHandler: createHandler,
		updateHandler: updateHandler,
		deleteHandler: deleteHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
	}
}

// Create 创建产品
//
//	@Summary		创建产品
//	@Description	系统管理员创建可订阅的产品
//	@Tags			Admin - Product Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		product.CreateProductDTO					true	"产品信息"
//	@Success		201		{object}	response.DataResponse[product.ProductDTO]	"产品创建成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误或产品名已存在"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		500		{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/admin/products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req product.CreateProductDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.createHandler.Handle(c.Request.Context(), product.CreateProductCommand(req))
	if err != nil {
		if errors.Is(err, productDomain.ErrProductNameExists) {
			response.Conflict(c, "产品名称已存在")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, "产品创建成功", result)
}

// List 产品列表
//
//	@Summary		产品列表
//	@Description	分页获取产品列表
//	@Tags			Admin - Product Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			params	query		ListProductsQuery							false	"查询参数"
//	@Success		200		{object}	response.PagedResponse[product.ProductDTO]	"产品列表"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		500		{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/admin/products [get]
func (h *ProductHandler) List(c *gin.Context) {
	var query ListProductsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.listHandler.Handle(c.Request.Context(), query.ToQuery())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := response.NewPaginationMeta(int(result.Total), query.GetPage(), query.GetLimit())
	response.List(c, "获取成功", result.Items, meta)
}

// Get 产品详情
//
//	@Summary		产品详情
//	@Description	获取产品详细信息
//	@Tags			Admin - Product Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int											true	"产品ID"
//	@Success		200	{object}	response.DataResponse[product.ProductDTO]	"产品详情"
//	@Failure		401	{object}	response.ErrorResponse						"未授权"
//	@Failure		403	{object}	response.ErrorResponse						"权限不足"
//	@Failure		404	{object}	response.ErrorResponse						"产品不存在"
//	@Failure		500	{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/admin/products/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的产品ID")
		return
	}

	result, err := h.getHandler.Handle(c.Request.Context(), product.GetProductQuery{
		ID: uint(id),
	})
	if err != nil {
		if errors.Is(err, productDomain.ErrProductNotFound) {
			response.NotFound(c, "产品不存在")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "获取成功", result)
}

// Update 更新产品
//
//	@Summary		更新产品
//	@Description	更新产品信息
//	@Tags			Admin - Product Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		int											true	"产品ID"
//	@Param			request	body		product.UpdateProductDTO					true	"更新信息"
//	@Success		200		{object}	response.DataResponse[product.ProductDTO]	"更新成功"
//	@Failure		400		{object}	response.ErrorResponse						"参数错误"
//	@Failure		401		{object}	response.ErrorResponse						"未授权"
//	@Failure		403		{object}	response.ErrorResponse						"权限不足"
//	@Failure		404		{object}	response.ErrorResponse						"产品不存在"
//	@Failure		409		{object}	response.ErrorResponse						"产品名称已存在"
//	@Failure		500		{object}	response.ErrorResponse						"服务器内部错误"
//	@Router			/api/admin/products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的产品ID")
		return
	}

	var req product.UpdateProductDTO
	if err = c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.updateHandler.Handle(c.Request.Context(), product.UpdateProductCommand{
		ID:          uint(id),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Status:      req.Status,
	})
	if err != nil {
		if errors.Is(err, productDomain.ErrProductNotFound) {
			response.NotFound(c, "产品不存在")
			return
		}
		if errors.Is(err, productDomain.ErrProductNameExists) {
			response.Conflict(c, "产品名称已存在")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "更新成功", result)
}

// Delete 删除产品
//
//	@Summary		删除产品
//	@Description	删除产品（软删除）
//	@Tags			Admin - Product Management
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		int						true	"产品ID"
//	@Success		200	{object}	response.MessageResponse	"删除成功"
//	@Failure		401	{object}	response.ErrorResponse	"未授权"
//	@Failure		403	{object}	response.ErrorResponse	"权限不足"
//	@Failure		404	{object}	response.ErrorResponse	"产品不存在"
//	@Failure		500	{object}	response.ErrorResponse	"服务器内部错误"
//	@Router			/api/admin/products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "无效的产品ID")
		return
	}

	if err := h.deleteHandler.Handle(c.Request.Context(), product.DeleteProductCommand{
		ID: uint(id),
	}); err != nil {
		if errors.Is(err, productDomain.ErrProductNotFound) {
			response.NotFound(c, "产品不存在")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, "删除成功", nil)
}
