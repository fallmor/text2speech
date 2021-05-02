package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	pb "github.com/fallmor/say-gprc/api"

	"github.com/Sirupsen/logrus"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen to: ")
	flag.Parse()

	logrus.Infof("listenning to port %d", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logrus.Fatalf("Could not listen to port %d: %v", *port, err)
	}
	s := grpc.NewServer()
	pb.RegisterTextToSpeechServer(s, server{})
	err = s.Serve(lis)
	if err != nil {
		logrus.Fatalf("Could not lserver: %v", err)
	}
}

type server struct{}

func (server) Say(ctx context.Context, text *pb.Text) (*pb.Speech, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, fmt.Errorf("could not create tmp file %s", err)
	}
	if f.Close(); err != nil {
		return nil, fmt.Errorf("couldn't close file %s: %v", f.Name(), err)
	}
	cmd := exec.Command("flite", "-t", text.Text, "-o", f.Name())
	if data, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("flite failed %s", data)
	}
	data, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("could not read file %v", err)
	}
	return &pb.Speech{Audio: data}, nil
}
