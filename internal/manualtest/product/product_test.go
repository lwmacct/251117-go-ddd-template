package product_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lwmacct/251117-go-ddd-template/internal/application/product"
	"github.com/lwmacct/251117-go-ddd-template/internal/manualtest"
)

// productBasePath 返回产品 API 基础路径
func productBasePath() string {
	return "/api/admin/products"
}

// productPath 返回单个产品 API 路径
func productPath(id uint) string {
	return fmt.Sprintf("%s/%d", productBasePath(), id)
}

// TestProductCRUD 测试产品完整 CRUD 流程。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestProductCRUD ./internal/manualtest/product/
func TestProductCRUD(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 测试 1: 创建产品
	t.Log("\n测试 1: 创建产品")
	createReq := product.CreateProductDTO{
		Name:        fmt.Sprintf("testproduct_%d", time.Now().UnixNano()),
		Description: "这是一个测试产品",
		Price:       99.99,
		Status:      "active",
	}
	createdProduct, err := manualtest.Post[product.ProductDTO](c, productBasePath(), createReq)
	require.NoError(t, err, "创建产品失败")
	productID := createdProduct.ID
	t.Cleanup(func() {
		if productID > 0 {
			_ = c.Delete(productPath(productID))
		}
	})
	t.Logf("  创建成功! 产品 ID: %d", productID)
	t.Logf("  名称: %s, 价格: %.2f", createdProduct.Name, createdProduct.Price)

	// 验证字段
	assert.Equal(t, createReq.Name, createdProduct.Name, "名称不匹配")
	assert.Equal(t, createReq.Description, createdProduct.Description, "描述不匹配")
	assert.InEpsilon(t, createReq.Price, createdProduct.Price, 0.01, "价格不匹配")
	assert.Equal(t, "active", createdProduct.Status, "默认状态应为 active")

	// 测试 2: 获取产品列表
	t.Log("\n测试 2: 获取产品列表")
	products, meta, err := manualtest.GetList[product.ProductDTO](c, productBasePath(), map[string]string{
		"page":  "1",
		"limit": "10",
	})
	require.NoError(t, err, "获取产品列表失败")
	t.Logf("  产品数量: %d", len(products))
	if meta != nil {
		t.Logf("  总数: %d, 总页数: %d", meta.Total, meta.TotalPages)
	}

	// 验证列表中包含创建的产品
	productIDs := manualtest.ExtractIDs(products, func(p product.ProductDTO) uint { return p.ID })
	assert.Contains(t, productIDs, productID, "列表中应包含新创建的产品")

	// 测试 3: 获取产品详情
	t.Log("\n测试 3: 获取产品详情")
	productDetail, err := manualtest.Get[product.ProductDTO](c, productPath(productID), nil)
	require.NoError(t, err, "获取产品详情失败")
	t.Logf("  详情: %s - %s", productDetail.Name, productDetail.Description)
	assert.Equal(t, productID, productDetail.ID, "产品 ID 不匹配")

	// 测试 4: 更新产品
	t.Log("\n测试 4: 更新产品")
	newPrice := 199.99
	newStatus := "inactive"
	updateReq := product.UpdateProductDTO{
		Price:  &newPrice,
		Status: &newStatus,
	}
	updatedProduct, err := manualtest.Put[product.ProductDTO](c, productPath(productID), updateReq)
	require.NoError(t, err, "更新产品失败")
	t.Logf("  更新成功! 价格: %.2f, 状态: %s", updatedProduct.Price, updatedProduct.Status)

	// 验证更新后的字段
	assert.InEpsilon(t, newPrice, updatedProduct.Price, 0.01, "价格未更新")
	assert.Equal(t, newStatus, updatedProduct.Status, "状态未更新")

	// 测试 5: 删除产品
	t.Log("\n测试 5: 删除产品")
	err = c.Delete(productPath(productID))
	require.NoError(t, err, "删除产品失败")
	t.Log("  删除成功!")

	// 标记已删除，避免 t.Cleanup 重复删除
	productID = 0

	t.Log("\n产品 CRUD 流程测试完成!")
}

// TestProductListPagination 测试产品列表分页。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestProductListPagination ./internal/manualtest/product/
func TestProductListPagination(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	t.Log("测试产品列表分页...")

	// 第一页
	page1, meta1, err := manualtest.GetList[product.ProductDTO](c, productBasePath(), map[string]string{
		"page":  "1",
		"limit": "5",
	})
	require.NoError(t, err, "获取第一页失败")
	t.Logf("  第 1 页: %d 条记录", len(page1))
	if meta1 != nil {
		t.Logf("  总数: %d, 总页数: %d", meta1.Total, meta1.TotalPages)
	}

	// 验证总数
	assert.Positive(t, meta1.Total, "总数应有效")

	// 如果有第二页，获取并验证
	if meta1.TotalPages > 1 {
		page2, meta2, err := manualtest.GetList[product.ProductDTO](c, productBasePath(), map[string]string{
			"page":  "2",
			"limit": "5",
		})
		require.NoError(t, err, "获取第二页失败")
		t.Logf("  第 2 页: %d 条记录", len(page2))
		if meta2 != nil {
			assert.Equal(t, meta1.Total, meta2.Total, "两页总数应一致")
		}

		// 验证无重复 ID
		ids1 := manualtest.ExtractIDs(page1, func(p product.ProductDTO) uint { return p.ID })
		ids2 := manualtest.ExtractIDs(page2, func(p product.ProductDTO) uint { return p.ID })
		for _, id := range ids1 {
			assert.NotContains(t, ids2, id, "两页不应有相同 ID")
		}
	}

	t.Log("\n分页测试完成!")
}

