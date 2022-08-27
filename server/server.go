package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"

	"go-tokens/server/utils"
	pb "go-tokens/token_manager"
	ts "go-tokens/token_sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"gopkg.in/yaml.v2"
)

var (
	t        *utils.SafeTokens
	nodeId   *string
	file     *string
	nodeList *utils.NodeList
)

type tokenServer struct {
	pb.UnimplementedTokenManagerServer
}

type syncServer struct {
	ts.UnimplementedTokenSyncServer
}

type TokenList struct {
	Tokens []utils.YamlToken `yaml:"tokens"`
}

func init() {
	nodeId = flag.String("node", "", "Node Id")
	file = flag.String("config", "../config.yaml", "Config file")
}

func (s *tokenServer) WriteToken(ctx context.Context, in *pb.WriteRequest) (*pb.Partial, error) {
	p, _ := peer.FromContext(ctx)
	log.Println(in.Id + " WRITE " + p.Addr.String())
	res, err := t.Write(in.Id, in.Name, in.Low, in.High, in.Mid)
	if err != nil {
		log.Println(err)
		return &pb.Partial{}, err
	}
	t.SendRequestsForReplication(in.Id)
	t.ListTokens()
	return &pb.Partial{Value: res}, nil
}

func (s *tokenServer) ReadToken(ctx context.Context, in *pb.Id) (*pb.Final, error) {
	p, _ := peer.FromContext(ctx)
	log.Println(in.Id + " READ " + p.Addr.String())
	res, err := t.Read(in.Id, nodeList, *nodeId)
	if err != nil {
		log.Println(err)
		return &pb.Final{}, err
	}
	t.ListTokens()
	return &pb.Final{Value: res}, nil
}

func (s *syncServer) Replicate(ctx context.Context, in *ts.TokenBroadcast) (*ts.Response, error) {
	log.Println(in.Id + " SYNC")
	t.Update(in)
	t.ListTokens()
	return &ts.Response{Msg: "Token with Id: " + in.Id + "is Replicated Successfully"}, nil
}

func (s *syncServer) GetTimeStamp(ctx context.Context, in *ts.TokenId) (*ts.TimeStamp, error) {
	log.Println(in.Id + " TS")
	time := t.GetTokenTimeStamp(in.Id)
	return &ts.TimeStamp{Ts: uint64(time)}, nil
}

func (s *syncServer) GetTokenDetails(ctx context.Context, in *ts.TokenId) (*ts.TokenBroadcast, error) {
	log.Println(in.Id + " TS")
	domain := &ts.Domain{Low: t.GetLow(in.Id), Mid: t.GetMid(in.Id), High: t.GetHigh(in.Id)}
	state := &ts.State{Partial: t.GetPartial(in.Id), Final: t.GetFinal(in.Id)}
	return &ts.TokenBroadcast{Id: in.Id, Name: t.GetName(in.Id), Ts: uint64(t.GetTokenTimeStamp(in.Id)), Domain: domain, State: state}, nil

}

// host string, port string, tokenIds []string
func main() {

	flag.Parse()
	yamlBytes, err := os.ReadFile(*file)
	if err != nil {
		log.Fatal((err))
	}
	tokenList := &TokenList{}
	err = yaml.Unmarshal(yamlBytes, tokenList)
	if err != nil {
		log.Fatal((err))
	}
	nodeList := &utils.NodeList{}
	err = yaml.Unmarshal(yamlBytes, nodeList)
	if err != nil {
		log.Fatal((err))
	}

	config := &utils.Config{}
	for _, node := range nodeList.Nodes {
		if node.Id == *nodeId {
			config.Id = node.Id
			config.Host = node.Host
			config.Port = node.Port
		}
	}

	t = utils.InitMap()
	t.AddTokens(tokenList.Tokens)

	lis, err := net.Listen("tcp", config.Host+":"+config.Port)
	if err != nil {
		log.Fatalf("\nFailed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTokenManagerServer(s, &tokenServer{})
	ts.RegisterTokenSyncServer(s, &syncServer{})

	log.Printf("\nServer listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("\nFailed to serve: %v", err)
	}

}

// func main() {
// 	flag.Parse()
// 	yamlBytes, err := os.ReadFile(*file)
// 	if err != nil {
// 		log.Fatal((err))
// 	}
// 	config := &Config{}
// 	err = yaml.Unmarshal(yamlBytes, config)
// 	if err != nil {
// 		log.Fatal((err))
// 	}

// }
