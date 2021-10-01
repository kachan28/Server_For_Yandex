package products

import (
	"discountDealer/logger"
	"discountDealer/models/productModels"
	"discountDealer/models/queryModels"
	"discountDealer/repository/productRepository"
	"discountDealer/repository/productRepository/filtersAndOptions"
	"discountDealer/x"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetProduct(c *fiber.Ctx) error {
	log := logger.New("Get product handler")
	product := new(productModels.Product)

	product.Articul = c.Params(x.ArticulParameter)

	err := productRepository.Init().Get(c.UserContext(), product)
	if err != nil {
		log.Error("Can't get product", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func Find(c *fiber.Ctx) error {
	log := logger.New("Find products handler")
	//parsing title query
	query := &queryModels.FindQuery{Root: queryModels.NewRoot(log)}
	err := c.QueryParser(query)
	if err != nil {
		log.Error("Error while parsing query", zap.Error(err))
	}
	//parsing equal filters
	qEqFilter := &queryModels.FindQueryEqFilter{Root: queryModels.NewRoot(log)}
	err = c.QueryParser(qEqFilter)
	if err != nil {
		log.Error("Error while parsing equal filters", zap.Error(err))
	}
	//parsing comparative filters
	qCompFilter := &queryModels.FindQueryCompFilter{Root: queryModels.NewRoot(log)}
	err = c.QueryParser(qCompFilter)
	if err != nil {
		log.Error("Error while parsing comparative filters", zap.Error(err))
	}
	//parsing sort options
	qOptions := &queryModels.FindQuerySort{Root: queryModels.NewRoot(log)}
	err = c.QueryParser(qOptions)
	if err != nil {
		log.Error("Error while parsing sort query", zap.Error(err))
	}
	//parsing num page
	qPage := &queryModels.FindQueryPage{Root: queryModels.NewRoot(log)}
	err = c.QueryParser(qPage)
	if err != nil {
		log.Error("Error while parsing page query", zap.Error(err))
	}
	//creating filter
	filter := filtersAndOptions.NewFilter()
	//adding equal filters iterating struct
	qEqFilter.Parse(qEqFilter, filter)
	//adding comparative filters
	qCompFilter.Parse(qCompFilter, filter)
	//finding by title
	query.Parse(query, filter)
	//creating options
	opts := filtersAndOptions.NewOpts()
	//setting sort
	qOptions.Parse(qOptions, opts)
	//setting pagination
	qPage.Parse(qPage, opts)

	products := new([]productModels.Product)
	productRepository.Init().GetByFilter(products, filter, opts)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"products": products,
	})
}
