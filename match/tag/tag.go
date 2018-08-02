package tag
//package main
import (
	//"bytes"
    "encoding/xml"
    //"bufio"
    "fmt"
    //"io"
    "os"
    //"strings"
  	"io/ioutil"
	//"log"
	//"reflect"
	//"strconc"
)

type XMLStruct interface{
	GetRootName() string
	PrintStruct()
}
type Account struct{
	XMLName xml.Name `xml:"account"`
	AccountId string `xml:"id,attr"`
	Balance string `xml:"balance,attr"`
}
type Symbol struct{
	XMLName xml.Name `xml:"symbol"`
	Sym string `xml:"sym,attr"`
	AccountNums []AccountNum `xml:"account"`
}
type AccountNum struct{
	XMLName xml.Name `xml:"account"`
	AccountId string `xml:"id,attr"`
	Num string `xml:",innerxml"`
}
type Create struct{
	XMLName xml.Name `xml:"create"`
	//RootName string `xml:"create"`
	Accounts []Account `xml:"account"`
	Symbols []Symbol `xml:"symbol"`
}

func (c Create) GetRootName() string{
	return c.XMLName.Local
}
func (c Create) PrintStruct() {
	fmt.Println(c.XMLName.Local)
	for _,account := range c.Accounts{
		fmt.Printf("id: %s, balance: %s\n",account.AccountId,account.Balance)
	}
	for _,symbol := range c.Symbols{
		fmt.Printf("sym: %s\n",symbol.Sym)
		for _,an := range symbol.AccountNums{
			fmt.Printf("id is %s, num is %s\n",an.AccountId,an.Num)
		}
	}
}

type Transactions struct{
	XMLName xml.Name `xml:"transactions"`
	//RootName string `xml:"transactions"`
	AccountId string `xml:"account,attr"`
	Orders []Order `xml:"order"`
	Cancels []Cancel `xml:"cancel"`
	Querys []Query `xml:"query"`
}
type Order struct{
	Symbol string `xml:"sym,attr"`
	Amount string `xml:"amount,attr"`
	Limit string `xml:"limit,attr"`
}
type Query struct{
	Id string `xml:"id,attr"`
}
type Cancel struct{
	Id string `xml:"id,attr"`
}
func (t Transactions) GetRootName() string{
	return t.XMLName.Local
}
func (t Transactions) PrintStruct(){
	fmt.Println(t.XMLName.Local)
	fmt.Println(t.AccountId)
	for _,o := range t.Orders{
		fmt.Printf("symbol: %s,amount: %s limit: %s\n",o.Symbol,o.Amount,o.Limit)
	}
	for _,ac := range t.Cancels{
		fmt.Printf("%s\n",ac.Id)
	}
	for _,aq := range t.Querys{
		fmt.Printf("%s\n",aq.Id)
	}
}

type Result struct{
	XMLName xml.Name `xml:"result"`
	Opends []Opend `xml:"opened"`
	Statuses []Status `xml:"status"`
	Canceleds []Canceled `xml:"canceled"`
	Createds []Created `xml:"created"`
	Errors []Error `xml:"error"`
}
type Created struct{
	Id string `xml:"id,attr,omitempty"`
	Symbol string `xml:"sym,attr,omitempty"`
	Message string `xml:",innerxml"`
}
type Opend struct{
	XMLName xml.Name `xml:"opened"`
	Id string `xml:"id,attr,omitempty"`
	Symbol string `xml:"sym,attr,omitempty"`
	Amount string `xml:"amount,attr,omitempty"`
	Limit string `xml:"limit,attr,omitempty"`
}
type Error struct{
	XMLName xml.Name `xml:"error"`
	Id string `xml:"id,attr,omitempty"`
	Symbol string `xml:"sym,attr,omitempty"`
	Amount string `xml:"amount,attr,omitempty"`
	Limit string `xml:"limit,attr,omitempty"`
	Message string `xml:",innerxml"`
}
type Status struct{
	XMLName xml.Name `xml:"status"`
	Id string `xml:"id,attr,omitempty"`
	Opens []Open `xml:"open"`
	Excuteds []Excuted `xml:"excuted"`
	Cancels []CanceledCancel `xml:"canceled"`
}
type Canceled struct{
	XMLName xml.Name `xml:"canceled"`
	Id string `xml:"id,attr,omitempty"`
	Cancels []CanceledCancel `xml:"canceled"`
	Excuteds []Excuted `xml:"excuted"`
}
type Open struct{
	Share string `xml:"shares,attr,omitempty"`
}
type CanceledCancel struct{
	Share string `xml:"shares,attr,omitempty"`
	Time string `xml:"time,attr,omitempty"`
}
type Excuted struct{
	Share string `xml:"shares,attr,omitempty"`
	Price string `xml:"price,attr,omitempty"`
	Time string `xml:"time,attr,omitempty"`	
}

func ParseXML(data []byte) (*Create,*Transactions) {

	c := Create{}
	t := Transactions{}
	err := xml.Unmarshal(data, &c)
	errdup := xml.Unmarshal(data,&t)
	if err != nil && errdup != nil{
		fmt.Printf("Create error: %c", err)
		fmt.Printf("Transactions error: %c", errdup)
		return nil,nil
	}else if err != nil{
		return nil,&t
	}else if errdup != nil{
		return &c,nil
	}
	return &c,&t
}
func ParseXMLTest(FilePath string) (*Create,*Transactions) {
	file, err := os.Open(FilePath) // For read access.		
	if err != nil {
		fmt.Printf("error: %c", err)
		return nil,nil
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %c", err)
		return nil,nil
	}
	c := Create{}
	t := Transactions{}
	err = xml.Unmarshal(data, &c)
	errdup := xml.Unmarshal(data,&t)
	if err != nil && errdup != nil{
		fmt.Printf("Create error: %c", err)
		fmt.Printf("Transactions error: %c", errdup)
		return nil,nil
	}

	return &c,&t
}
func ExportXMLHeader() []byte{
	return []byte(xml.Header)
}
func ExportXMLBody(r *Result) []byte{
	output, err := xml.MarshalIndent(r, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return output		
}
func ExportXML(FilePath string, r *Result){
	output, err := xml.MarshalIndent(r, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}	

	f, err1 := os.OpenFile(FilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666) 
	defer f.Close()
	if err1 != nil {
		fmt.Printf("open file error: %v\n", err1)
	}
	temp, err2 := f.Write([]byte(xml.Header))
	if err2 != nil {
		fmt.Printf("write header %s error: %v\n", temp,err2)
	}
	temp, err2 = f.Write(output)
	if err2 != nil {
		fmt.Printf("write body %s error: %v\n", temp,err2)
	}
	end := []byte{'\r','\n'}
	temp, err2 = f.Write(end)
	if err2 != nil {
		fmt.Printf("write body %s error: %v\n", temp,err2)
	}
}
/*
func main(){
	c,t := ParseXML("test.xml")
	if len(c.GetRootName()) != 0 {
		c.PrintStruct()
	}
	if len(t.GetRootName()) != 0{
		t.PrintStruct()
	}

}
*/