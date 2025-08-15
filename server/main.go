package main

import (
    "context"
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "net"
    "os"
    "time"

    pb "github.com/OpeThompson/tenac/proto"
    "github.com/OpeThompson/tenac/internal"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type Node struct {
    NodeID string `json:"node_id"`
    Host   string `json:"host"`
    Port   int    `json:"port"`
}

type server struct {
    pb.UnimplementedKeyValueStoreServer
    store  *internal.Store
    nodeID string
}

func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
    s.store.Set(req.Key, req.Value)
    return &pb.SetResponse{Success: true}, nil
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
    val, ok := s.store.Get(req.Key)
    return &pb.GetResponse{Value: val, Found: ok}, nil
}

func (s *server) HealthCheck(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
    log.Printf("[%s] HealthCheck from: %s", s.nodeID, req.FromNodeId)
    return &pb.HealthResponse{Healthy: true}, nil
}

func loadNodes() ([]Node, error) {
    f, err := os.Open("nodes.json")
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var nodes []Node
    if err := json.NewDecoder(f).Decode(&nodes); err != nil {
        return nil, err
    }
    return nodes, nil
}

func startHeartbeat(selfID string, nodes []Node) {
    for _, n := range nodes {
        if n.NodeID == selfID {
            continue
        }
        peer := n

        go func() {
            for {
                conn, err := grpc.Dial(
                    fmt.Sprintf("%s:%d", peer.Host, peer.Port),
                    grpc.WithTransportCredentials(insecure.NewCredentials()),
                )
                if err != nil {
                    log.Printf("[%s] Cannot dial %s: %v", selfID, peer.NodeID, err)
                    time.Sleep(5 * time.Second)
                    continue
                }

                client := pb.NewKeyValueStoreClient(conn)
                ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
                resp, err := client.HealthCheck(ctx, &pb.HealthRequest{FromNodeId: selfID})
                cancel()
                conn.Close()

                if err != nil || resp == nil || !resp.Healthy {
                    log.Printf("[%s] Peer %s is DOWN", selfID, peer.NodeID)
                } else {
                    log.Printf("[%s] Peer %s is ALIVE", selfID, peer.NodeID)
                }

                time.Sleep(5 * time.Second)
            }
        }()
    }
}

func main() {
    nodeID := flag.String("node_id", "node1", "ID of this node")
    port := flag.Int("port", 50051, "Port to run this node on")
    flag.Parse()

    nodes, err := loadNodes()
    if err != nil {
        log.Fatalf("Could not load nodes.json: %v", err)
    }

    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    s := grpc.NewServer()
    store := internal.NewStore()
    pb.RegisterKeyValueStoreServer(s, &server{store: store, nodeID: *nodeID})

    log.Printf("[%s] Starting Tenac node on :%d", *nodeID, *port)
    startHeartbeat(*nodeID, nodes)

    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
