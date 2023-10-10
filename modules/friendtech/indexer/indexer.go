package indexer

import (
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/weeaa/nft/database/db"
	"github.com/weeaa/nft/modules/friendtech"
	"github.com/weeaa/nft/modules/friendtech/utils"
	"github.com/weeaa/nft/modules/twitter"
	"github.com/weeaa/nft/pkg/files"
	"github.com/weeaa/nft/pkg/logger"
	"github.com/weeaa/nft/pkg/tls"
	"io"
	"net/url"
	"strconv"
	"time"
)

const indexer = "Friend Tech Indexer"

func New(db *db.DB, proxyFilePath string, rotateEachReq bool) (*Indexer, error) {
	f, err := files.ReadJSON[map[string]int]("latestUserID.json")
	if err != nil {
		return nil, err
	}

	proxyList, err := tls.ReadProxyFile(proxyFilePath)
	if err != nil {
		return nil, err
	}

	return &Indexer{
		UserCounter:   f["id"],
		DB:            db,
		ProxyList:     proxyList,
		Client:        tls.New(tls.RandProxyFromList(proxyList)),
		RotateEachReq: rotateEachReq,
	}, nil
}

func (s *Indexer) StartIndexer() {
	logger.LogStartup(indexer)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.LogError(indexer, fmt.Errorf(""))
				s.StartIndexer()
			}
		}()

		s.Index()
	}()
}

// Index stores every user
// from the platform in a postgres database.
func (s *Indexer) Index() {
	for {
		req := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "https", Host: fren_utils.ProdBaseApi, Path: "/users/by-id/" + fmt.Sprint(s.UserCounter)},
			Host:   fren_utils.ProdBaseApi,
			Header: http.Header{},
		}

		resp, err := s.Client.Do(req)
		if err != nil {
			logger.LogError(indexer, err)
			time.Sleep(5 * time.Second)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusNotFound {
				logger.LogError(indexer, fmt.Errorf("status not found for id: %d", s.UserCounter))
				time.Sleep(4 * time.Second)
				continue
			} else if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusForbidden {
				tls.HandleRateLimit(s.Client, s.ProxyList, indexer)
				continue
			}
			logger.LogError(indexer, fmt.Errorf("status %s for id: %d", resp.Status, s.UserCounter))
			time.Sleep(10 * time.Second)
			continue
		}

		if s.RotateEachReq && resp.StatusCode != http.StatusForbidden {
			if err = tls.RotateProxy(s.Client, s.ProxyList); err != nil {
				logger.LogError(indexer, err)
				continue
			}
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.LogError(indexer, err)
			continue
		}

		var u friendtech.UserInformation
		if err = json.Unmarshal(body, &u); err != nil {
			logger.LogError(indexer, err)
			continue
		}

		if err = resp.Body.Close(); err != nil {
			logger.LogError(indexer, err)
			continue
		}

		var nitter twitter.NitterResponse
		var importance fren_utils.Importance

		{
			nitter, err = twitter.FetchNitter(u.TwitterName, s.Client)
			if err != nil {
				logger.LogError(indexer, err)
				continue
			}

			followers, _ := strconv.Atoi(nitter.Followers)

			importance = fren_utils.AssertImportance(followers, fren_utils.Followers)
			if err != nil {
				logger.LogError(indexer, err)
				continue
			}
		}

		_ = importance

		//todo move to API req rather as it will be ran on diff programs
		/*
			if err = s.DB.InsertFriendTechIndexer(&database.FriendTechIndexer{
				BaseAddress: u.Address,
				UserID:      fmt.Sprint(u.Id),
			}, context.Background()); err != nil {
				logger.LogError(indexer, err)
				continue
			}
		*/

		logger.LogInfo(indexer, fmt.Sprintf("inserted %d | %s", u.Id, u.TwitterName))

		s.UserCounter++

		if err = files.WriteJSON("latestUserID.json", map[string]int{"id": u.Id}); err != nil {
			logger.LogError(indexer, err)
		}

		time.Sleep(3 * time.Second)
	}
}