// TestProductDuplicateName 测试重复名称处理。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestProductDuplicateName ./internal/manualtest/product/
func TestProductDuplicateName(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 步骤 1: 创建第一个产品
	t.Log("\n步骤 1: 创建第一个产品")
	productName := fmt.Sprintf("duplicate_%d", time.Now().UnixNano())
	createReq := product.CreateProductDTO{
		Name:        productName,
		Description: "第一个产品",
		Price:       10.0,
	}
	createdProduct, err := manualtest.Post[product.ProductDTO](c, productBasePath(), createReq)
	require.NoError(t, err, "创建第一个产品失败")
	productID := createdProduct.ID
	t.Cleanup(func() {
		if productID > 0 {
			_ = c.Delete(productPath(productID))
		}
	})
	t.Logf("  第一个产品创建成功! ID: %d", productID)

	// 步骤 2: 尝试创建同名产品（应失败）
	t.Log("\n步骤 2: 尝试创建同名产品（应失败）")
	duplicateReq := product.CreateProductDTO{
		Name:        productName, // 相同名称
		Description: "第二个产品",
		Price:       20.0,
	}
	_, err = manualtest.Post[product.ProductDTO](c, productBasePath(), duplicateReq)
	require.Error(t, err, "创建同名产品应该失败")
	t.Logf("  预期失败: %v", err)

	t.Log("\n重复名称处理测试完成!")
}

// TestProductStatusToggle 测试状态切换。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestProductStatusToggle ./internal/manualtest/product/
func TestProductStatusToggle(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 创建测试产品
	t.Log("\n步骤 1: 创建测试产品")
	createReq := product.CreateProductDTO{
		Name:   fmt.Sprintf("statusproduct_%d", time.Now().UnixNano()),
		Price:  50.0,
		Status: "active",
	}
	createdProduct, err := manualtest.Post[product.ProductDTO](c, productBasePath(), createReq)
	require.NoError(t, err, "创建产品失败")
	productID := createdProduct.ID
	t.Cleanup(func() {
		if productID > 0 {
			_ = c.Delete(productPath(productID))
		}
	})
	t.Logf("  创建成功! 初始状态: %s", createdProduct.Status)
	assert.Equal(t, "active", createdProduct.Status, "初始状态应为 active")

	// 测试 2: 更新为 inactive
	t.Log("\n步骤 2: 更新状态为 inactive")
	inactiveStatus := "inactive"
	statusReq := product.UpdateProductDTO{
		Status: &inactiveStatus,
	}
	updatedProduct, err := manualtest.Put[product.ProductDTO](c, productPath(productID), statusReq)
	require.NoError(t, err, "更新状态失败")
	t.Logf("  状态更新为: %s", updatedProduct.Status)
	assert.Equal(t, "inactive", updatedProduct.Status, "状态应为 inactive")

	// 测试 3: 恢复为 active
	t.Log("\n步骤 3: 恢复状态为 active")
	activeStatus := "active"
	statusReq.Status = &activeStatus
	updatedProduct, err = manualtest.Put[product.ProductDTO](c, productPath(productID), statusReq)
	require.NoError(t, err, "恢复状态失败")
	t.Logf("  状态恢复为: %s", updatedProduct.Status)
	assert.Equal(t, "active", updatedProduct.Status, "状态应为 active")

	t.Log("\n状态切换测试完成!")
}

// TestProductPriceUpdate 测试价格更新。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestProductPriceUpdate ./internal/manualtest/product/
func TestProductPriceUpdate(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 创建测试产品
	t.Log("\n步骤 1: 创建测试产品")
	createReq := product.CreateProductDTO{
		Name:  fmt.Sprintf("priceproduct_%d", time.Now().UnixNano()),
		Price: 100.0,
	}
	createdProduct, err := manualtest.Post[product.ProductDTO](c, productBasePath(), createReq)
	require.NoError(t, err, "创建产品失败")
	productID := createdProduct.ID
	t.Cleanup(func() {
		if productID > 0 {
			_ = c.Delete(productPath(productID))
		}
	})
	t.Logf("  创建成功! 初始价格: %.2f", createdProduct.Price)
	assert.InEpsilon(t, 100.0, createdProduct.Price, 0.01, "初始价格应为 100.0")

	// 测试 2: 更新价格
	t.Log("\n步骤 2: 更新价格")
	prices := []float64{149.99, 199.99, 99.99}
	for i, newPrice := range prices {
		priceReq := product.UpdateProductDTO{
			Price: &newPrice,
		}
		updatedProduct, err := manualtest.Put[product.ProductDTO](c, productPath(productID), priceReq)
		require.NoError(t, err, "更新价格失败")
		t.Logf("  更新 %d: 新价格 %.2f", i+1, updatedProduct.Price)
		assert.InEpsilon(t, newPrice, updatedProduct.Price, 0.01, "价格未正确更新")
	}

	t.Log("\n价格更新测试完成!")
}

