package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "go-tokens/token_manager"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	write *bool
	read  *bool
	id    *string
	name  *string
	low   *uint64
	mid   *uint64
	high  *uint64
	port  *string
)

func init() {
	write = flag.Bool("write", false, "Write Command")
	read = flag.Bool("read", false, "Read Command")
	id = flag.String("id", "", "Token Id")
	name = flag.String("name", "", "Token Name")
	low = flag.Uint64("low", 0, "Low")
	mid = flag.Uint64("mid", 0, "Mid")
	high = flag.Uint64("high", 0, "High")
	port = flag.String("port", "", "Port")
}

type TokenList struct {
	Tokens []YamlToken `yaml:"tokens"`
}

type YamlToken struct {
	Id      string   `yaml:"id"`
	Writer  string   `yaml:"writer"`
	Readers []string `yaml:"readers"`
}

// type Config struct {
// 	Host string `yaml:"host"`
// 	Port string `yaml:"port"`
// }

func main() {
	flag.Parse()

	conn, err := grpc.Dial("localhost:"+*port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTokenManagerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if *write {
		if len(*id) > 0 && len(*name) > 0 && *low <= *mid && *mid < *high {
			r, err := c.WriteToken(ctx, &pb.WriteRequest{Id: *id, Name: *name, Low: *low, High: *high, Mid: *mid})
			if err == nil {
				log.Printf("\nWrite: %d", r.GetValue())
			} else {
				log.Println("Failed to Write:", err)
			}
		} else {
			log.Println("Please enter correct values: (id & name should not be empty) & (low <= mid && mid < high)")
		}
	} else if *read {
		if len(*id) > 0 {
			r, err := c.ReadToken(ctx, &pb.Id{Id: *id})
			if err == nil {
				log.Printf("\nRead: %d", r.GetValue())
			} else {
				log.Println("Failed to Read", err)
			}
		} else {
			log.Println("Id not provided")
		}
	}

}
