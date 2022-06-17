package pando

import (
	"context"
	crand "crypto/rand"
	"encoding/base64"
	golegs "github.com/filecoin-project/go-legs"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/sync"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipld/go-ipld-prime/traversal/selector"
	consumerSdk "github.com/kenlabs/pando/sdk/pkg/consumer"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

const (
	//privateKeyStr  = "CAESQAycIStrQXBoxgf2pEazDLoZbL8WCLX5GIb69dl4x2mJMpukCAPbzq1URPtKen4Bpxfz9et2exWhfAfZ/RG30ts="
	pandoAddr      = "/ip4/52.14.211.248/tcp/9013"
	pandoApi       = "http://52.14.211.248:9011"
	pandoPeerID    = "12D3KooWNU48MUrPEoYh77k99RbskgftfmSm3CdkonijcM5VehS9"
	providerPeerID = "12D3KooWNnK4gnNKmh6JUzRb34RqNcBahN5B8v18DsMxQ8mCqw81"
)

const (
	connectTimeout = time.Minute
	syncTimeout    = 10 * time.Minute
	syncDepth      = 50
)

type Count struct {
	Success int64
	Failed  int64
}

type DealTaskStats struct {
	Retrieval Count
	Storage   Count
}

func loadOrInitPeerKey(kf string) (crypto.PrivKey, error) {
	data, err := ioutil.ReadFile(kf)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		k, _, err := crypto.GenerateEd25519Key(crand.Reader)
		if err != nil {
			return nil, err
		}

		data, err := crypto.MarshalPrivateKey(k)
		if err != nil {
			return nil, err
		}

		if err := ioutil.WriteFile(kf, data, 0600); err != nil {
			return nil, err
		}

		return k, nil
	}
	return crypto.UnmarshalPrivateKey(data)
}

func encodePrivateKey(privateKey crypto.PrivKey) (string, error) {
	privateKeyBytes, err := crypto.MarshalPrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	privateKeyStr := base64.StdEncoding.EncodeToString(privateKeyBytes)
	if err != nil {
		return "", err
	}
	return privateKeyStr, nil
}

func GetMinerDealTaskStats(keyFile string) map[string]*DealTaskStats {
	privateKey, err := loadOrInitPeerKey(keyFile)
	if err != nil {
		logger.WithError(err).Fatal("call loadOrInitPeerKey failed")
	}
	peerID, err := peer.IDFromPrivateKey(privateKey)
	if err != nil {
		logger.WithError(err).Fatal("call IDFromPrivateKey failed")
	}
	logger.Debugf("consumer peerID: %v", peerID.String())

	privateKeyStr, err := encodePrivateKey(privateKey)
	if err != nil {
		logger.WithError(err).Fatal("call encodePrivateKey failed")
	}

	ch := make(chan *FinishedTask, 10000)
	rtMap := make(map[string]*DealTaskStats)

	go func() {
		// init store (mem for example), use custom link system for calculating
		ds := datastore.NewMapDatastore()
		mds := sync.MutexWrap(ds)
		bs := blockstore.NewBlockstore(mds)
		lsys := MkLinkSystem(bs, ch)
		defer close(ch)

		consumer, err := consumerSdk.NewDAGConsumer(privateKeyStr, pandoApi, connectTimeout, &lsys, syncTimeout)
		if err != nil {
			logger.WithError(err).Error("call consumerSdk.NewDAGConsumer failed")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), syncTimeout)
		defer cancel()
		go func() {
			select {
			case <-ctx.Done():
				err := consumer.Close()
				if err != nil {
					logger.WithError(err).Error("close DAGConsumer failed")
				}
			}
		}()

		sel := golegs.ExploreRecursiveWithStopNode(selector.RecursionLimitDepth(syncDepth), nil, nil)

		err = consumer.Start(pandoAddr, pandoPeerID, providerPeerID, sel)
		if err != nil {
			logger.WithError(err).Error("call consumer.Start failed")
		}
	}()

	count := 0
	for t := range ch {
		count++
		if t.RetrievalTask != nil {
			rt, ok := rtMap[t.RetrievalTask.Miner]
			if !ok {
				rt = &DealTaskStats{}
				rtMap[t.RetrievalTask.Miner] = rt
			}
			if t.ErrorMessage == nil {
				rt.Retrieval.Success += 1
			} else {
				rt.Retrieval.Failed += 1
			}
		}
		if t.StorageTask != nil {
			rt, ok := rtMap[t.StorageTask.Miner]
			if !ok {
				rt = &DealTaskStats{}
				rtMap[t.StorageTask.Miner] = rt
			}
			if t.ErrorMessage == nil {
				rt.Storage.Success += 1
			} else {
				rt.Storage.Failed += 1
			}
		}
	}
	logger.Debugf("received %d tasks", count)
	return rtMap
}
