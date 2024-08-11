package main

import (
	"02.08.2024-L0/database"
	"02.08.2024-L0/models"
	"02.08.2024-L0/services"
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
	OrdersCache map[string]*models.Order
}

func main() {
	err, appStash := newAppStash()
	if err != nil {
		log.Fatal(err)
	}
	defer appStash.closeAllConnections()

	err = database.StartupCacheFromDB(appStash.DbConn, appStash.OrdersCache)
	if err != nil {
		log.Fatal(err)
	}

	err = services.SubscribeNATS(appStash.NatsConn, appStash.getData)
	if err != nil {
		log.Fatal(err)
	}

	services.StartupServ(appStash.OrdersCache)
}

func (appStash *AppStash) getData(m *stan.Msg) {
	var order models.Order
	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		fmt.Printf("%s, skipping \n", err)
		return
	}

	err = database.WriteToDB(appStash.DbConn, order)
	if err != nil {
		log.Fatal(err)
	}
}

func newAppStash() (error, *AppStash) {
	err, db := database.ConnectDB()
	if err != nil {
		return err, nil
	}

	err, sc := services.ConnectNats(clientID)
	if err != nil {
		return err, nil
	}

	err, uidsCount := database.GetUIDsCount(db)
	if err != nil {
		log.Fatal(err)
	}
	ordersCache := make(map[string]*models.Order, uidsCount)

	return nil, &AppStash{DbConn: db, NatsConn: sc, OrdersCache: ordersCache}
}

func (appStash *AppStash) closeAllConnections() {
	appStash.DbConn.Close()
	appStash.NatsConn.Close()
}
