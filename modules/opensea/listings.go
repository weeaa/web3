package opensea

import (
	"fmt"
	"github.com/foundVanting/opensea-stream-go/entity"
	"github.com/mitchellh/mapstructure"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/pkg/logger"
	"math/big"
	"time"
)

func (c *Client) GetListings(client discord.Client, collections []string, floorBelow float64) {

	logger.LogStartup(moduleName + " Listings")

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

				var floorBelow20, double float64
				floorBelow20 = l.PriceInfo.Floor / floorBelow
				double = floorBelow20 * 2
				floorMinus20 := l.PriceInfo.Floor - double

				l.PriceInfo.PriceBefore, _ = l.PriceInfo.Price.Float64()

				if l.PriceInfo.PriceBefore <= floorMinus20 {
					if err = client.SendNotification(discord.Webhook{
						Username:  "OpenSea",
						AvatarUrl: client.AvatarImage,
						Embeds: []discord.Embed{
							{
								Title:       l.Collection,
								Description: fmt.Sprintf("%v just listed %v for `%2f %v` at <t:%v>.", l.Seller, l.Item, l.PriceInfo.PriceBefore, l.Symbol, time.Now().Unix()),
								Thumbnail: discord.EmbedThumbnail{
									Url: l.Image,
								},
								Url:       l.CollectionLink,
								Color:     client.Color,
								Timestamp: discord.GetTimestamp(),
								Footer: discord.EmbedFooter{
									Text:    client.FooterText,
									IconUrl: client.FooterImage,
								},
								Fields: []discord.EmbedFields{
									{
										Name:   "Collection Slug",
										Value:  slug,
										Inline: true,
									},
									{
										Name:   "Price (wei)",
										Value:  fmt.Sprint(l.Price),
										Inline: true,
									},
									{
										Name:   "Floor Price",
										Value:  fmt.Sprintf("%f ETH", l.PriceInfo.Floor),
										Inline: true,
									},
								},
							},
						},
					}, moduleName); err != nil {
						logger.LogError(moduleName, err)
					}
				}
			})
		}(collection)
	}
}
