package products

import (
	"Databases/mongodb"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sahilm/fuzzy"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type product struct {
	ID        int
	Shop      string
	Articul   string
	URL       string
	Title     string
	Images    []string
	Old_price int
	New_price int
	Sizes     []string
}

type productAnswer struct {
	Status   string    `json:"status"`
	Products []product `json:"products"`
}

var shoplist = [13]string{
	"Adidas",
	"Brandshop",
	"Endclothing",
	"NewBalance",
	"Reebok",
	"Rendevouz",
	"Revolve",
	"Basketshop",
	"Coggles",
	"Eastdane",
	"Funkydunky",
	"Nike",
	"Puma",
}

var dbList = []string{"Cluster1", "Cluster2", "Cluster3", "Cluster4"}

func ProductInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shop := vars["shop"]
	articul := vars["articul"]
	options := options.Find()
	filter := bson.M{"shop": shop, "articul": articul}
	client, ctx, productCollection, cancel, err := mongodb.GetItems(filter, options, dbList)
	defer cancel()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	var tovar product
	var products []product
	for _, product := range productCollection {
		for product.Next(context.Background()) {
			_ = product.Decode(&tovar)
			products = append(products, tovar)
		}
	}
	jsonResponse, err := json.Marshal(productAnswer{"Ok", products})
	w.Write(jsonResponse)
}

func FindProductHandler(w http.ResponseWriter, r *http.Request) {
	var tovar product
	var products []product

	params := r.URL.Query()
	fmt.Println(params)
	shops := params["shop"]
	filter := bson.M{}

	var productCollection []*mongo.Cursor
	var client *mongo.Client
	var ctx context.Context
	// var cancel context.CancelFunc
	var err error

	// defer cancel()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	var tovartitles []string
	if shops != nil {
		if shops[0] != "" {
			for _, shop := range shops {
				filter["shop"] = shop
				client, ctx, productCollection, _, err = mongodb.GetItems(filter, nil, dbList)
				if err != nil {
					panic(err)
				}
				for _, product := range productCollection {
					for product.Next(context.Background()) {
						_ = product.Decode(&tovar)
						products = append(products, tovar)
						tovartitles = append(tovartitles, tovar.Title)
					}
				}
			}
		}
	} else {
		client, ctx, productCollection, _, err = mongodb.GetItems(filter, nil, dbList)
		if err != nil {
			panic(err)
		}
		for _, product := range productCollection {
			for product.Next(context.Background()) {
				_ = product.Decode(&tovar)
				products = append(products, tovar)
				tovartitles = append(tovartitles, tovar.Title)
			}
		}
	}
	title := params["title"]

	var keys fuzzy.Matches
	var filteredProducts []product
	if title != nil {
		keys = fuzzy.Find(title[0], tovartitles)
		for _, key := range keys {
			filteredProducts = append(filteredProducts, products[key.Index])
		}
	} else {
		filteredProducts = append(filteredProducts, products...)
	}
	jsonResponse, err := json.Marshal(productAnswer{"Ok", filteredProducts})
	if err != nil {
		log.Println(err)
	}
	w.Write(jsonResponse)
}

//ProductsHandler for Registration user for getting products
func GetProductsForMainHandler(w http.ResponseWriter, r *http.Request) {
	options := options.Find()
	options.SetLimit(10)
	var productlist []product
	for _, shop := range shoplist {
		filter := bson.M{"shop": shop}
		client, ctx, productCollection, cancel, err := mongodb.GetItems(filter, options, dbList)
		defer cancel()
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()
		var tovar product
		var products []product
		for _, product := range productCollection {
			for product.Next(context.Background()) {
				_ = product.Decode(&tovar)
				products = append(products, tovar)
			}
		}
		productlist = append(productlist, products...)
	}
	jsonResponse, err := json.Marshal(productAnswer{"Ok", productlist})
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(jsonResponse)
	}
}
