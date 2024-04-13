package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fest/config"
	dao "fest/db"
	"fmt"

	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var (
	writeWait            time.Duration
	maxMessageSize       int64
	pingWait             time.Duration
	pongWait             time.Duration
	pingPeriod           time.Duration
	closeGracePeriod     time.Duration
	brPref               string
	wsAllowedOrigin      string
	isControlSessionOpen = false
	WsAsignConns         = make(map[string]string)
	db                   dao.Database
)

func init() {
	cfg, _ := config.LoadConfig("config.json")
	WsConfig := cfg.WsConfig

	writeWait = time.Duration(WsConfig.WriteWait) * time.Second               // Time allowed to write a message to the peer.
	maxMessageSize = WsConfig.MaxMessageSize                                  // Maximum message size allowed from peer.
	pingWait = time.Duration(WsConfig.PingWait) * time.Second                 // Time allowed to read the next pong message from the peer.
	pongWait = time.Duration(WsConfig.PongWait) * time.Second                 // Time allowed to read the next pong message from the peer.
	pingPeriod = time.Duration(WsConfig.PingPeriod) * time.Second             // Send pings to peer with this period. Must be less than pongWait.
	closeGracePeriod = time.Duration(WsConfig.CloseGracePeriod) * time.Second // Time to wait before force close on connection.
	brPref = WsConfig.BrPref
	wsAllowedOrigin = WsConfig.WsAllowedOrigin

	go hub.run()
	//log.Println("Config values: ", cfg)

	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbConnection.DBHost,
		cfg.DbConnection.DBPort,
		cfg.DbConnection.DBUser,
		cfg.DbConnection.DBPass,
		cfg.DbConnection.DBName,
	)

	db, _ = dao.NewDB(psqlconn)

}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		deviceOrigin := r.Header.Get("Origin")
		if wsAllowedOrigin == deviceOrigin {
			return true
		} else {
			log.Println("Origin not allowed:", deviceOrigin)
			return false
		}

	},
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		log.Println("{WS}:Error:", reason.Error())
	},
}

func ServeWs(w http.ResponseWriter, r *http.Request) {

	var ws *websocket.Conn
	var err error
	var token string
	var login string
	var tData []byte

	devId := r.Header.Get("DeviceId")
	if devId == "" {
		rw := r.URL.RawQuery
		qv := strings.Split(rw, "&")
		if qv[0] != "" {
			vv := strings.Split(qv[0], "=")
			devId = vv[1]
			if devId == "" || devId == "null" {
				suf, _ := HashPass(r.Header.Get("Sec-WebSocket-Key"))
				devId = brPref + strings.ToUpper(suf[12:18])
			}
		}
		if qv[1] != "" {
			vv := strings.Split(qv[1], "=")
			tData, err = base64.StdEncoding.DecodeString(vv[1] + strings.Repeat("=", len(vv)-2))
			if err != nil {
				log.Println("decode error:", err)
				return
			}
			ud := strings.Split(string(tData), "|")
			login = ud[0]
			token = ud[1]
		}
		log.Printf("user Data: %s ; %s", login, token)
		ok, dberr := db.AuthUser(login)
		if dberr != nil {
			log.Println("Auth Error:", err)
			return
		}
		if !ok {
			log.Println("Access Deny:")
			return
		}

		ok, dberr = db.UpdateSession(devId, token, login, true)
		if dberr != nil {
			log.Println("Auth Error:", err)
			return
		}
		if !ok {
			log.Println("Access Deny:")
			return
		}
	}

	log.Println("incoming WS request from: ", devId)

	swp := r.Header.Get("Sec-WebSocket-Protocol")
	if swp != "" {
		headers := http.Header{"Sec-Websocket-Protocol": {swp}}
		ws, err = upgrader.Upgrade(w, r, headers)
	} else {
		ws, err = upgrader.Upgrade(w, r, nil)
	}

	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	//defer ws.Close() !!!! Important

	conn := &Conn{send: make(chan []byte, 256), ws: ws, deviceId: devId}
	hub.register <- conn
	go conn.writePump(db)
	conn.readPump(db)

	//WsConnections[deviceId] = ws
	//if(ws != nil) {
	//	//ws.SetPingHandler(ping)
	//	//ws.SetPongHandler(pong)
	//	log.Println("Create new Ws Connection: succes, device: ", deviceId)
	//	go wsProcessor(ws, db, deviceId)
	//} else {
	//	log.Println("Ws Connectionfor device: ", deviceId, " not created.")
	//}

}

