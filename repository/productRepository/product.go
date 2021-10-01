package productRepository

import (
	"context"
	"discountDealer/conf"
	"discountDealer/logger"
	"discountDealer/models/productModels"
	"discountDealer/repository/productRepository/filtersAndOptions"
	"fmt"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var ProductRepositoryContext = context.Background()

type productDB struct {
	collection *mgm.Collection
}

func Init() productDB {
	ProductRepositoryContext = logger.InsertLogger(
		ProductRepositoryContext,
		logger.ProductsRepositoryKey,
		logger.New("Product Repository"),
	)
	err := mgm.SetDefaultConfig(nil, conf.Config.ProductsDBName, options.Client().ApplyURI(makeConnectionURL()))
	if err != nil {
		logger.ExtractLogger(
			ProductRepositoryContext,
			logger.ProductsRepositoryKey,
		).Fatal("Can't connect to products database", zap.Error(err))
	}

	db := productDB{}
	db.collection = mgm.Coll(&productModels.Product{})
	return db
}

func (p productDB) Get(ctx context.Context, product *productModels.Product) error {
	log := logger.ExtractLogger(
		ProductRepositoryContext,
		logger.ProductsRepositoryKey,
	)
	log.Info("Getting product with ID", zap.String("ID", product.Articul))

	result := p.collection.FindOne(ctx, map[string]string{"articul": product.Articul})
	err := result.Err()
	if err != nil {
		log.Error("Error while getting product by ID", zap.Error(err))
		return err
	}

	result.Decode(product)
	return nil
}

func (p productDB) GetByFilter(products *[]productModels.Product, filters *filtersAndOptions.Filter, opts *filtersAndOptions.Opts) error {
	log := logger.ExtractLogger(
		ProductRepositoryContext,
		logger.ProductsRepositoryKey,
	)
	log.Info("Getting product(s) by filters", zap.Any("filters", filters.ToBson()))

	err := p.collection.SimpleFind(products, filters.ToBson(), opts.ToBson())
	if err != nil {
		log.Error("Error while getting products", zap.Error(err))
		return err
	}
	return nil
}

func (p productDB) Insert(ctx context.Context, product *productModels.Product) error {
	log := logger.ExtractLogger(
		ProductRepositoryContext,
		logger.ProductsRepositoryKey,
	)

	log.Info("Inserting product", zap.Any("product", product))

	_, err := p.collection.InsertOne(ctx, product)
	if err != nil {
		log.Error("Error while inserting product", zap.Error(err))
		return err
	}

	return nil
}

func (p productDB) Update(product *productModels.Product) error {
	log := logger.ExtractLogger(
		ProductRepositoryContext,
		logger.ProductsRepositoryKey,
	)
	log.Info("Updating product", zap.Any("product", product))

	err := p.collection.Update(product)
	if err != nil {
		log.Error("Error while updating product", zap.Error(err))
		return err
	}

	return nil
}

func (p productDB) Delete(product *productModels.Product) error {
	log := logger.ExtractLogger(
		ProductRepositoryContext,
		logger.ProductsRepositoryKey,
	)
	log.Info("Deleting product", zap.Any("product", product))

	err := p.collection.Delete(product)
	if err != nil {
		log.Error("Error while deleting product", zap.Any("product", product))
		return err
	}

	return nil
}

func makeConnectionURL() string {
	return fmt.Sprintf(
		"mongodb://%v:%v@%v",
		conf.Config.ProductsDBLogin,
		conf.Config.ProductsDBPassword,
		conf.Config.ProductsDBHost,
	)
}
