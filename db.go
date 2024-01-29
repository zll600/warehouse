package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func generateUserID(numUsers int) []int {
	var users []int
	for i := 1; i <= numUsers; i++ {
		users = append(users, i)
	}
	return users
}

func generateRandomWeight() float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	weightProbabilities := map[float64]float64{
		2.0: 0.4, // 40% 的订单为 2kg
		4.0: 0.2, // 20% 的订单为 4kg
		6.0: 0.1, // 10% 的订单为 6kg
		8.0: 0.3, // 30% 的订单为 8kg
	}
	totalOrders := 100
	weights := make([]float64, 0, totalOrders)

	for weight, prob := range weightProbabilities {
		numOrders := int(prob * float64(totalOrders))
		for i := 0; i < numOrders; i++ {
			weights = append(weights, weight)
		}
	}

	idx := r.Intn(len(weights))
	return weights[idx]
}

func initdb() {
	os.Remove("./warehouse.db")

	db, err := sql.Open("sqlite3", "./warehouse.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	create table warehouse(id integer primary key, uid integer, weight double, created_at datetime default current_timestamp);
	create index idx_uid ON warehouse(uid);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into warehouse(uid, weight) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	users := generateUserID(10)
	for i := 0; i < 100; i++ {
		uid := users[rand.Intn(len(users))]
		weight := generateRandomWeight()
		log.Printf("insert data id: %d, uid: %d, weight: %f\n", i, uid, weight)
		_, err := stmt.Exec(uid, weight)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func printAllData(uid int) {
	db, err := sql.Open("sqlite3", "./warehouse.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, weight, created_at FROM warehouse WHERE uid = ?`, uid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	sum := 0
	for rows.Next() {
		var id int
		var weight float64
		var createAt string

		err = rows.Scan(&id, &weight, &createAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("id: %d, weight: %f, create_at: %s\n", id, weight, createAt)
		var res int
		res, err = calc(weight)
		if err != nil {
			log.Fatal(res)
		}
		sum += res
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("user %d total cost: %d\n", uid, sum)
}
