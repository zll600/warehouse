package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const MaxWeight = 100.0 // 最大重量 100kg

func generateUserID(numUsers int) []int {
	var users []int
	for i := 1; i <= numUsers; i++ {
		users = append(users, i)
	}
	return users
}

/*
1. 生成一个介于 0 到 1 之间的随机数 x。
2. 使用一个函数将 x 映射到一个重量值，这个函数应该能够保证较小的重量有更高的概率被选中。
3. 为了实现 1/W 分布，我们可以采用反函数的方法。这意味着我们可以使用类似于 weight = 100 * (1 - sqrt(x)) 的函数，其中 x 是 0 到 1 之间的随机数。
*/
func generateRandomWeight() float64 {
	x := rand.Float64()
	weight := 100 * (1 - math.Sqrt(x))
	return weight
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

	users := generateUserID(1000)
	for i := 0; i < 100000; i++ {
		uid := users[rand.Intn(len(users))]
		weight := generateRandomWeight()
		log.Printf("insert data id: %d, uid: %d, weight: %f\n", i, uid, weight)
		_, err = stmt.Exec(uid, weight)
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