// readPump pumps messages from the websocket connection to the hub.
func (c *Conn) readPump(db dao.Database) {
	defer func() {
		hub.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	devId := c.deviceId

	for {
		if c != nil {
			mt, message, err := c.ws.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Printf("error: %v", err)
				}
				break
			}

			log.Printf("[WS]:recive: %s, type: %d", message, mt)
			msg := string(message)

			//if strings.Contains(msg, "WS: connection success") {
			//	if !sendMsg(c, "{\"action\":\"connect\",\"success\":true, \"device\":\""+devId+"\" }") {
			//		break
			//	} else {
			//		log.Println("{action:connect,success:true}", c.deviceId)
			//	}
			//}

			if strings.Contains(msg, "\"action\":\"connect\",\"success\":true") {
				//get Session data
				//var p ent.WsCommand
				//err := json.Unmarshal([]byte(msg), &p)
				//if err != nil {
				//	if !sendMsg(c, "{\"action\":\"createSession\",\"success\":false}") {
				//		break
				//	} else {
				//		log.Println("{action:createSession,success:false}", devId)
				//	}
				//	return
				//}
				//bData, err := json.Marshal(p.Data)
				//if err != nil {
				//	if !sendMsg(c, "{\"action\":\"createSession\",\"success\":false}") {
				//		break
				//	} else {
				//		log.Println("{action:createSession,success:false}", devId)
				//	}
				//	return
				//}
				//var sData ent.TildaSessionData
				//err = json.Unmarshal(bData, &sData)
				//if err != nil {
				//	if !sendMsg(c, "{\"action\":\"createSession\",\"success\":false}") {
				//		break
				//	} else {
				//		log.Println("{action:createSession,success:false}", devId)
				//	}
				//	return
				//}
				//Check Session
				//ok, err := db.CheckSession(c.deviceId, sData.Login, sData.Token)
				//if !ok && err != nil {
				//	if !sendMsg(c, "{\"action\":\"createSession\",\"success\":false}") {
				//		break
				//	} else {
				//		log.Println("{action:createSession,success:false}", devId)
				//	}
				//	return
				//} else {
				//	if ok { //Update Session
				//		ok, err = db.UpdateSession(devId)
				//		if !ok || !sendMsg(c, "{\"action\":\"updateSession\",\"success\":true, \"device\":\""+devId+"\" }") {
				//			break
				//		} else {
				//			log.Println("{action:updateSession,success:true}", devId)
				//		}
				//	} else { //Creeate session
				//		ok, err = db.CreateSession(devId, sData.Login, sData.Token)
				//		if !ok || err != nil {
				//			if !sendMsg(c, "{\"action\":\"createSession\",\"success\":false}") {
				//				break
				//			} else {
				//				log.Println("{action:createSession,success:false}", devId)
				//			}
				//			return
				//		}
				//		if !sendMsg(c, "{\"action\":\"createSession\",\"success\":true, \"device\":\""+devId+"\" }") {
				//			break
				//		} else {
				//			log.Println("{action:createSession,success:true}", devId)
				//		}
				//	}
				//}

				log.Println("{action:connect, success:true}", devId)

				if !sendMsg(c, "{\"action\":\"setDeviID\",\"device\":\""+devId+"\" }") {
					break
				} else {
					log.Println("{action:setDeviID,success:true}", devId)
				}

			}

			if strings.Contains(msg, "\"action\":\"getUserData\"") {

			}

			if strings.Contains(msg, "\"action\":\"getApartData\"") {

			}

			//###############################################################

		} else {
			log.Println("# Connection lost    #")
			hub.unregister <- c
			break
		}

	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Conn) writePump(db dao.Database) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		if c != nil {
			select {
			case message, ok := <-c.send:
				if !ok {
					// The hub closed the channel.
					c.write(websocket.CloseMessage, []byte{})
					return
				}

				c.ws.SetWriteDeadline(time.Now().Add(writeWait))
				w, err := c.ws.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				w.Write(message)
				// Add queued chat messages to the current websocket message.
				n := len(c.send)
				for i := 0; i < n; i++ {
					//w.Write(newline)
					w.Write(<-c.send)
				}

				if err := w.Close(); err != nil {
					return
				}
			case <-ticker.C:
				//log.Println("############################ Ping ##############################")

				if err := c.write(websocket.PingMessage, []byte{}); err != nil {
					log.Println("# Send ping to " + c.GetDeviceId() + " failed.  #")
					return
				} else {
					log.Println("# Send ping to " + c.GetDeviceId() + " success.  #")
				}

			}
		} else {
			hub.unregister <- c
			log.Println("# Connection lost    #")
		}
	}
}

func sendMsg(c *Conn, m string) bool {
	if c != nil {
		c.send <- []byte(m)
		log.Println("[WS]:send:", m, " succes")
		return true
	} else {
		log.Println("[WS]:send:", m, " failed")
		return false
	}

}

//func broadCastSend(msg string) bool {
//	hub.
//}

func unAssign(devId string) {
	for key, val := range WsAsignConns {
		if val == devId {
			conn := hub.getConnByDevId(key)
			if conn != nil {
				hub.getConnByDevId(key).send <- []byte("{\"action\":\"unassign\",\"device\":\"" + val + "\"}")
			}
			delete(WsAsignConns, key)
		}
	}
}

func SendMsgByDevId(devId string, msg string) {
	conn := hub.getConnByDevId(devId)
	sendMsg(conn, msg)
}

func f(tt string, begin int, end int) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Catch Panic in ", r)
		}
	}()
	ss := tt[begin:end]
	return ss, nil
}

func HashPass(p string) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(p))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
