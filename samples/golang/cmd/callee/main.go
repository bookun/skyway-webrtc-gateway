package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/bookun/skyway-webrtc-gateway/samples/golang/client"
)

func run(APIKey, peerID, gwDomain string, gwPort int) error {
	gwClient := client.NewClient(APIKey, gwDomain, gwPort)

	token, err := gwClient.CreatPeer(peerID, "localhost")
	if err != nil {
		return err
	}

	videoID, _, videoPort, err := gwClient.CreateMedia(true)
	fmt.Printf("\x1b[31mVIDEO PORT%d\x1b[0m\n", videoPort)
	if err != nil {
		return err
	}

	audioID, _, _, err := gwClient.CreateMedia(false)
	if err != nil {
		return err
	}

	videoRTCPID, _, _, err := gwClient.CreateRTCP()
	if err != nil {
		return err
	}

	audioRTCPID, _, _, err := gwClient.CreateRTCP()
	if err != nil {
		return err
	}

	constraints := client.Constraints{
		Video:               true,
		VideoReceiveEnabled: true,
		Audio:               true,
		AudioReceiveEnabled: true,
		VideoParams: struct {
			BandWidth   int    `json:"band_width"`
			Codec       string `json:"codec"`
			MediaID     string `json:"media_id"`
			RTCPID      string `json:"rtcp_id"`
			PayloadType int    `json:"payload_type"`
		}{
			BandWidth:   1500,
			Codec:       "H264",
			MediaID:     videoID,
			RTCPID:      videoRTCPID,
			PayloadType: 100,
		},
		AudioParams: struct {
			BandWidth   int    `json:"band_width"`
			Codec       string `json:"codec"`
			MediaID     string `json:"media_id"`
			RTCPID      string `json:"rtcp_id"`
			PayloadType int    `json:"payload_type"`
		}{
			BandWidth:   1500,
			Codec:       "opus",
			MediaID:     audioID,
			RTCPID:      audioRTCPID,
			PayloadType: 111,
		},
	}

	redirect := client.MediaRedirect{
		Video: struct {
			IPv4 string `json:"ip_v4"`
			Port int    `json:"port"`
		}{
			IPv4: "127.0.0.1",
			Port: 20000,
		},
		Audio: struct {
			IPv4 string `json:"ip_v4"`
			Port int    `json:"port"`
		}{
			IPv4: "127.0.0.1",
			Port: 20001,
		},
		VideoRTCP: struct {
			IPv4 string `json:"ip_v4"`
			Port int    `json:"port"`
		}{
			IPv4: "127.0.0.1",
			Port: 20010,
		},
		AudioRTCP: struct {
			IPv4 string `json:"ip_v4"`
			Port int    `json:"port"`
		}{
			IPv4: "127.0.0.1",
			Port: 20011,
		},
	}

	mediaConnectionID, err := gwClient.WaitCall(peerID, token)
	if err != nil {
		return err
	}

	if err := gwClient.Answer(mediaConnectionID, constraints, redirect); err != nil {
		return err
	}

	//defer gwClient.MediaConnectionClose(mediaConnectionID)
	defer gwClient.PeerClose(peerID, token)

	// シグナル用のチャネル定義
	quit := make(chan os.Signal)

	// 受け取るシグナルを設定
	signal.Notify(quit, os.Interrupt)

	fmt.Println("listening")

	<-quit // ここでシグナルを受け取るまで以降の処理はされない
	return nil
}

const rs2Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randPeerID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = rs2Letters[rand.Intn(len(rs2Letters))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	APIKey := os.Getenv("APIKEY")
	peerID := os.Getenv("PEER_ID")
	gwDomain := os.Getenv("DOMAIN")
	gwPort := os.Getenv("PORT")

	if APIKey == "" {
		log.Fatal("APIKey env must be exported")
	}

	if peerID == "" {
		peerID = randPeerID(10)
	}
	fmt.Printf("\x1b[33mpeerID: %s\x1b[0m\n", peerID)

	if gwDomain == "" {
		log.Fatal("gwDomain env must be exported")
	}
	if gwPort == "" {
		log.Fatal("gwPort env must be exported")
	}

	gwPortNum, err := strconv.Atoi(gwPort)
	if err != nil {
		log.Fatal("gwPort env must be exported")
	}

	if err := run(APIKey, peerID, gwDomain, gwPortNum); err != nil {
		log.Fatal(err)
	}
}
