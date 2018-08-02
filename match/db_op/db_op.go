package db_op

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
	"fmt"
)
type ACCOUNT struct {
	Id                 int      `db:"id"`
	Balance			   float64  `db:"balance"`
}

type ORDER struct {
	Trans_id int64
	Acc_id  int
	Sym_id  string
	Share   float64
	Price   float64
	Time    string
}
type QUERY struct {
	Trans_id int64
	Acc_id  int
	Sym_id  string
	Share   float64
	Price   float64
	Time    string
	Status  string
}
func main() {
	db, err := sqlx.Connect("mysql", "root:SJCzq@tcp(localhost:3306)/market")
	if err != nil {
		log.Fatalln(err)
	}
	Setup_db(db)
	fmt.Println("setup")
	Insert_account(db, 1, 1000) //args: int id, int balance
	Insert_account(db, 2, 1000)
	Insert_account(db, 3, 1000)
	Insert_account(db, 4, 1000)
	Insert_account(db, 5, 1000)
	Insert_symbol(db, "BTC") //args: string id
	Insert_symbol(db, "ATC")
	Insert_accsymbol(db, 1, "BTC", 500)
	Insert_accsymbol(db, 1, "BTC", 500)
	//Insert_accsymbol(db, 3, "BTC", 500)
	//Insert_accsymbol(db, 4, "BTC", 500)
	//Insert_accsymbol(db, 5, "BTC", 500)
	//Insert_accsymbol(db, 1, "ATC", 500)
	//Insert_accsymbol(db, 2, "ATC", 500)
	//Insert_accsymbol(db, 3, "ATC", 500)
	//Insert_accsymbol(db, 4, "ATC", 500)
	//Insert_accsymbol(db, 5, "ATC", 500)
	//Match(db, 1, "BTC", -1, 50)
	//Match(db, 1, "BTC", 1, 50)
	//Match(db, 1, "BTC", 1, 60)
	//Match(db, 2, "BTC", 1, 70)
	//Match(db, 3, "BTC", 1, 80)
	//Match(db, 4, "BTC", 1, 90)
	//Match(db, 1, "BTC", 1, 100)
	//Match(db, 1, "BTC", -5, 50)
	//Match(db, 1, "ATC", -1, 30)
	//Match(db, 1, "BTC", 5, 30)
	//Match(db, 1, "BTC", -1, 20)
	//Match(db, 1, "BTC", -1, 20)
	//Match(db, 1, "BTC", -1, 20)
	//Match(db, 1, "BTC", -1, 20)
	//Match(db, 1, "BTC", -1, 20)
	//Cancel_order(db, 1)
	//Cancel_order(db, 4)
	//Match(db, 2, "ATC", -1, 40)
	//Match(db, 3, "ATC", -1, 50)
	//Match(db, 4, "ATC", -1, 60)
	//Match(db, 4, "ATC", 2, 100)
	//Match(db, 5, "ATC", 4, 85)
}


func Setup_db(db *sqlx.DB){
	Drop_table(db)
	schema0 := `CREATE TABLE IF NOT EXISTS account(
		id INTEGER,
		balance DOUBLE CHECK (balance >= 0),
   		PRIMARY KEY ( id )
   		
	);`
	db.Exec(schema0)

	schema1 := `CREATE TABLE IF NOT EXISTS symbol(
		id VARCHAR(100),
		PRIMARY KEY ( id )
	);`
	db.Exec(schema1)

	schema2 := `CREATE TABLE IF NOT EXISTS acc_sym(
		acc_id INTEGER,
		sym_id VARCHAR(100),
		share DOUBLE,
		FOREIGN KEY (acc_id)
      		REFERENCES account(id)
      		ON UPDATE CASCADE ON DELETE RESTRICT,
      	FOREIGN KEY  (sym_id)
      		REFERENCES symbol(id)
      		ON UPDATE CASCADE ON DELETE RESTRICT
	);`
	db.Exec(schema2)

	schema3 := `CREATE TABLE IF NOT EXISTS pool(
		trans_id INT UNSIGNED AUTO_INCREMENT,
		acc_id INTEGER,
		sym_id VARCHAR(100),
		share DOUBLE,
		time TIMESTAMP,
		price DOUBLE,
		PRIMARY KEY (trans_id)
	);`
	db.Exec(schema3)

	schema4 := `CREATE TABLE IF NOT EXISTS query(
		trans_id INTEGER,
		acc_id INTEGER,
		sym_id VARCHAR(100),
		share DOUBLE,
		price DOUBLE,
		time TIMESTAMP,
		status VARCHAR(100)
	);`
	db.Exec(schema4)
}
func Drop_table(db *sqlx.DB){
	_, err := db.Exec("DROP TABLE IF EXISTS acc_sym")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS account")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS symbol")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS pool")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("DROP TABLE IF EXISTS query")
	if err != nil {
		panic(err)
	}
}
//insert account balance into DB
func Insert_account (db *sqlx.DB, id int, balance float64) bool{
	tx, err := db.Begin()
	defer tx.Commit()
	_, err= tx.Exec("INSERT INTO account (id, balance) VALUES (?, ?)", id, balance)
	if err!=nil{
		//log.Fatalln(err)
		//err = tx.Commit()
		return false
	}
	//err = tx.Commit()
	return true
	//accounts := []ACCOUNT{}
	//
	//db.Select(&accounts, "select * from account")
	//log.Println("users...")
	//log.Println(accounts)
}

