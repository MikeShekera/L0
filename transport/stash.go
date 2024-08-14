package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/MikeShekera/L0/database"
	"github.com/MikeShekera/L0/models"
	"github.com/nats-io/stan.go"
)

type AppStash struct {
	DbConn      *sql.DB
	NatsConn    stan.Conn
	OrdersCache map[string]models.Order
	UIDSCount   int64
}

func (appStash *AppStash) CloseAllConnections() {
	appStash.DbConn.Close()
	appStash.NatsConn.Close()
}

func NewAppStash(db *sql.DB, sc stan.Conn, clientID string) (error, *AppStash) {
	err, uidsCount := database.GetUIDsCount(db)
	if err != nil {
		return err, nil
	}
	ordersCache := make(map[string]models.Order, uidsCount)

	return nil, &AppStash{DbConn: db, NatsConn: sc, OrdersCache: ordersCache, UIDSCount: uidsCount}
}

func (appStash *AppStash) GetData(m *stan.Msg) {
	var order models.Order
	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		fmt.Printf("%s, skipping \n", err)
		return
	}

	if _, ok := appStash.OrdersCache[order.OrderUID]; ok {
		fmt.Println(fmt.Sprintf("Order with uid: %s already in DB", order.OrderUID))
		return
	}

	err = database.WriteToDB(appStash.DbConn, order)
	if err != nil {
		fmt.Printf("%s, skipping \n", err)
		return
	}
}
