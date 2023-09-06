package opensea

import (
	"fmt"
	"github.com/foundVanting/opensea-stream-go/entity"
	"github.com/mitchellh/mapstructure"
	"github.com/weeaa/nft/discord"
	"github.com/weeaa/nft/pkg/logger"
)

func (c *Client) MonitorSales(client discord.Client, collections []string) {

	logger.LogStartup(moduleNameSales)
	defer logger.LogShutDown(moduleNameSales)

	for _, collection := range collections {
		go func(slug string) {
			c.StreamClient.OnItemSold(fmt.Sprintf(slug), func(response any) {
				var ItemSoldEvent entity.ItemSoldEvent
				s := &Sale{}

				if err := mapstructure.Decode(response, &ItemSoldEvent); err != nil {
					logger.LogError(moduleNameSales, err)
				}

				s.Item = "[" + ItemSoldEvent.Payload.PayloadItemAndColl.Item.Metadata.Name + "]" + "(" + ItemSoldEvent.Payload.PayloadItemAndColl.Item.Permalink + ")"
				s.Seller = "[OpenSea Member]" + "(https://etherscan.io/address/" + ItemSoldEvent.Payload.Maker.Address + ")"
				s.Username = "[" + ItemSoldEvent.Payload.Maker.Address + "]" + "(" + "https://opensea.io/" + ItemSoldEvent.Payload.Maker.Address + ")"
				s.Collection = ItemSoldEvent.Payload.PayloadItemAndColl.Collection.Slug
				s.CollectionLink = "https://opensea.io/collection/" + ItemSoldEvent.Payload.PayloadItemAndColl.Collection.Slug
				s.Image = ItemSoldEvent.Payload.PayloadItemAndColl.Item.Metadata.ImageUrl

				//todo fill wh later
				if err := client.SendNotification(discord.Webhook{}, moduleNameSales); err != nil {
					logger.LogError(moduleNameSales, err)
				}
			})
		}(collection)
	}
}
