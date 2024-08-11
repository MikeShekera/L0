package main

import (
	"02.08.2024-L0/OrderJsonStructure"
	"02.08.2024-L0/ServerPart"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
)

const (
	clientID = "main"
)

type AppStash struct {
	DbConn      *sql.DB
	NatsConn    stan.Conn
	OrdersCache map[string]*OrderJsonStructure.Order
}

func main() {
	err, appStash := newAppStash()
	if err != nil {
		log.Fatal(err)
	}
	defer appStash.closeAllConnections()

	err = ServerPart.StartupCacheFromDB(appStash.DbConn, appStash.OrdersCache)
	if err != nil {
		log.Fatal(err)
	}

	err = ServerPart.SubscribeNATS(appStash.NatsConn, appStash.getData)
	if err != nil {
		log.Fatal(err)
	}

	ServerPart.StartupServ(appStash.OrdersCache)

	select {}
}

func (appStash *AppStash) getData(m *stan.Msg) {
	var order OrderJsonStructure.Order
	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		fmt.Printf("%s, skipping \n", err)
		return
	}

	err = ServerPart.WriteToDB(appStash.DbConn, order)
	if err != nil {
		log.Fatal(err)
	}
}

func newAppStash() (error, *AppStash) {
	err, db := ServerPart.ConnectDB()
	if err != nil {
		return err, nil
	}

	err, sc := ServerPart.ConnectNats(clientID)
	if err != nil {
		return err, nil
	}

	err, uidsCount := ServerPart.GetUIDsCount(db)
	if err != nil {
		log.Fatal(err)
	}
	ordersCache := make(map[string]*OrderJsonStructure.Order, uidsCount)

	return nil, &AppStash{DbConn: db, NatsConn: sc, OrdersCache: ordersCache}
}

func (appStash *AppStash) closeAllConnections() {
	appStash.DbConn.Close()
	appStash.NatsConn.Close()
}
