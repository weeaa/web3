package indexer

import (
	"context"
	"encoding/json"
	"errors"
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
	"os"
	"strings"
	"time"
)

const indexer = "Friend Tech Indexer"

func New(db *db.DB, proxyFilePath string, rotateEachReq bool, delay time.Duration, filePath string) (*Indexer, error) {
	if !strings.Contains(filePath, "json") {
		s := strings.SplitAfter(filePath, ".")
		return nil, fmt.Errorf("expected json formatted file, got %s", s[len(s)-1])
	}

	f, err := files.ReadJSON[map[string]uint](filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			files.CreateFile(filePath)
			if err = files.WriteJSON(filePath, map[string]int{"id": 11}); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	if f["id"] < 11 {
		return nil, fmt.Errorf("[%s] id can't be below 11 (got %d)", filePath, f["id"])
	}

	proxyList, err := tls.ReadProxyFile(proxyFilePath)
	if err != nil {
		return nil, err
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

		if err = s.DB.Indexer.InsertUser(&models.FriendTechIndexer{UserID: fmt.Sprint(u.Id), BaseAddress: u.Address, TwitterUsername: u.TwitterUsername}, context.Background()); err != nil {
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
