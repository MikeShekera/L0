package database

import (
	"database/sql"
	"fmt"
	"github.com/MikeShekera/L0/models"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "0000"
	dbname   = "L0"

	ordersInsertString = `
        INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        ON CONFLICT (order_uid) DO UPDATE SET 
            track_number = EXCLUDED.track_number,
            entry = EXCLUDED.entry,
            locale = EXCLUDED.locale,
            internal_signature = EXCLUDED.internal_signature,
            customer_id = EXCLUDED.customer_id,
            delivery_service = EXCLUDED.delivery_service,
            shardkey = EXCLUDED.shardkey,
            sm_id = EXCLUDED.sm_id,
            date_created = EXCLUDED.date_created,
            oof_shard = EXCLUDED.oof_shard;
    `

	deliveryInsertString = `
        INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (order_uid) DO UPDATE SET 
            name = EXCLUDED.name,
            phone = EXCLUDED.phone,
            zip = EXCLUDED.zip,
            city = EXCLUDED.city,
            address = EXCLUDED.address,
            region = EXCLUDED.region,
            email = EXCLUDED.email;
    `

	paymentInsertString = `
        INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        ON CONFLICT (order_uid) DO 
        UPDATE
        SET 
            transaction = EXCLUDED.transaction,
            request_id = EXCLUDED.request_id,
            currency = EXCLUDED.currency,
            provider = EXCLUDED.provider,
            amount = EXCLUDED.amount,
            payment_dt = EXCLUDED.payment_dt,
            bank = EXCLUDED.bank,
            delivery_cost = EXCLUDED.delivery_cost,
            goods_total = EXCLUDED.goods_total,
            custom_fee = EXCLUDED.custom_fee;
    `

	itemsInsertString = `
            INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
            ON CONFLICT (order_uid, chrt_id) DO 
            UPDATE
            SET 
                track_number = EXCLUDED.track_number,
                price = EXCLUDED.price,
                rid = EXCLUDED.rid,
                name = EXCLUDED.name,
                sale = EXCLUDED.sale,
                size = EXCLUDED.size,
                total_price = EXCLUDED.total_price,
                nm_id = EXCLUDED.nm_id,
                brand = EXCLUDED.brand,
                status = EXCLUDED.status;
        `

	receiveCustomIDString = `SELECT 
        o.order_uid,
        o.track_number,
        o.entry,
        o.locale,
        o.internal_signature,
        o.customer_id,
        o.delivery_service,
        o.shardkey,
        o.sm_id,
        o.date_created,
        o.oof_shard,
        
        d.name,
        d.phone,
        d.zip,
        d.city,
        d.address,
        d.region,
        d.email,
        
        p.transaction,
        p.request_id,
        p.currency,
        p.provider,
        p.amount,
        p.payment_dt,
        p.bank,
        p.delivery_cost,
        p.goods_total,
        p.custom_fee,
          
		i.chrt_id,
		i.track_number AS item_track_number,
		i.price,
		i.rid,
		i.name AS item_name,
		i.sale,
		i.size,
		i.total_price,
		i.nm_id,
		i.brand,
		i.status
	
	FROM 
		orders o
	LEFT JOIN 
		delivery d ON o.order_uid = d.order_uid
	LEFT JOIN 
		payment p ON o.order_uid = p.order_uid
	LEFT JOIN 
		items i ON o.order_uid = i.order_uid
	WHERE 
		o.order_uid = $1;`

	recieveItemsString = `SELECT * FROM items WHERE order_uid = $1`

	receiveAllDataString = `SELECT 
    o.order_uid,
    o.track_number,
    o.entry,
    o.locale,
    o.internal_signature,
    o.customer_id,
    o.delivery_service,
    o.shardkey,
    o.sm_id,
    o.date_created,
    o.oof_shard,
    
    d.name AS delivery_name,
    d.phone AS delivery_phone,
    d.zip AS delivery_zip,
    d.city AS delivery_city,
    d.address AS delivery_address,
    d.region AS delivery_region,
    d.email AS delivery_email,
    
    p.transaction AS payment_transaction,
    p.request_id AS payment_request_id,
    p.currency AS payment_currency,
    p.provider AS payment_provider,
    p.amount AS payment_amount,
    p.payment_dt AS payment_payment_dt,
    p.bank AS payment_bank,
    p.delivery_cost AS payment_delivery_cost,
    p.goods_total AS payment_goods_total,
    p.custom_fee AS payment_custom_fee
    
FROM 
    orders o
LEFT JOIN 
    delivery d ON o.order_uid = d.order_uid
LEFT JOIN 
    payment p ON o.order_uid = p.order_uid
`
)

func ConnectDB() (error, *sql.DB) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err, nil
	}
	return nil, db
}

func WriteToDB(db *sql.DB, order models.Order) error {
	_, err := db.Exec(
		ordersInsertString, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return err
	}

	delivery := order.Delivery
	_, err = db.Exec(
		deliveryInsertString, order.OrderUID, delivery.Name, delivery.Phone, delivery.Zip, delivery.City,
		delivery.Address, delivery.Region, delivery.Email,
	)
	if err != nil {
		return err
	}

	payment := order.Payment
	_, err = db.Exec(
		paymentInsertString, order.OrderUID, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider,
		payment.Amount, payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee,
	)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = db.Exec(
			itemsInsertString, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status,
		)
		if err != nil {
			return err
		}
	}
	fmt.Println("Insert Complete")
	return nil
}

func GetUIDsCount(db *sql.DB) (error, int64) {
	var ordersCount int64
	row := db.QueryRow(`Select count(*) from orders`)
	err := row.Scan(&ordersCount)
	if err != nil {
		return err, 0
	}
	return nil, ordersCount
}

func StartupCacheFromDB(db *sql.DB, cacheMap map[string]*models.Order) error {
	rows, err := db.Query(receiveAllDataString)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry,
			&order.Locale, &order.InternalSignature, &order.CustomerID,
			&order.DeliveryService, &order.ShardKey, &order.SmID,
			&order.DateCreated, &order.OofShard,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
			&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
			&order.Delivery.Email,
			&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
			&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
			&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee,
		)
		if err != nil {
			return err
		}
		cacheMap[order.OrderUID] = &order
		getItemsByOrderID(db, &order)
	}
	return nil
}

func getItemsByOrderID(db *sql.DB, data *models.Order) {
	rows, err := db.Query(recieveItemsString, data.OrderUID)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var item models.Item
		var unusedValue string
		err := rows.Scan(
			&unusedValue,
			&item.ChrtID, &item.TrackNumber, &item.Price,
			&item.RID, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID,
			&item.Brand, &item.Status,
		)
		if err != nil {
			log.Fatal(err)
		}
		data.Items = append(data.Items, item)
	}
}
