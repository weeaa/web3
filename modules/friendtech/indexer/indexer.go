package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/weeaa/nft/database/db"
	"github.com/weeaa/nft/database/models"
	"github.com/weeaa/nft/modules/friendtech/constants"
	"github.com/weeaa/nft/modules/friendtech/utils"
	"github.com/weeaa/nft/pkg/files"
	"github.com/weeaa/nft/pkg/logger"
	"github.com/weeaa/nft/pkg/tls"
	"io"
	"net/url"
	"time"
)

const indexer = "Friend Tech Indexer"

func New(db *db.DB, proxyFilePath string, rotateEachReq bool, delay time.Duration) (*Indexer, error) {
	f, err := files.ReadJSON[map[string]uint]("id.json")
	if err != nil {
		return nil, err
	}

	proxyList, err := tls.ReadProxyFile(proxyFilePath)
	if err != nil {
		return nil, err
	}

	if f["id"] < 11 {
		return nil, fmt.Errorf("UserID can't be below 11 (%d)", f["id"])
	}

	return &Indexer{
		userCounter:   f["id"],
		DB:            db,
		ProxyList:     proxyList,
		Delay:         delay,
		Client:        tls.New(tls.RandProxyFromList(proxyList)),
		RotateEachReq: rotateEachReq,
	}, nil
}

func (s *Indexer) StartIndexer() {
	logger.LogStartup(indexer)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.LogError(indexer, fmt.Errorf("panic recovered: %v", r))
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
			URL:    &url.URL{Scheme: "https", Host: constants.ProdBaseApi, Path: "/users/by-id/" + fmt.Sprint(s.userCounter)},
			Host:   constants.ProdBaseApi,
			Header: http.Header{
				"user-agent": {constants.IphoneUserAgent},
			},
		}

		resp, err := s.Client.Do(req)
		if err != nil {
			logger.LogError(indexer, err)
			time.Sleep(DefaultDelay)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusNotFound {
				// logger.LogError(indexer, fmt.Errorf("status not found for id: %d", s.userCounter))
				time.Sleep(s.Delay)
				continue
			} else if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusForbidden {
				tls.HandleRateLimit(s.Client, s.ProxyList, indexer)
				continue
			}
			time.Sleep(10 * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.LogError(indexer, err)
			continue
		}

		var u fren_utils.UserInformation
		if err = json.Unmarshal(body, &u); err != nil {
			logger.LogError(indexer, err)
			continue
		}

		if err = resp.Body.Close(); err != nil {
			logger.LogError(indexer, err)
			continue
		}

		if err = s.DB.Indexer.InsertUser(&models.FriendTechIndexer{UserID: fmt.Sprint(u.Id), BaseAddress: u.Address}, context.Background()); err != nil {
			logger.LogError(indexer, err)
		} else {
			logger.LogInfo(indexer, fmt.Sprintf("inserted %d | %s", u.Id, u.TwitterName))
		}

		s.userCounter++

		if err = files.WriteJSON("id.json", map[string]int{"id": u.Id}); err != nil {
			logger.LogError(indexer, err)
		}

		if s.RotateEachReq {
			if err = tls.RotateProxy(s.Client, s.ProxyList); err != nil {
				logger.LogError(indexer, err)
				continue
			}
		}

		time.Sleep(s.Delay)
	}
}

/*
	nitter, err := s.NitterClient.FetchNitter(u.TwitterName)
	if err != nil {
		logger.LogError(indexer, err)
		continue
	}

	followers, _ := strconv.Atoi(nitter.Followers)

	importance := fren_utils.AssertImportance(followers, fren_utils.Followers)
	if err != nil {
		logger.LogError(indexer, err)
		continue
	}
*/
