# skyway-gateway Go client sample

skyway-gateway を利用するためのクライアントをGoで書いたサンプルです。


## Requirements

* Docker
* docker-compose
* gstreamer

## Get started

1. .envを書き換える

```
$ cp .env.sample .env
```

その後 `.env` 内の APIKEY を Skyway から入手したAPIKEYに変更する

2. skyway-gateway および skyway-js-sdkを用いたサンプルアプリを立ち上げる

```
$ docker-compose up -d
```

3. サンプルアプリケーションを開く

http://localhost:8080/examples/p2p-media

4. サンプルアプリケーションからCallされるGWクライアントを起動する

```
$ docker-compose exec camera /bin/bash
$ ./callee
```

ここでは `PEER_ID` を指定することもできる

```
$ PEER_ID=hoge ./calee
```

実行後に赤文字で表示されるVIDEO_PORTは `6.` で使用する


5. `3.` で開いたブラウザから `4.` で指定したPeerIDにCallする

6. ストリームを流す

    `4.` とは別のターミナルで
```
$ docker-compose exec camera /bin/bash
$ gst-launch-1.0 -v videotestsrc ! video/x-raw,framerate=20/1 ! videoscale ! videoconvert ! x264enc ! rtph264pay pt=100 ! udpsink host=gw port=50001
```

ここでの `port=500001` は適宜 `4.` で表示されている VIDEO_PORTに書き換える


7. ブラウザにストリームが映る
