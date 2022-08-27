package utils

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	ts "go-tokens/token_sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SafeTokens struct {
	mu     sync.RWMutex
	tokens map[string]Token
}

type Domain struct {
	low  uint64
	mid  uint64
	high uint64
}

type State struct {
	patrial uint64
	final   uint64
}

type Token struct {
	name    string
	domain  Domain
	state   State
	writer  string
	readers []string
	ts      int64
}

type NodeList struct {
	Nodes []Config `yaml:"nodes"`
}

type Config struct {
	Id   string `yaml:"id"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type YamlToken struct {
	Id      string   `yaml:"id"`
	Writer  string   `yaml:"writer"`
	Readers []string `yaml:"readers"`
}

func InitMap() *SafeTokens {
	return &SafeTokens{tokens: make(map[string]Token)}
}

func (c *SafeTokens) AddTokens(yamlTokens []YamlToken) {
	for _, token := range yamlTokens {
		c.tokens[token.Id] = Token{
			name:    "",
			domain:  Domain{},
			state:   State{},
			writer:  token.Writer,
			readers: token.Readers,
		}
	}
}

func Hash(name string, nonce uint64) uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s %d", name, nonce)))
	return binary.BigEndian.Uint64(hasher.Sum(nil))
}

func ArgMinX(name string, start uint64, end uint64) uint64 {
	var hash = Hash(name, start)
	var x = start
	for i := start + 1; i < end; i++ {
		if Hash(name, i) < hash {
			hash = Hash(name, i)
			x = i
		}
	}
	return x
}

func (c *SafeTokens) Write(id string, name string, low uint64, high uint64, mid uint64) (uint64, error) {
	_, val := c.tokens[id]
	if val {
		c.mu.Lock()
		token := c.tokens[id]
		token.name = name
		token.domain = Domain{
			low:  low,
			high: high,
			mid:  mid,
		}
		token.state = State{
			patrial: ArgMinX(name, low, mid),
			final:   0,
		}
		token.ts = time.Now().Unix()
		c.tokens[id] = token
		c.mu.Unlock()
		return token.state.patrial, nil
	} else {
		return 0, errors.New("Token with given id:" + id + " does not exists")
	}
}

func (c *SafeTokens) Read(id string, nodeList *NodeList, nodeId string) (uint64, error) {
	_, val := c.tokens[id]
	if val {
		c.mu.Lock()
		token := c.tokens[id]
		if token.ts == 0 {
			wToken := c.GetWriterTokenDetails(id)
			token.name = wToken.Name
			token.domain = Domain{
				low:  wToken.Domain.Low,
				high: wToken.Domain.High,
				mid:  wToken.Domain.Mid,
			}
			token.state = State{
				patrial: wToken.State.Partial,
				final:   wToken.State.Final,
			}
			token.ts = int64(wToken.Ts)
			var x = ArgMinX(token.name, token.domain.mid, token.domain.high)
			if x < token.state.patrial {
				token.state.final = x
			} else {
				token.state.final = token.state.patrial
			}
			c.tokens[id] = token
			c.mu.Unlock()
			return token.state.final, nil
		}
		wts := c.GetWriterTimestamp(id)
		if wts == 0 {
			c.mu.Unlock()
			return 0, errors.New("Cannot get timestamp from writer for:" + id)
		}
		if wts == token.ts {
			var x = ArgMinX(token.name, token.domain.mid, token.domain.high)
			if x < token.state.patrial {
				token.state.final = x
			} else {
				token.state.final = token.state.patrial
			}
			c.tokens[id] = token
			c.mu.Unlock()
			return token.state.final, nil
		} else if token.ts < wts {
			config := &Config{}
			for _, node := range nodeList.Nodes {
				if node.Id == nodeId {
					config.Id = node.Id
					config.Host = node.Host
					config.Port = node.Port
				}
			}
			isMajority := c.MajorityWriteBack(id, config.Host+config.Port)
			if isMajority {
				wToken := c.GetWriterTokenDetails(id)
				token.name = wToken.Name
				token.domain = Domain{
					low:  wToken.Domain.Low,
					high: wToken.Domain.High,
					mid:  wToken.Domain.Mid,
				}
				token.state = State{
					patrial: wToken.State.Partial,
					final:   wToken.State.Final,
				}
				token.ts = int64(wToken.Ts)
				var x = ArgMinX(token.name, token.domain.mid, token.domain.high)
				if x < token.state.patrial {
					token.state.final = x
				} else {
					token.state.final = token.state.patrial
				}
				c.tokens[id] = token
				c.mu.Unlock()
				return token.state.final, nil
			} else {
				c.mu.Unlock()
				return token.state.final, errors.New("Cannot form majority")
			}
		}
		c.mu.Unlock()
		return 0, errors.New("No conditions matched")
	} else {
		return 0, errors.New("Token with given id:" + id + " does not exists")
	}
}

func (c *SafeTokens) Update(in *ts.TokenBroadcast) {
	_, val := c.tokens[in.Id]
	if val {
		c.mu.Lock()
		token := c.tokens[in.Id]
		token.name = in.Name
		token.domain = Domain{
			low:  in.Domain.Low,
			high: in.Domain.High,
			mid:  in.Domain.Mid,
		}
		token.state = State{
			patrial: in.State.Partial,
			final:   in.State.Final,
		}
		token.ts = int64(in.Ts)
		c.tokens[in.Id] = token
		c.mu.Unlock()
	}
}

func (c *SafeTokens) SendRequestsForReplication(id string) {
	token := c.tokens[id]
	for _, addr := range token.readers {
		log.Println(addr)
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Could not connect: %v", err)
		}
		defer conn.Close()
		c := ts.NewTokenSyncClient(conn)
		bctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		domain := &ts.Domain{Low: token.domain.low, Mid: token.domain.mid, High: token.domain.high}
		state := &ts.State{Partial: token.state.patrial, Final: token.state.final}
		r, err := c.Replicate(bctx, &ts.TokenBroadcast{Id: id, Name: token.name, Ts: uint64(token.ts), Domain: domain, State: state})
		if err == nil {
			log.Println("\nSync Resp:", r.GetMsg())
		} else {
			log.Println("Failed to Sync:", err)
		}
	}
}

func (c *SafeTokens) GetWriterTimestamp(id string) int64 {
	token := c.tokens[id]
	conn, err := grpc.Dial(token.writer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	w := ts.NewTokenSyncClient(conn)
	bctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	r, err := w.GetTimeStamp(bctx, &ts.TokenId{Id: id})
	if err == nil {
		return int64(r.Ts)
	} else {
		log.Println("Failed to get timestamp:", err)
		return 0
	}
}

func (c *SafeTokens) GetWriterTokenDetails(id string) *ts.TokenBroadcast {
	token := c.tokens[id]
	conn, err := grpc.Dial(token.writer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	w := ts.NewTokenSyncClient(conn)
	bctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	r, err := w.GetTokenDetails(bctx, &ts.TokenId{Id: id})
	if err == nil {
		return r
	} else {
		log.Println("Failed to get timestamp:", err)
		return &ts.TokenBroadcast{}
	}
}

func (c *SafeTokens) MajorityWriteBack(id string, self string) bool {
	token := c.tokens[id]
	timeStamps := make([]int64, 0)
	for _, addr := range token.readers {
		if addr != self {
			conn, err := grpc.Dial(token.writer, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Could not connect: %v", err)
			}
			defer conn.Close()
			w := ts.NewTokenSyncClient(conn)
			bctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			r, err := w.GetTimeStamp(bctx, &ts.TokenId{Id: id})
			if err == nil {
				timeStamps = append(timeStamps, int64(r.Ts))
			} else {
				log.Println("Failed to get timestamp:", err)
			}
		}
	}
	size := len(token.readers)
	counter := 1
	ts := timeStamps[0]
	for i, t := range timeStamps {
		if i != 0 {
			if t == ts {
				counter += 1
			}
		}
	}
	if counter > size/2 {
		return true
	} else {
		return false
	}
}

func (c *SafeTokens) GetTokenTimeStamp(id string) int64 {
	return c.tokens[id].ts
}

func (c *SafeTokens) GetName(id string) string {
	return c.tokens[id].name
}
func (c *SafeTokens) GetLow(id string) uint64 {
	return c.tokens[id].domain.low
}
func (c *SafeTokens) GetMid(id string) uint64 {
	return c.tokens[id].domain.mid
}
func (c *SafeTokens) GetHigh(id string) uint64 {
	return c.tokens[id].domain.high
}
func (c *SafeTokens) GetPartial(id string) uint64 {
	return c.tokens[id].state.patrial
}
func (c *SafeTokens) GetFinal(id string) uint64 {
	return c.tokens[id].state.final
}

func (c *SafeTokens) GetTokens() *map[string]Token {
	return &c.tokens
}

func (c *SafeTokens) ListTokens() {
	log.Println("------------------------ Tokens ------------------------")
	c.mu.Lock()
	for k, v := range c.tokens {
		log.Println(k, v)
	}
	c.mu.Unlock()
	log.Println("------------------------- End --------------------------")
}
