package searchRepository

import (
	"bytes"
	"context"
	"discountDealer/conf"
	"discountDealer/logger"
	"discountDealer/models/manticoreModels"
	"discountDealer/x"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

var ManticoreClientContext = context.Background()

type ManticoreClient struct {
	logger *zap.Logger
}

func Init() *ManticoreClient {
	ManticoreClientContext = logger.InsertLogger(
		ManticoreClientContext,
		logger.ManticoreClientKey,
		logger.New("Search repository"),
	)
	client := new(ManticoreClient)
	client.logger = logger.ExtractLogger(ManticoreClientContext, logger.ManticoreClientKey)
	return client
}

func (m *ManticoreClient) FindProduct(title *string) ([]string, error) {
	data := map[string]interface{}{
		"index": conf.Config.ManticoreIndex,
		"query": map[string]interface{}{
			"match": map[string]string{
				"*": *title,
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		m.logger.Error("Can't encode data to json", zap.Error(err))
		return nil, err
	}

	resp, err := http.Post(m.formURL(conf.Config.ManticoreSearchURL), x.JsonContent, bytes.NewBuffer(jsonData))
	if err != nil {
		m.logger.Error("Error while searching products", zap.Error(err))
		return nil, err
	}

	res := new(manticoreModels.ManticoreSearchResp)
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		m.logger.Error("Error while decoding data from json", zap.Error(err))
		return nil, err
	}

	articuls := make([]string, 0)
	for _, product := range res.Hits.Hits {
		articuls = append(articuls, product.Source.Articul)
	}
	return articuls, nil
}

func (m *ManticoreClient) formURL(host string) string {
	return fmt.Sprintf("http://%s", host)
}
