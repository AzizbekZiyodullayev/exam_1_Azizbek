package controller

import (
	"market/models"
	"sort"
)

// TASK - 1; Shop cart boyicha default holati time sort bolishi kerak.

func (c *Controller) Sort(req *models.ShopCartGetListRequest) (*models.ShopCartGetListResponse, error) {
	var resp = &models.ShopCartGetListResponse{}
	var orderDateFilter []*models.ShopCart
	getorder, err := c.ShopCartGetList(req)
	if err != nil {
		return nil, err
	}
	for _, ord := range getorder.ShopCarts {
		orderDateFilter = append(orderDateFilter, ord)

	}
	sort.Slice(orderDateFilter, func(i, j int) bool {
		return orderDateFilter[i].Time > orderDateFilter[j].Time
	})
	resp.Count = len(orderDateFilter)
	resp.ShopCarts = orderDateFilter
	return resp, nil
}

// TASK - 2

func (c *Controller) Filter(req *models.ShopCartGetListRequest) ([]*models.ShopCart, error) {
	var orderDateFilter []*models.ShopCart
	getorder, err := c.ShopCartGetList(req)
	if err != nil {
		return nil, err
	}
	for _, ord := range getorder.ShopCarts {
		if ord.Time >= req.FromTime && ord.Time < req.ToTime {
			orderDateFilter = append(orderDateFilter, ord)
		}
	}
	return orderDateFilter, nil
}

// TASK - 3; Client history chiqish kerak. Ya'ni sotib olgan mahsulotlari korsatish kerak

func (c *Controller) HistoryUser(req *models.UserPrimaryKey) (map[string][]models.History, error) {
	var (
		orders   = []models.History{}
		orderMap = make(map[string][]models.History)
	)
	getOrder, err := c.ShopCartGetList(&models.ShopCartGetListRequest{})
	if err != nil {
		return nil, err
	}

	getUser, err := c.UserGetById(&models.UserPrimaryKey{Id: req.Id})
	if err != nil {
		return nil, err
	}

	for _, v := range getOrder.ShopCarts {
		getproduct, err := c.ProductGetById(&models.ProductPrimaryKey{v.ProductId})
		if err != nil {
			return nil, err
		}

		if v.UserId == req.Id {
			if v.Status == true {
				orders = append(orders, models.History{
					ProductName: getproduct.Name,
					Count:       v.Count,
					Total:       v.Count * getproduct.Price,
					Time:        v.Time,
				})
			}
		}
	}
	orderMap[getUser.Name] = orders
	return orderMap, nil
}

// TASK - 4; Client qancha pul mahsulot sotib olganligi haqida hisobot.

func (c *Controller) UserCash(req *models.UserPrimaryKey) (map[string]int, error) {
	user := make(map[string]int)

	getorder, err := c.ShopCartGetList(&models.ShopCartGetListRequest{})
	if err != nil {
		return nil, err
	}

	getuser, err := c.UserGetById(req)

	for _, value := range getorder.ShopCarts {
		if value.UserId == req.Id {
			if value.Status == true {
				getproduct, err := c.ProductGetById(&models.ProductPrimaryKey{Id: value.ProductId})
				if err != nil {
					return nil, err
				}
				user[getuser.Name] += value.Count * getproduct.Price
			}
		}
	}
	return user, nil
}

// TASK - 5; Productlarni Qancha sotilgan boyicha hisobot

func (c *Controller) ProductCountSold() (map[string]int, error) {
	product := make(map[string]int)

	getorder, err := c.ShopCartGetList(&models.ShopCartGetListRequest{})
	if err != nil {
		return nil, err
	}

	for _, value := range getorder.ShopCarts {
		getproduct, err := c.ProductGetById(&models.ProductPrimaryKey{Id: value.ProductId})
		if err != nil {
			return nil, err
		}
		if value.Status == true {
			product[getproduct.Name] += value.Count
		}
	}
	return product, nil
}

// TASK - 6; Top 10 ta sotilayotgan mahsulotlarni royxati.

func (c *Controller) TopProducts() ([]*models.ProductsHistory, error) {
	var (
		prodctsMap = make(map[string]int)
		products   []*models.ProductsHistory
	)

	getOrder, err := c.ShopCartGetList(&models.ShopCartGetListRequest{})
	if err != nil {
		return nil, err
	}

	for _, value := range getOrder.ShopCarts {
		getProduct, err := c.ProductGetById(&models.ProductPrimaryKey{Id: value.ProductId})
		if err != nil {
			return nil, err
		}
		if value.Status == true {
			prodctsMap[getProduct.Name] += value.Count
		}
	}
	for k, v := range prodctsMap {
		products = append(products, &models.ProductsHistory{
			Name:  k,
			Count: v,
		})
	}

	sort.Slice(products, func(i, j int) bool {
		return products[i].Count > products[j].Count
	})

	return products, nil
}

// TASK - 7; Top 10 ta Past sotilayotgan mahsulotlarni royxati.

func (c *Controller) FailureProducts() ([]*models.ProductsHistory, error) {
	var (
		prodctsMap = make(map[string]int)
		products   []*models.ProductsHistory
	)

	getOrder, err := c.ShopCartGetList(&models.ShopCartGetListRequest{})
	if err != nil {
		return nil, err
	}

	for _, value := range getOrder.ShopCarts {
		getProduct, err := c.ProductGetById(&models.ProductPrimaryKey{Id: value.ProductId})
		if err != nil {
			return nil, err
		}
		if value.Status == true {
			prodctsMap[getProduct.Name] += value.Count
		}
	}
	for k, v := range prodctsMap {
		products = append(products, &models.ProductsHistory{
			Name:  k,
			Count: v,
		})
	}
	sort.Slice(products, func(i, j int) bool {
		return products[i].Count < products[j].Count
	})
	return products, nil
}