// TestProductWithInvalidData 测试无效数据处理。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestProductWithInvalidData ./internal/manualtest/product/
func TestProductWithInvalidData(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 测试 1: 名称太短
	t.Log("\n测试 1: 名称太短应失败")
	invalidReq := product.CreateProductDTO{
		Name:  "x", // 太短
		Price: 10.0,
	}
	_, err := manualtest.Post[product.ProductDTO](c, productBasePath(), invalidReq)
	require.Error(t, err, "名称太短应该失败")
	t.Logf("  预期失败: %v", err)

	// 测试 2: 负价格
	t.Log("\n测试 2: 负价格应失败")
	negativePrice := -10.0
	invalidReq2 := product.CreateProductDTO{
		Name:  "valid_name",
		Price: negativePrice,
	}
	_, err = manualtest.Post[product.ProductDTO](c, productBasePath(), invalidReq2)
	require.Error(t, err, "负价格应该失败")
	t.Logf("  预期失败: %v", err)

	// 测试 3: 无效状态
	t.Log("\n测试 3: 无效状态应失败")
	invalidStatus := "invalid_status_value"
	statusReq := product.UpdateProductDTO{
		Status: &invalidStatus,
	}
	// 先创建一个正常产品
	createReq := product.CreateProductDTO{
		Name:  fmt.Sprintf("invalidstatus_%d", time.Now().UnixNano()),
		Price: 10.0,
	}
	createdProduct, err := manualtest.Post[product.ProductDTO](c, productBasePath(), createReq)
	require.NoError(t, err)
	productID := createdProduct.ID
	t.Cleanup(func() {
		if productID > 0 {
			_ = c.Delete(productPath(productID))
		}
	})

	_, err = manualtest.Put[product.ProductDTO](c, productPath(productID), statusReq)
	require.Error(t, err, "无效状态应该失败")
	t.Logf("  预期失败: %v", err)

	t.Log("\n无效数据处理测试完成!")
}

// TestProductFilterByStatus 测试按状态筛选产品。
//
// 手动运行:
//
//	MANUAL=1 go test -v -run TestProductFilterByStatus ./internal/manualtest/product/
func TestProductFilterByStatus(t *testing.T) {
	c := manualtest.LoginAsAdmin(t)

	// 创建 active 产品
	t.Log("\n步骤 1: 创建 active 产品")
	activeReq := product.CreateProductDTO{
		Name:   fmt.Sprintf("active_%d", time.Now().UnixNano()),
		Price:  100.0,
		Status: "active",
	}
	activeProduct, err := manualtest.Post[product.ProductDTO](c, productBasePath(), activeReq)
	require.NoError(t, err)
	activeProductID := activeProduct.ID
	t.Cleanup(func() {
		if activeProductID > 0 {
			_ = c.Delete(productPath(activeProductID))
		}
	})

	// 创建 inactive 产品
	t.Log("\n步骤 2: 创建 inactive 产品")
	inactiveReq := product.CreateProductDTO{
		Name:   fmt.Sprintf("inactive_%d", time.Now().UnixNano()),
		Price:  50.0,
		Status: "inactive",
	}
	inactiveProduct, err := manualtest.Post[product.ProductDTO](c, productBasePath(), inactiveReq)
	require.NoError(t, err)
	inactiveProductID := inactiveProduct.ID
	t.Cleanup(func() {
		if inactiveProductID > 0 {
			_ = c.Delete(productPath(inactiveProductID))
		}
	})

	// 测试 3: 筛选 active 产品
	t.Log("\n步骤 3: 筛选 active 状态的产品")
	activeProducts, _, err := manualtest.GetList[product.ProductDTO](c, productBasePath(), map[string]string{
		"status": "active",
		"limit":  "100",
	})
	require.NoError(t, err, "获取 active 产品失败")
	t.Logf("  active 产品数量: %d", len(activeProducts))

	// 验证所有产品都是 active 状态
	for _, p := range activeProducts {
		assert.Equal(t, "active", p.Status, "筛选结果应只包含 active 产品")
	}

	// 测试 4: 筛选 inactive 产品
	t.Log("\n步骤 4: 筛选 inactive 状态的产品")
	inactiveProducts, _, err := manualtest.GetList[product.ProductDTO](c, productBasePath(), map[string]string{
		"status": "inactive",
		"limit":  "100",
	})
	require.NoError(t, err, "获取 inactive 产品失败")
	t.Logf("  inactive 产品数量: %d", len(inactiveProducts))

	// 验证所有产品都是 inactive 状态
	for _, p := range inactiveProducts {
		assert.Equal(t, "inactive", p.Status, "筛选结果应只包含 inactive 产品")
	}

	t.Log("\n状态筛选测试完成!")
}
