package ws

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestAsset_Auth(t *testing.T) {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	log.Info("miao")
	asset, _ := NewAsset()
	asset.Auth("a4382164-ed2htwf5tf-6d55e15e-701e5", "e7de9097-0adeb442-66b6f2d7-76752")

	asset.Loop()

}