//insert symbol into DB, 1 mean no such symbol, 2 mean insert error
func Insert_symbol (db *sqlx.DB, id string) bool{
	rows, err := db.Query("SELECT symbol.id FROM symbol WHERE symbol.id = ?;", id)
	if err!=nil{
		//log.Fatalln(err)
		return false
	}
	if rows.Next() == false {
		tx, err := db.Begin()
		_, err= tx.Exec("INSERT INTO symbol (id) VALUES (?)", id)
		if err!=nil{
			//log.Fatalln(err)
			fmt.Printf("%s error\r\n",id)
			return false
		}
		err = tx.Commit()
	}else{
		return false
	}
	return true
}

//insert account symbol info into DB
//@@@@@change
func Insert_accsymbol (db *sqlx.DB, acc_id int, sym_id string, share float64) bool{
	fmt.Println("start Insert_accsymbol function")
	defer fmt.Println("finish Insert_accsymbol function")
	var id int
	tx, err := db.Beginx()
	defer tx.Commit()
	tx.Get(&id,"SELECT acc_sym.acc_id FROM acc_sym WHERE acc_sym.acc_id = ? AND acc_sym.sym_id = ? FOR UPDATE;",acc_id, sym_id)
	//fmt.Println("pos1.1")
	fmt.Println(id)
	if id==0 {
		fmt.Println("IN IT")
		//fmt.Printf("INSERT INTO acc_sym (acc_id, sym_id, share) VALUES (%d, %s, %d)\r\n",acc_id, sym_id, share)
		//fmt.Println("pos1.2")
		_, err= tx.Exec("INSERT INTO acc_sym (acc_id, sym_id, share) VALUES (?, ?, ?)",acc_id, sym_id, share)
		//fmt.Println("pos1.3")
		if err!=nil{
			//fmt.Println("check err: ")
			//log.Fatalln(err)
			return false
		}
		//fmt.Println("pos1.6")
	} else{
		//fmt.Println("pos1.4")
		_, err = tx.Exec("UPDATE acc_sym SET acc_sym.share = acc_sym.share + ? WHERE acc_sym.sym_id = ? AND acc_sym.acc_id = ?;", share, sym_id,acc_id)
		//fmt.Println("pos1.5")
		//fmt.Println("pos1.7")
		if err != nil {
			log.Fatalln(err)
		}
	}
	return true
}

