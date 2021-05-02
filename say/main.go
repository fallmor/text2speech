package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/Sirupsen/logrus"
	pb "github.com/fallmor/say-gprc/api"
	"google.golang.org/grpc"
)

func main() {
	backend := flag.String("b", "localhost:8080", "address of Say Backend")
	output := flag.String("o", "output.wav", "output where the audio will be generated")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("usage:\n\t%s \"text to speech\"\n", os.Args[0])
		os.Exit(1)
	}
	conn, err := grpc.Dial(*backend, grpc.WithInsecure())
	if err != nil {
		logrus.Fatalf("Could not dial the backend %s: %v", *backend, err)
	}
	defer conn.Close()
	client := pb.NewTextToSpeechClient(conn)
	text := &pb.Text{Text: flag.Arg(0)}
	resp, err := client.Say(context.Background(), text)
	if err != nil {
		log.Fatalf("could not read text %s: %v", text.Text, err)
	}
	if err := ioutil.WriteFile(*output, resp.Audio, 0666); err != nil {
		logrus.Fatalf("Could not write the file to %s: %v", *output, err)
	}
}
