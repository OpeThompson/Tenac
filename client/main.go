package main

import (
    "context"
    "flag"
    "log"
    "time"

    pb "github.com/OpeThompson/tenac/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    addr := flag.String("addr", "localhost:50051", "server address")
    flag.Parse()

    conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("dial: %v", err)
    }
    defer conn.Close()

    c := pb.NewKeyValueStoreClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    if _, err := c.Set(ctx, &pb.SetRequest{Key: "name", Value: "Tenac"}); err != nil {
        log.Fatalf("Set: %v", err)
    }
    resp, err := c.Get(ctx, &pb.GetRequest{Key: "name"})
    if err != nil {
        log.Fatalf("Get: %v", err)
    }
    log.Printf("Value=%q Found=%v", resp.Value, resp.Found)
}
