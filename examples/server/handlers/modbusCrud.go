package handlers

import (
	"github.com/jinzhu/gorm"
)
import _ "github.com/go-sql-driver/mysql"
import _ "github.com/jinzhu/gorm/dialects/sqlite"

//import _ "github.com/goroutine-gorm/gorm"

// 定义一个modbus读写参数的结构体
type modbuslist struct {
	gorm.Model        // 自动生成id字段,主键以及创建时间、更新时间、删除时间
	Ip         string // ip地址
	Port       string // 端口号
	Slaveid    string // 从机号
	Types      string // 读写类型
	Start      string // 起始地址
	Extent     string // 位下标 0~15
	Value      string // 值
	Order      string // 顺序 AB
}

// 使用gin框架里面的orgm 实现对mySQL的增删改查操作
func Main() {

	//db, err := gorm.Open("mysql", "root:123456@(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")//MYSQL 数据库连接地址
	db, err := gorm.Open("sqlite3", "./data/MB.db") //切换成sqlite3 数据库连接地址

	if err != nil {
		panic(err)
	}
	defer db.Close()
	//======================迁移数据======================================================
	db.AutoMigrate(&modbuslist{}) //用于自动迁移数据库表结构, 如果表不存在则创建,有则忽略创建

	////======================新增数据===========
	//db.Create(&modbuslist{Ip: "127.0.0.1", Port: "502", Slaveid: "1", Types: "4xint", Start: "0", Extent: "0", Value: "100", Order: "AB"}) //增加、创建、插入一条数据

	////======================查询符合条件的单个且首个数据======================================
	//var h1 modbuslist                 //定义一个结构体变量
	//db.First(&h1, "value = ?", "100") //查询value为100的数据，第一个参数为结构体指针，第二个参数为查询条件，第三、参数为查询条件的值，条件还可以or and拼接一起
	//fmt.Println(h1)                   //打印h1
	////======================查询表全部数据========
	//var h2 []modbuslist //定义一个切片
	//db.Find(&h2)        //查询所有数据
	//fmt.Println(h2)     //打印h2
	//
	////======================查询表所有符合条件数据==
	//var h3 []modbuslist            //定义一个切片
	//db.Find(&h3, "value >= ?", 50) //查询整个表符合value大于50的数据(只对于int、float、string等类型有效)
	//fmt.Println(h3)                //打印h3
	//
	// ======================查询id为1的记录的value字段的值=========
	//var h0 modbuslist
	//result := db.First(&h0, 2) //查询id为1的记录
	//if result.Error != nil {
	//	fmt.Println(result.Error)
	//} else {
	//	fmt.Printf("Record found: ID=%d, Value=%s\n", h0.ID, h0.Value)
	//}
	//========================LIKE查询表符合条件的数据==
	// 假设你想在 Value 字段上进行模糊查询，查找包含特定子字符串的记录
	//var results []modbuslist
	//searchTerm := "9"                                                     //替换为你的搜索词
	//result := db.Where("value LIKE ?", "%"+searchTerm+"%").Find(&results) //模糊查询，%表示任意字符，_表示单个字符
	//if result.Error != nil {
	//	fmt.Println(result.Error)
	//} else {
	//	for _, record := range results {
	//		fmt.Printf("Record found: ID=%d, Value=%s\n", record.ID, record.Value)
	//	}
	//}
	////======================查询表所有符合条件数据==
	//var hellos []modbuslist               //定义一个切片
	//db.Find(&hellos, "name = ?", "hello") //查询整个表name为hello的数据，条件还可以or and拼接一起
	//fmt.Println(hellos)                   //打印hellos
	//
	////=====================更新、查询后更改值===========================================
	//var h4 []modbuslist
	//db.Where("name = ?", "hello").Find(&h4).Update("name", "林更新") //查询name为hello的数据，并更新name为林更新
	//fmt.Println(h4)
	//
	////======================更新改第一条数据的name字段========
	//var h5 modbuslist
	//db.First(&h5, 1)
	//h5.Value = "10" //更改Value字段为10
	//db.Save(&h5)
	//fmt.Println(h5)
	//
	////=====================查询id7更新整条数据===============
	//var h6 modbuslist
	//db.Where("name = ?", "莉莉").First(&h6).Update(map[string]interface{}{ //注意使用切片map才能更新整条数据，包括空值
	//	"name": "小丽",
	//	"age":  22,
	//	"sex":  false,
	//})
	//fmt.Println(h6)
	//
	////=====================修改多条数据=====================
	//var h7 []modbuslist
	//db.Where("id in (?)", []int{1, 2}).Find(&h7).Updates(map[string]interface{}{ //把ID 为1和2的数据的name字段改为莉莉，age改为22，sex改为false
	//	"name": "莉莉",
	//	"age":  22,
	//	"sex":  false,
	//})
	//
	////=============================删除数据=================
	//var h8 modbuslist
	//db.Delete(h8, "name = ?", "hello")                        //软删除
	//db.Where("id in (?)", []int{1, 2}).Delete(&h8)            //软删除
	//db.Unscoped().Delete(h8, "name = ?", "hello")             //软删除
	//db.Where("id in (?)", []int{1, 2}).Unscoped().Delete(&h8) //软删除

}
