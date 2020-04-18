# skyway-gateway Go client

skyway-gateway を利用するためのクライアントをGoで書いたサンプルです。


## Requirements

* Docker
* docker-compose
* gstreamer

## Get started

1. skyway-gateway および skyway-js-sdkを用いたサンプルアプリを立ち上げる

```
$ docker-compose up -d
```

2. サンプルアプリケーションからCallされるGWクライアントを起動する

```
$ cd cmd/callee
$ APIKEY=YOUR_API_KEY DOMAIN=YOUR_GW_DOMAIN PORT=YOUR_GW_PORT go run m
ain.go
```

ここでは `PEER_ID` を指定することもできる

```
$ APIKEY=YOUR_API_KEY DOMAIN=YOUR_GW_DOMAIN PORT=YOUR_GW_PORT PEER_ID=hoge go run m
ain.go
```

3. カメラの映像を擬似したストリームを流す

```

```

4. ブラウザから `2.` で指定したPeerIDにCallする
