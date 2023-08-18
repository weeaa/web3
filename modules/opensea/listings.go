package opensea

import (
	"github.com/foundVanting/opensea-stream-go/entity"
	"github.com/mitchellh/mapstructure"
	"github.com/weeaa/nft/pkg/logger"
	"math/big"
	"time"
)

func (c *Client) GetListings(collections []string) {
	for _, collection := range collections {
		go func(slug string) {
			var err error
			c.StreamClient.OnItemListed(slug, func(response any) {
				var ItemListedEvent entity.ItemListedEvent
				l := &Listing{}

				if err = mapstructure.Decode(response, &ItemListedEvent); err != nil {
					logger.LogError(moduleName, err)
				}

				l.Item = ItemListedEvent.Payload.PayloadItemAndColl.Item.Metadata.Name
				l.Item = "[" + l.Item + "]" + "(" + ItemListedEvent.Payload.PayloadItemAndColl.Item.Permalink + ")"
				l.Seller = ItemListedEvent.Payload.Maker.Address
				l.Seller = "[" + l.Seller + "]" + "(" + "https://opensea.io/" + l.Seller + ")"

				wei := new(big.Int)
				wei.SetString(ItemListedEvent.Payload.BasePrice, 10)
				l.PriceInfo.Price = weiToEther(wei)

				l.Collection = ItemListedEvent.Payload.PayloadItemAndColl.Collection.Slug
				l.CollectionLink = "https://opensea.io/collection/" + ItemListedEvent.Payload.PayloadItemAndColl.Collection.Slug
				l.Symbol = ItemListedEvent.Payload.PaymentToken.Symbol

				l.Image = ItemListedEvent.Payload.PayloadItemAndColl.Item.Metadata.ImageUrl
				l.Timestamp = time.Now().Unix()

				l.PriceInfo.Floor, err = c.GetFloor(l.Collection)
				if err != nil {
					logger.LogError(moduleName, err)
				}

				var floorBelow20 float64
				var double float64
				floorBelow20 = l.PriceInfo.Floor / 10
				double = floorBelow20 * 2
				floorMinus20 := l.PriceInfo.Floor - double

				l.PriceInfo.PriceBefore, _ = l.PriceInfo.Price.Float64()

				if l.PriceInfo.PriceBefore <= floorMinus20 {
					//todo: add webhook
				}
			})
		}(collection)
	}
}
