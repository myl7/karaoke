package server

import (
	"crypto/rand"

	"github.com/myl7/karaoke/pkg/utils"
	"golang.org/x/crypto/nacl/box"
	"google.golang.org/protobuf/proto"
)

// TODO: Replace this dummy Client as real impl
type Client struct {
	c ClientConfig
}

type ClientConfig struct {
	PK, SK *[32]byte

	PCs      map[string]PConfig
	LayerIdx map[int][]string
}

func NewClient(c ClientConfig) *Client {
	if c.PK == nil || c.SK == nil {
		var err error
		c.PK, c.SK, err = box.GenerateKey(rand.Reader)
		if err != nil {
			panic(err)
		}
	}

	return &Client{c: c}
}

// TODO: Replace this dummy runRound as real impl
func (c *Client) RunRound() ([]byte, error) {
	deadDrop := "0"
	body := []byte("OK")
	o := &Onion{
		Body:     body,
		NextHop:  "",
		DeadDrop: deadDrop,
	}
	oB, err := proto.Marshal(o)
	if err != nil {
		return nil, err
	}

	route := make([]string, len(c.c.LayerIdx))
	for i, ids := range c.c.LayerIdx {
		id := ids[0]
		route[i] = id
	}

	for i := len(route) - 1; i >= 0; i-- {
		id := route[i]
		pC := c.c.PCs[id]
		newOB := utils.BoxEnc(oB, pC.PK)
		if i == 0 {
			oB = newOB
			break
		}
		o = &Onion{
			Body:     newOB,
			NextHop:  id,
			DeadDrop: "",
		}
		oB, err = proto.Marshal(o)
		if err != nil {
			return nil, err
		}
	}

	return oB, nil
}