func Insert_accsymbol2 (db *sqlx.DB, tx *sqlx.Tx, acc_id int, sym_id string, share float64) bool{
	fmt.Println("start Insert_accsymbol2 function")
	defer fmt.Println("finish Insert_accsymbol2 function")
	var id int
	err := tx.Get(&id,"SELECT acc_sym.acc_id FROM acc_sym WHERE acc_sym.acc_id = ? AND acc_sym.sym_id = ? FOR UPDATE ;",acc_id, sym_id)

	//fmt.Println("pos1.1")
	if err!=nil {
		//fmt.Printf("INSERT INTO acc_sym (acc_id, sym_id, share) VALUES (%d, %s, %d)\r\n",acc_id, sym_id, share)
		//fmt.Println("pos1.2")
		_, err= tx.Exec("INSERT INTO acc_sym (acc_id, sym_id, share) VALUES (?, ?, ?)",acc_id, sym_id, share)
		//fmt.Println("pos1.3")
		if err!=nil{
			//fmt.Println("check err: ")
			//log.Fatalln(err)
			return false
		}
		//fmt.Println("pos1.6")

	} else{
		//fmt.Println("pos1.4")
		_, err = tx.Exec("UPDATE acc_sym SET acc_sym.share = acc_sym.share + ? WHERE acc_sym.sym_id = ? AND acc_sym.acc_id = ?;", share, sym_id,acc_id)
		//fmt.Println("pos1.5")
		//fmt.Println("pos1.7")
		if err != nil {
			log.Fatalln(err)
		}
	}
	return true
}

//insert transcation into BD
func Insert_pool (db *sqlx.DB, acc_id int, sym_id string, share float64, price float64) int64{
	fmt.Println("start Insert_pool")
	defer fmt.Println("finish Insert_pool")
	t := time.Now().Format("2006-01-02 15:04:05")
	tx, err := db.Begin()
	res, err := tx.Exec("INSERT INTO pool (acc_id, sym_id, share, price, time) VALUES (?, ?, ?, ?, ?)",acc_id, sym_id, share, price, t)
	err = tx.Commit()
	if err!=nil{
		log.Fatalln(err)
	}
	id, err := res.LastInsertId()
	return id
}

func Insert_query (db *sqlx.DB, tx *sqlx.Tx, trans_id int64, acc_id int, sym_id string, share float64, price float64, status string){
	t := time.Now().Format("2006-01-02 15:04:05")
	//tx, err = db.Begin()
	_, err := tx.Exec("INSERT INTO query (trans_id, acc_id, sym_id, share, price, status, time) VALUES (?, ?, ?, ?, ?, ?, ?)",trans_id, acc_id, sym_id, share, price, status, t)
	//err = tx.Commit()
	if err!=nil{
		log.Fatalln(err)
	}
}
//@@@@@
func Update_open (db *sqlx.DB, tx *sqlx.Tx, trans_id int64, minus_share float64){
	var share float64
	tx.Select(&share, "SELECT query.share FROM query WHERE trans_id = ? AND status = ? FOR UPDATE ;", trans_id, "open")
	_, err := tx.Exec("UPDATE query SET query.share = query.share - ? WHERE trans_id = ? AND status = ?;", minus_share, trans_id, "open")
	if err != nil {
		log.Fatalln(err)
	}
}
//@@@@@
func Update_exec (db *sqlx.DB, tx *sqlx.Tx,trans_id int64){
	var status string
	tx.Select(&status, "SELECT query.status FROM query WHERE trans_id = ? AND status = ? FOR UPDATE ;", trans_id, "open")
	_, err := tx.Exec("UPDATE query SET query.status = ? WHERE trans_id = ? AND status = ?;", "executed", trans_id, "open")
	if err != nil {
		log.Fatalln(err)
	}
}

func Cancel_order (db *sqlx.DB, trans_id int64)bool{
	fmt.Println("start Cancel")
	defer fmt.Println("finish cancel")
	orders := []ORDER{}
	tx, _ :=db.Beginx()
	//defer tx.Commit()
	tx.Select(&orders, "select * from pool WHERE trans_id = ?", trans_id)
	if len(orders)==0{
		tx.Commit()
		return false
	}
	if(orders[0].Share>0){
		Update_balance(db, tx, orders[0].Acc_id, float64(orders[0].Share)*orders[0].Price)
	}
	if(orders[0].Share<0){
		Update_share(db, tx, orders[0].Acc_id, orders[0].Sym_id, (-1)*orders[0].Share)
	}
	_, err := tx.Exec("DELETE FROM pool WHERE pool.trans_id = ?", trans_id)
	if err != nil {
		log.Fatalln(err)
	}
	var status string
	tx.Select(&status, "SELECT query.status FROM query WHERE trans_id = ? AND status = ? FOR UPDATE ;", trans_id, "open")
	_, err = tx.Exec("UPDATE query SET query.status = ? WHERE trans_id = ? AND status = ?;", "canceled", trans_id, "open")
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}
	return true
}

