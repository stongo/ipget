package main

import (
	"context"
	"log"
	"sync"

	iface "github.com/ipfs/interface-go-ipfs-core"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

func connect(ctx context.Context, ipfs iface.CoreAPI, peers []string) error {
	var wg sync.WaitGroup
	pinfos := make(map[peer.ID]*peerstore.PeerInfo, len(peers))
	for _, addrStr := range peers {
		addr, err := ma.NewMultiaddr(addrStr)
		if err != nil {
			return err
		}
		pii, err := peerstore.InfoFromP2pAddr(addr)
		if err != nil {
			return err
		}
		pi, ok := pinfos[pii.ID]
		if !ok {
			pi = &peerstore.PeerInfo{ID: pii.ID}
			pinfos[pi.ID] = pi
		}
		pi.Addrs = append(pi.Addrs, pii.Addrs...)
	}

	wg.Add(len(pinfos))
	for _, pi := range pinfos {
		go func(pi *peerstore.PeerInfo) {
			defer wg.Done()
			err := ipfs.Swarm().Connect(ctx, *pi)
			if err != nil {
				log.Printf("failed to connect to %s: %s", pi.ID, err)
			}
		}(pi)
	}
	wg.Wait()
	return nil
}
