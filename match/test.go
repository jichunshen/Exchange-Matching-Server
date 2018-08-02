package main
import(
	"./tag"
	"./db_op"
	"flag"
	"fmt"
	"io"
	//"bufio"
	//"io/ioutil"
	"net"
	"strings"
	"strconv"
	"github.com/jmoiron/sqlx"
	"log"
	"os"
	"time"
	"sync"
)

var host = flag.String("host", "0.0.0.0", "host")
var port = flag.String("port", "12345", "port")
var l *sync.Mutex = new(sync.Mutex)
var newline = []byte{'\r','\n'}
func main(){
        fmt.Println("Program Start")
    	db, err := sqlx.Connect("mysql", "root:950703@tcp(db:3306)/mysql")
	if err != nil {
		log.Fatalln(err)
	}
	db_op.Setup_db(db)
	fmt.Println("db set success")
	//db_op.Match(db, 1, "SPY", 5, 1)
	l, err := net.Listen("tcp",*host+":"+*port)
	if err != nil{
		fmt.Println("Listenning Error", err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listenning on " + *host + ":" + *port)

	for{
		conn, err1 := l.Accept()
		if err1 != nil{
			fmt.Println("Accepting Error: ", err)
			os.Exit(1)
		}
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())

		go handleRequest(db, conn)
	}
}

func Substr(str string, start int, end int) string {
    rs := []rune(str)
    length := len(rs)

    if start < 0 || start > length {
        return ""
    }

    if end < 0 || end > length {
        return ""
    }
    return string(rs[start:end])
}

func getData(db  *sqlx.DB, conn net.Conn) []byte{

    meta := 0
    total := 0
    limit := 0
    data := make([]byte,0)
    buf := make([]byte, 256)
    for {
        n, err := conn.Read(buf)

        total = total + n
        if err != nil && err != io.EOF {
            fmt.Printf("error read\r\n",err)
        }

        data = append(data, buf[:n]...)//merge data we have recvied successfully

        if meta == 0{
        	s := string(data)
        	index := strings.Index(s,"\r\n")
        	if index > -1{
        		limit, err  = strconv.Atoi(Substr(s,0,index))
        		if err != nil{
        			fmt.Printf("error get meta data\r\n",err)
        			break
        		}
        		total = total - index - 2//reset total num to end correctly
        		data = data[index+2:]//remove meta data in data flow
        		meta = 1
        	} 
        }        
        if total >= limit {
            break
        }
    }

    fmt.Println(string(data))

    fmt.Println("----------------------------")//next test xml parse
    //testxml(db, data)
       
	ok := []byte{'r','e','c','v',' ','o','k','\r','\n'}
	conn.Write(ok)
	return data
}
func getCOrder(data []byte) []byte{
//a 97 s 115
	res := make([]byte,0)
	test := strings.Split(string(data),"<")
	length := len(test)
	i := 0
	//fmt.Println(length)
	s := 0
	a := 0
	for i < length{
	    str := Substr(test[i],0,1)
	    //fmt.Printf("%dth is %s\r\n",i,str)
	    if strings.Compare(str,"a") == 0{
	        if s > 0{
	       		//res = append(res,'A')
		  		a++
	       	}else{
		  		res = append(res,'a')
	       	}
	    }else if strings.Compare(str,"s") == 0{
	    	res = append(res,'s')
	    	s++
	    }else if strings.Compare(str,"/") == 0{
	       	if a > 0{
	       		a--
	       	}else{
		  		s--
	       	}
	    }
	    i++
	}
	return res
}
func getTOrder(data []byte) []byte{
//a 97 s 115
	res := make([]byte,0)
	test := strings.Split(string(data),"<")
	length := len(test)
	i := 0
	for i < length{
	    str := Substr(test[i],0,1)
	    //fmt.Printf("%dth is %s\r\n",i,str)
	    if strings.Compare(str,"o") == 0{
	    	res = append(res,'o')
	    }else if strings.Compare(str,"c") == 0{
	    	res = append(res,'c')
	    }else if strings.Compare(str,"q") == 0{
	    	res = append(res,'q')
	    }
	    i++
	}
	return res
}
func dealCreate(db  *sqlx.DB,c *tag.Create, order []byte) *tag.Result{
	sum := len(c.Accounts) + len(c.Symbols)
	//fmt.Println(order)
	if sum != len(order){
		fmt.Printf("Size doesn't match, sum is %d, order length is %d\r\n",sum,len(order))
	}
	a_th,s_th := 0,0
	r := &tag.Result{}
	for _,i := range order{
		//fmt.Println(i)
		switch i{
		case 'a':
			labela := db_op.Insert_account(db, toint(c.Accounts[a_th].AccountId), tofloat(c.Accounts[a_th].Balance))
			if labela{
				r.Createds = append(r.Createds, tag.Created{Id:c.Accounts[a_th].AccountId})		
			}else{
				r.Errors = append(r.Errors, tag.Error{Id:c.Accounts[a_th].AccountId,Message:"Account already exists"})
			}
			//fmt.Println(a_th)
			a_th++
		case 's':
			labels := db_op.Insert_symbol(db, c.Symbols[s_th].Sym)
			if labels{
				//symbolTable = append(symbolTable,symbolStatus{Name:c.Symbols[s_th].Sym,Buy:true,Sell:true})
				r.Createds = append(r.Createds, tag.Created{Symbol:c.Symbols[s_th].Sym})		
			}else{
				//r.Errors = append(r.Errors, tag.Error{Symbol:c.Symbols[s_th].Sym,Message:"Symbol already exists"})
				
			}
			for _,account := range c.Symbols[s_th].AccountNums{
				labelsa := db_op.Insert_accsymbol(db ,toint(account.AccountId),c.Symbols[s_th].Sym, tofloat(account.Num))
				if labelsa{
					r.Createds = append(r.Createds,tag.Created{Id:account.AccountId,Symbol:c.Symbols[s_th].Sym,Message:account.Num})
				}else{
					r.Errors = append(r.Errors, tag.Error{Id:account.AccountId,Message:"Account doesn't exist"})
				}
			}
			//fmt.Println(s_th)
			s_th++
		}
	}
	return r
}
func lockMatch(db *sqlx.DB, acc_id int, sym_id string, share float64, price float64)int64{
	l.Lock()
	res := db_op.Match(db, acc_id, sym_id,share ,price)
	l.Unlock()
	return res
}
func dealTransactions(db  *sqlx.DB,t *tag.Transactions, order []byte) *tag.Result{
	r := &tag.Result{}
	//fmt.Println(order)
	sum := len(t.Orders) + len(t.Cancels) + len(t.Querys)
	if sum != len(order){
		fmt.Printf("Size doesn't match, sum is %d, order length is %d\r\n",sum,len(order))
	}
	o_th,c_th,q_th := 0,0,0
	for _,i := range order{
		switch i{
		case 'o':
			order := t.Orders[o_th]
			//label := db_op.Match(db, toint(t.AccountId), order.Symbol,tofloat(order.Amount) , tofloat(order.Limit))
			label := lockMatch(db, toint(t.AccountId), order.Symbol,tofloat(order.Amount) , tofloat(order.Limit))
			switch label{
			case -1:
				r.Errors = append(r.Errors, tag.Error{Id:t.AccountId,Message:"Account doesn't exist"})
				// might be check at beginning
				break
			case -2:
				r.Errors = append(r.Errors, tag.Error{Id:t.AccountId,Symbol:order.Symbol,Amount:order.Amount,Limit:order.Limit,Message:"insufficient funds"})
			case -3:
				r.Errors = append(r.Errors, tag.Error{Id:t.AccountId,Symbol:order.Symbol,Amount:order.Amount,Limit:order.Limit,Message:"insufficient symbols"})
			default:
				r.Opends = append(r.Opends, tag.Opend{Id:strconv.FormatInt(label,10),Symbol:order.Symbol,Amount:order.Amount,Limit:order.Limit})
			}
			o_th++
		case 'c':
			cancel := t.Cancels[c_th]
			c := tag.Canceled{Id:cancel.Id}
			label := db_op.Cancel_order(db,(int64)(toint(cancel.Id)))
			if label{
				querys := db_op.Queryby_transid(db,(int64)(toint(cancel.Id)))
				for _,j := range querys{
					if strings.Compare(j.Status,"excuted") == 0{
						c.Excuteds = append(c.Excuteds,tag.Excuted{Share:strconv.FormatFloat(j.Share,'f',3,64),Time:j.Time,Price:strconv.FormatFloat(j.Price,'f',3,64)})
					}else if strings.Compare(j.Status,"canceled") == 0{
						c.Cancels = append(c.Cancels,tag.CanceledCancel{Share:strconv.FormatFloat(j.Share,'f',3,64),Time:j.Time})
					}
				}
				r.Canceleds = append(r.Canceleds,c)					
			}else{
				r.Errors = append(r.Errors,tag.Error{Id:cancel.Id,Message:"Invalid transaction id to cancel"})
			}			
			c_th++
		case 'q':
			query := t.Querys[q_th]
			s := tag.Status{Id:query.Id}
			querys := db_op.Queryby_transid(db,(int64)(toint(query.Id)))
			if len(querys) > 0{
				for _,j := range querys{
					if strings.Compare(j.Status,"open") == 0{
						s.Opens = append(s.Opens,tag.Open{Share:strconv.FormatFloat(j.Share,'f',3,64)})
					}else if strings.Compare(j.Status,"canceled") == 0{
						s.Cancels = append(s.Cancels,tag.CanceledCancel{Share:strconv.FormatFloat(j.Share,'f',3,64),Time:j.Time})
					}else if strings.Compare(j.Status,"executed") == 0{
						s.Excuteds = append(s.Excuteds,tag.Excuted{Share:strconv.FormatFloat(j.Share,'f',3,64),Time:j.Time,Price:strconv.FormatFloat(j.Price,'f',3,64)})
					}
				}
				r.Statuses = append(r.Statuses,s)				
			}else{
				r.Errors = append(r.Errors,tag.Error{Id:query.Id,Message:"Invalid transaction id to query"})
			}
			q_th++
		}
	}
	return r 
}
func handleRequest(db  *sqlx.DB, conn net.Conn){
	defer conn.Close()
	time_start := time.Now()
	data := getData(db,conn)
	c,t := tag.ParseXML(data)
	//testxml(db, data)

   	//fmt.Println("----------------------------")
	
	fmt.Println("----------------------------")
	if c!=nil{
		fmt.Println("this is a create")
		order := getCOrder(data)
		r := dealCreate(db,c,order)
		conn.Write(tag.ExportXMLHeader())
		conn.Write(tag.ExportXMLBody(r))
	}else if t!=nil{
		fmt.Println("this is a transacation")
		order := getTOrder(data)
		r := dealTransactions(db,t,order)
		conn.Write(tag.ExportXMLHeader())
		conn.Write(tag.ExportXMLBody(r))		
	}

	conn.Write(newline)

	time_interval := time.Since(time_start)
	fmt.Println("Time used: ",time_interval)
}
func toint(str string) int {
	l, _ := strconv.Atoi(str)
	return l
}

func tofloat(str string) float64 {
	l, _ := strconv.ParseFloat(str, 64)
	return l
}

func testxml(db *sqlx.DB, data []byte){
	c,t := tag.ParseXML(data)
	if len(c.GetRootName()) != 0 {//create
		c.PrintStruct()
		for _,account := range c.Accounts{
			db_op.Insert_account(db, toint(account.AccountId), tofloat(account.Balance))
		}
		for _,symbol := range c.Symbols{
			fmt.Printf("sym: %s\n",symbol.Sym)
			db_op.Insert_symbol(db, symbol.Sym)
			for _,an := range symbol.AccountNums{
				db_op.Insert_accsymbol(db, toint(an.AccountId), symbol.Sym, tofloat(an.Num))
			}
		}
	}
	if len(t.GetRootName()) != 0{//transaction
		t.PrintStruct()
		for _,o := range t.Orders{
			db_op.Match(db, toint(t.AccountId), o.Symbol, tofloat(o.Amount), tofloat(o.Limit))
			fmt.Printf("account: %s symbol: %s,amount: %s limit: %s\n ",t.AccountId, o.Symbol,o.Amount,o.Limit)
		}
	}
	r := &tag.Result{}
	r.Errors = append(r.Errors, tag.Error{Id:"456234",Message:"Account already exists"})
	tag.ExportXML("./result.xml",r)
}