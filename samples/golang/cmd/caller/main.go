package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"

	"github.com/bookun/skyway-webrtc-gateway/samples/golang/client"
)

func run(APIKey, peerID, targetPeerID, gwDomain string, gwPort int) error {
	gwClient := client.NewClient(APIKey, gwDomain, gwPort)

	token, err := gwClient.CreatPeer(peerID, "localhost")
	fmt.Printf("token: %s\n", token)
	if err != nil {
		return err
	}

	videoID, _, videoPort, err := gwClient.CreateMedia(true)
	fmt.Printf("videoID: %s\n", videoID)
	fmt.Printf("video_port: %d\n", videoPort)
	if err != nil {
		return err
	}

	audioID, _, _, err := gwClient.CreateMedia(false)
	fmt.Printf("audioID: %s\n", audioID)
	if err != nil {
		return err
	}

	videoRTCPID, _, _, err := gwClient.CreateRTCP()
	fmt.Printf("videoRTCPID: %s\n", videoRTCPID)
	if err != nil {
		return err
	}

	audioRTCPID, _, _, err := gwClient.CreateRTCP()
	fmt.Printf("audioRTCPID: %s\n", audioRTCPID)
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

	mediaConnectionID, err := gwClient.Call(token, peerID, targetPeerID, constraints, redirect)
	if err != nil {
		return err
	}
	fmt.Printf("mediaConnectionID: %s", mediaConnectionID)
	if err := gwClient.WatchMediaConnection(mediaConnectionID); err != nil {
		log.Fatal(err)
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
	APIKey := os.Getenv("APIKey")
	peerID := os.Getenv("PEER_ID")
	targetPeerID := os.Getenv("TARGET_PEER_ID")
	gwDomain := os.Getenv("DOMAIN")
	gwPort := os.Getenv("PORT")

	if APIKey == "" {
		log.Fatal("APIKey env must be exported")
	}

	if peerID == "" {
		peerID = randPeerID(10)
	}
	fmt.Printf("peerID: %s\n", peerID)

	if targetPeerID == "" {
		log.Fatal("targetPeerID env must be exported")
	}

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

	if err := run(APIKey, peerID, targetPeerID, gwDomain, gwPortNum); err != nil {
		log.Fatal(err)
	}
}
