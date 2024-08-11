package main

import (
	"github.com/MikeShekera/L0/database"
	"github.com/MikeShekera/L0/services"
	"github.com/MikeShekera/L0/transport"
	_ "github.com/lib/pq"
	"log"
)

const (
	clientID = "main"
)

func main() {
	err, db := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	err, sc := services.ConnectNats(clientID)
	if err != nil {
		log.Fatal(err)
	}

	err, appStash := transport.NewAppStash(db, sc, clientID)
	if err != nil {
		log.Fatal(err)
	}
	defer appStash.CloseAllConnections()

	err = database.StartupCacheFromDB(appStash.DbConn, appStash.OrdersCache, appStash.UIDSCount)
	if err != nil {
		log.Fatal(err)
	}

	err = services.SubscribeNATS(appStash.NatsConn, appStash.GetData)
	if err != nil {
		log.Fatal(err)
	}

	services.StartupServ(appStash)
}