func Del_bytransid (db *sqlx.DB, tx *sqlx.Tx, trans_id int64){
	_, err := tx.Exec("DELETE FROM pool WHERE pool.trans_id = ?", trans_id)
	if err != nil {
		log.Fatalln(err)
	}
}

func Update_balance(db *sqlx.DB, tx *sqlx.Tx, acc_id int, add_amount float64) bool{
	fmt.Println("start Update_balance")
	defer fmt.Println("finish Update_balance")
	var balance float64
	err := tx.Get(&balance,"SELECT account.balance FROM account WHERE id = ? FOR UPDATE;",acc_id)
	if err!=nil{
		return false
		log.Fatalln(err)
	}
	//var balance float64
	//for rows.Next() {
	//	err = rows.Scan(&balance)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//}
	fmt.Println("test Update_balance pos1")
	if (balance > 0 && balance+add_amount>=0) || balance == 0 && add_amount > 0 {
		//tx, err := db.Begin()
		_, err = tx.Exec("UPDATE account SET account.balance = account.balance + ? WHERE id = ?;", add_amount, acc_id)
		//err = tx.Commit()
		fmt.Println("update by ", add_amount, "id ", acc_id)
		if err != nil {
			log.Fatalln(err)
		}
	} else{
		fmt.Println("Insufficient balance!")
		return false
	}
	return true
}
// func Update_balance2(db *sqlx.DB, acc_id int, add_amount float64) bool{
// 	fmt.Println("start Update_balance2")
// 	defer fmt.Println("finish Update_balance2")
// 	tx, _ :=db.Beginx()
// 	defer tx.Commit()
// 	var balance float64
// 	err := tx.Get(&balance,"SELECT account.balance FROM account WHERE id = ? FOR UPDATE;",acc_id)
// 	if err!=nil{
// 		log.Fatalln(err)
// 	}
// 	//var balance float64
// 	//for rows.Next() {
// 	//	err = rows.Scan(&balance)
// 	//	if err != nil {
// 	//		log.Fatalln(err)
// 	//	}
// 	//}
// 	fmt.Println("test Update_balance pos1")
// 	if (balance > 0 && balance+add_amount>=0) || balance == 0 && add_amount > 0 {
// 		//tx, err := db.Begin()
// 		_, err = tx.Exec("UPDATE account SET account.balance = account.balance + ? WHERE id = ?;", add_amount, acc_id)
// 		//err = tx.Commit()
// 		fmt.Println("update by ", add_amount, "id ", acc_id)
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 	} else{
// 		fmt.Println("Insufficient balance!")
// 		return false
// 	}
// 	return true
// }
func Update_share(db *sqlx.DB, tx *sqlx.Tx, acc_id int, sym_id string, add_share float64) bool{
	fmt.Println("start Update_share")
	defer fmt.Println("finish Update_share")
	var share float64
	err:= tx.Get(&share,"SELECT acc_sym.share FROM acc_sym WHERE acc_sym.acc_id = ? AND acc_sym.sym_id = ? FOR UPDATE;",acc_id, sym_id)
	if err!=nil {
		return false
		log.Fatalln(err)
	}
	//for rows.Next() {
	//	err = rows.Scan(&share)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//}

	if (share > 0 && share+add_share>=0)|| share == 0 && add_share > 0 {
		//tx, err := db.Begin()
		_, err = tx.Exec("UPDATE acc_sym SET acc_sym.share = acc_sym.share + ? WHERE acc_sym.acc_id = ? AND acc_sym.sym_id = ?;", add_share, acc_id, sym_id)
		//err = tx.Commit()

		if err != nil {
			log.Fatalln(err)
		}
	} else{
		fmt.Println("Insufficient share!")
		return false
	}
	return true
}
func Update_pool(db *sqlx.DB, tx *sqlx.Tx, trans_id int64, share float64){
	//tx, err := db.Begin()
	var share2 float64
	tx.Select(&share2, "SELECT pool.share FROM pool WHERE trans_id = ? FOR UPDATE ;", trans_id)
	_, err := tx.Exec("UPDATE pool SET pool.share = ? WHERE pool.trans_id = ?;", share, trans_id)
	//err = tx.Commit()

	if err != nil {
		log.Fatalln(err)
	}
}

func Queryby_transid(db *sqlx.DB, trans_id int64) []QUERY{
	fmt.Println("start query")
	defer fmt.Println("finish query")
	querys := []QUERY{}
	db.Select(&querys, "select * from query WHERE trans_id = ?", trans_id)
	fmt.Println(querys)
	return querys
}
func Match(db *sqlx.DB, acc_id int, sym_id string, share float64, price float64)int64{
	fmt.Println("start Match")
	defer fmt.Println("finish Match")
	var trans_id int64
	rows,_ := db.Query("SELECT * FROM account WHERE account.id = ?;", acc_id)
	if rows.Next() == false{
		trans_id = -1
		return trans_id
	}
	if share>0 {
		trans_id = Match_transbuy(db, acc_id, sym_id, share, price)
	} else {
		trans_id = Match_transsell(db, acc_id, sym_id, share, price)
	}
	return trans_id
}

func Match_transbuy(db *sqlx.DB, acc_id int, sym_id string, share float64, price float64)int64{
	fmt.Println("start Matchbuy")
	defer fmt.Println("finish Matchbuy")
	orders := []ORDER{}
	//fmt.Println("pos1")
	//fmt.Println("pos2")
	trans_id := Insert_pool(db, acc_id, sym_id, 0, price)

	tx, _ :=db.Beginx()
	defer tx.Commit()
	if Update_balance(db,tx,acc_id, (-1)*price*float64(share)) ==false{
		//deal with insufficient balance
		fmt.Println("deal with insufficient balance")
		return -2
	}
	Insert_accsymbol2(db,tx ,acc_id, sym_id,0)
	//fmt.Println("pos1")
	//Insert_accsymbol(db,acc_id,sym_id,0)
	//fmt.Println("pos2")
	//trans_id := Insert_pool(db, acc_id, sym_id, 0, price)
	tx.Select(&orders, "select * from pool WHERE sym_id = ? AND share <0 AND ? >= price ORDER BY price ASC, time ASC FOR UPDATE ", sym_id, price)
	for i:=0; share > 0 && i < len(orders) ;i++ {
		switch {
		case (-1)*share == orders[i].Share:
			//trans_id := Insert_pool(db, acc_id, sym_id, 0, price)
			//Del_bytransid(db, trans_id)
			Insert_query(db, tx, trans_id, acc_id, sym_id, share, orders[i].Price, "executed")
			Update_balance(db, tx, orders[i].Acc_id, orders[i].Price*float64((-1)*orders[i].Share))
			Update_balance(db, tx, acc_id, (price-orders[i].Price)*float64((-1)*orders[i].Share))
			Update_share(db, tx, acc_id, sym_id, share)
			Del_order(db, tx, orders[i])
			Update_exec (db, tx, orders[i].Trans_id)
			share = 0
		case share > (-1)*orders[i].Share:
			//trans_id := Insert_pool(db, acc_id, sym_id, 0, price)
			Insert_query(db, tx, trans_id, acc_id, sym_id, (-1)*orders[i].Share, orders[i].Price, "executed")
			share += orders[i].Share
			Update_balance(db, tx, orders[i].Acc_id, orders[i].Price*float64((-1)*orders[i].Share))
			Update_balance(db, tx, acc_id, (price-orders[i].Price)*float64((-1)*orders[i].Share))
			Update_share(db, tx, acc_id, sym_id, orders[i].Share)
			Del_order(db, tx, orders[i])
			Update_exec (db, tx, orders[i].Trans_id)
		case share < (-1)*orders[i].Share:
			//trans_id := Insert_pool(db, acc_id, sym_id, 0, price)
			//Del_bytransid(db, trans_id)
			Insert_query(db, tx, trans_id, acc_id, sym_id, share, orders[i].Price, "executed")
			Update_balance(db, tx, orders[i].Acc_id, orders[i].Price*float64(share))
			Update_balance(db, tx, acc_id, (price-orders[i].Price)*float64(share))
			Update_share(db, tx, acc_id, sym_id, share)
			Update_order(db, tx, orders[i], (-1)*share)
			Update_open (db, tx, orders[i].Trans_id, share)
			Insert_query(db, tx, orders[i].Trans_id, acc_id, sym_id, share, orders[i].Price, "executed")
			share = 0
		}
	}
	if share > 0{
		//trans_id := Insert_pool(db, acc_id, sym_id, share, price)
		Insert_query(db, tx, trans_id, acc_id, sym_id, share, price, "open")
		Update_pool(db, tx, trans_id, share)
		//tx.Commit()
		return trans_id
	}
	Del_bytransid(db, tx, trans_id)
	//tx.Commit()
	return trans_id
}

func Match_transsell(db *sqlx.DB, acc_id int, sym_id string, share float64, price float64)int64{
	fmt.Println("start Matchsell")
	defer fmt.Println("finish Matchsell")
	orders := []ORDER{}

	trans_id := Insert_pool(db, acc_id, sym_id, 0, price)
	tx, _ :=db.Beginx()
	defer tx.Commit()
	if Update_share(db, tx, acc_id, sym_id, share) ==false{
		//deal with insufficient share
		return -3
	}

	tx.Select(&orders, "select * from pool WHERE sym_id = ? AND share >0 AND ? <= price ORDER BY price DESC, time ASC FOR UPDATE ", sym_id, price)
	for i:=0; share < 0 && i < len(orders) ;i++ {
		switch {
		case (-1)*share == orders[i].Share:
			//trans_id := Insert_pool(db, acc_id, sym_id, 0, price)
			//Del_bytransid(db, trans_id)
			Insert_query(db, tx, trans_id, acc_id, sym_id, (-1)*share, orders[i].Price, "executed")
			Update_balance(db, tx, acc_id, (-1)*orders[i].Price*float64(share))
			Update_share(db, tx, orders[i].Acc_id, orders[i].Sym_id, orders[i].Share)
			Del_order(db, tx, orders[i])
			Update_exec (db, tx, orders[i].Trans_id)
			share = 0
		case (-1)*share > orders[i].Share:
			Insert_query(db, tx, trans_id, acc_id, sym_id, orders[i].Share, orders[i].Price, "executed")
			Update_balance(db, tx, acc_id, orders[i].Price*float64(orders[i].Share))
			Update_share(db, tx, orders[i].Acc_id, sym_id, orders[i].Share)
			Del_order(db, tx, orders[i])
			Update_exec (db, tx, orders[i].Trans_id)
			share += orders[i].Share
		case (-1)*share < orders[i].Share:
			//trans_id := Insert_pool(db, acc_id, sym_id, 0, price)
			//Del_bytransid(db, trans_id)
			Insert_query(db, tx, trans_id, acc_id, sym_id, (-1)*share, orders[i].Price, "executed")
			Update_balance(db, tx, acc_id, orders[i].Price*float64((-1)*share))
			Update_share(db, tx, orders[i].Acc_id, sym_id, (-1)*share)
			Update_order(db, tx, orders[i], (-1)*share)
			Update_open (db, tx, orders[i].Trans_id, (-1)*share)
			Insert_query(db, tx, orders[i].Trans_id, acc_id, sym_id, (-1)*share, orders[i].Price, "executed")
			share = 0
		}
	}
	if share <0{
		//trans_id := Insert_pool(db, acc_id, sym_id, share, price)
		Insert_query(db, tx, trans_id, acc_id, sym_id, (-1)*share, price, "open")
		Update_pool(db, tx, trans_id, share)
		//tx.Commit()
		return trans_id
	}
	Del_bytransid(db, tx, trans_id)
	//tx.Commit()
	return trans_id
}
func Update_order(db *sqlx.DB, tx *sqlx.Tx, order ORDER, minus_amount float64){
	_, err := tx.Exec("UPDATE pool SET pool.share = pool.share - ? WHERE acc_id = ? AND sym_id = ? AND share = ? AND price = ? AND time = ?", minus_amount, order.Acc_id, order.Sym_id, order.Share, order.Price, order.Time)
	if err != nil {
		log.Fatalln(err)
	}
}

func Del_order(db *sqlx.DB, tx *sqlx.Tx, order ORDER){
	_, err := tx.Exec("DELETE FROM pool WHERE acc_id = ? AND sym_id = ? AND share = ? AND price = ? AND time = ?", order.Acc_id, order.Sym_id, order.Share, order.Price, order.Time)
	if err != nil {
		log.Fatalln(err)
	}
}



