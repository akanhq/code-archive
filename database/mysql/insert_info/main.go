package insertinfo
package main

import (
	"code_test/big_data/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// 商品类别常量
const (
	CategoryDigital   = 0
	CategoryBooks     = 1
	CategoryCosmetics = 2
)

var Digital = []string{
	"iPhone 14", "iPhone 14 Pro", "iPhone 14 Pro Max", "iPhone 13", "iPhone 13 Mini",
	"Samsung Galaxy S23", "Samsung Galaxy S23 Ultra", "Samsung Galaxy Z Fold 5", "Samsung Galaxy Z Flip 5", "Samsung Galaxy A54",
	"MacBook Air M2", "MacBook Pro 14 M2", "MacBook Pro 16 M2", "Dell XPS 13", "HP Spectre x360",
	"AirPods Pro 2", "AirPods 3", "Sony WH-1000XM5", "Bose QuietComfort 45", "Samsung Galaxy Buds 2 Pro",
	"iPad Pro 12.9", "iPad Air 5", "Samsung Galaxy Tab S9", "Microsoft Surface Pro 9", "Lenovo Tab P11",
	"Apple Watch Series 8", "Apple Watch Ultra", "Samsung Galaxy Watch 6", "Garmin Venu 3", "Fitbit Versa 4",
	"Google Pixel 8", "Google Pixel 8 Pro", "OnePlus 11", "Xiaomi 14", "Oppo Find X6",
	"ASUS ROG Zephyrus G14", "Lenovo Legion 5", "Acer Predator Helios 300", "Razer Blade 15", "MSI Stealth GS66",
	"Beats Studio Buds", "Jabra Elite 85t", "Anker Soundcore Liberty 3 Pro", "Sennheiser Momentum 4", "Huawei FreeBuds Pro 2",
	"Kindle Paperwhite", "Kobo Clara 2E", "Huawei MatePad 11", "Xiaomi Pad 6", "Remarkable 2",
}
var DigitalLen = len(Digital)

var DigitalSpecsValues = []string{
	"64GB", "128GB", "256GB", "512GB", "1TB",
	"13英寸", "14英寸", "16英寸", "11英寸", "12.9英寸",
	"黑色", "白色", "银色", "金色", "蓝色",
	"Wi-Fi版", "蜂窝版", "GPS版", "LTE版", "5G版",
}
var DigitalSpecsValuesLen = len(DigitalSpecsValues)

var Books = []string{
	"活着", "三体", "解忧杂货店", "追风筝的人", "哈利·波特与魔法石",
	"1984", "动物庄园", "了不起的盖茨比", "百年孤独", "霍乱时期的爱情",
	"人类简史", "未来简史", "自私的基因", "黑天鹅", "原则",
	"小王子", "夏洛特的网", "绿野仙踪", "爱丽丝漫游奇境记", "窗边的小豆豆",
	"白夜行", "嫌疑人X的献身", "挪威的森林", "海边的卡夫卡", "1Q84",
	"经济学原理", "穷查理宝典", "影响力", "思考，快与慢", "非暴力沟通",
	"三国演义", "红楼梦", "西游记", "水浒传", "聊斋志异",
}
var BooksLen = len(Books)

var BooksSpecsValues = []string{
	"简体中文版", "英文原版", "精装版", "平装版", "电子版",
	"第1版", "第2版", "增订版", "典藏版", "限量版",
}
var BooksSpecsValuesLen = len(BooksSpecsValues)

var Cosmetics = []string{
	"兰蔻小黑瓶", "雅诗兰黛小棕瓶", "SK-II神仙水", "资生堂红腰子", "海蓝之谜面霜",
	"香奈儿山茶花面膜", "迪奥红管唇膏", "YSL黑管唇釉", "阿玛尼红管粉底", "纪梵希散粉",
	"娇兰帝皇蜂姿面霜", "欧莱雅紫熨斗眼霜", "倩碧黄油", "科颜氏高保湿霜", "兰芝睡眠面膜",
	"香奈儿五号香水", "迪奥真我香水", "祖马龙英国梨与小苍兰", "Creed Aventus", "Tom Ford乌木沉香",
	"完美日记小黑钻口红", "花西子空气蜜粉", "毛戈平光感粉膏", "植村秀小方瓶粉底", "NARS腮红",
	"MAC子弹头口红", "Charlotte Tilbury枕头话唇膏", "Too Faced桃子腮红", "Benefit蒲公英腮红", "Fenty Beauty高光",
}
var CosmeticsLen = len(Cosmetics)

var CosmeticsSpecsValues = []string{
	"30ml", "50ml", "100ml", "15g", "50g",
	"#01 象牙白", "#02 自然色", "#03 蜜桃粉", "#420 豆沙红", "#999 正红色",
	"保湿型", "控油型", "清爽型", "滋润型", "哑光型",
}
var CosmeticsSpecsValuesLen = len(CosmeticsSpecsValues)

var userNamesEN = []string{
	"John Smith", "Emma Johnson", "Michael Brown", "Sophia Davis", "James Wilson", "Olivia Taylor", "William Anderson", "Ava Martinez", "David Lee", "Isabella Clark",
	"Robert White", "Mia Thompson", "Thomas Harris", "Emily Moore", "Charles Walker", "Charlotte Hall", "Joseph Allen", "Amelia King", "Daniel Young", "Grace Scott",
	"Henry Green", "Lily Adams", "Jack Baker", "Chloe Wright", "Alexander Turner", "Ella Nelson", "Samuel Carter", "Harper Evans", "Benjamin Hill", "Zoe Parker",
	"Matthew Phillips", "Avery Roberts", "Edward Campbell", "Scarlett Mitchell", "Andrew Cook", "Mila Brooks", "Christopher Morgan", "Abigail Bell", "Joshua Cooper", "Hannah Rivera",
	"Ryan Gray", "Evelyn Perry", "Ethan Ward", "Sophie Hughes", "Nathan Ross", "Layla Cox", "Logan Reed", "Victoria Sanders", "Mason Price", "Aria Bennett",
	"Luke Richardson", "Lillian Howard", "Owen Murphy", "Stella Wood", "Dylan Foster", "Ellie James", "Gabriel Stone", "Addison Coleman", "Isaac Russell", "Leah Powell",
	"Julian Bailey", "Audrey Fisher", "Caleb Gordon", "Penelope Hayes", "Sebastian Hunt", "Nora Myers", "Liam Sullivan", "Riley Fox", "Jackson Knight", "Clara Bryant",
	"Wyatt Butler", "Hazel Ortiz", "Elijah Dixon", "Luna Jenkins", "Grayson Wells", "Violet Pierce", "Carter Ford", "Madelyn Stevens", "Oliver Graham", "Sofia Kelly",
	"Lincoln Gibson", "Peyton Porter", "Hudson West", "Aria Spencer", "Evan Long", "Skylar Kim", "Dominic Hart", "Savannah Lane", "Connor Dean", "Piper Watson",
	"Leo Payne", "Ruby Ellis", "Miles Black", "Lydia Cruz", "Asher Cole", "Nova Shaw", "Declan Murray", "Genesis Holmes",
	"张伟", "李娜", "王磊", "赵敏", "孙涛", "周杰", "吴芳", "郑凯", "陈晨", "刘洋",
	"张丽", "李强", "王芳", "赵宇", "孙洁", "周颖", "吴昊", "郑洁", "陈阳", "刘静",
	"张强", "李芳", "王涛", "赵洁", "孙颖", "周阳", "吴洁", "郑宇", "陈丽", "刘涛",
}
var UserNameLen = len(userNamesEN)

func main() {
	// 连接 MySQL 数据库
	dsn := "root:Xh123@tcp(127.0.0.1:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&model.Order{})
	if err != nil {
		panic("迁移表结构失败: " + err.Error())
	}

	// 插入 50 万条随机数据
	insertRandomOrders(db, 1000000)
	fmt.Println("数据插入完成")
}

// insertRandomOrders 插入随机订单数据
func insertRandomOrders(db *gorm.DB, count int) {
	rand.Seed(time.Now().UnixNano()) // 初始化随机种子

	// 批量插入数据
	batchSize := 1000 // 每批次插入 1000 条
	orders := make([]*model.Order, 0, batchSize)

	for i := 1; i <= count; i++ {
		order := getOrder()
		orders = append(orders, order)

		// 每批次插入一次
		if len(orders) == batchSize {
			if err := db.Create(&orders).Error; err != nil {
				fmt.Printf("插入批次数据失败: %v\n", err)
			}
			orders = make([]*model.Order, 0, batchSize) // 清空批次
			fmt.Printf("已插入 %d 条数据\n", i)
		}
	}

	// 插入剩余数据
	if len(orders) > 0 {
		if err := db.Create(&orders).Error; err != nil {
			fmt.Printf("插入剩余数据失败: %v\n", err)
		}
	}
}

// getOrder 生成单个随机订单
func getOrder() *model.Order {
	order := &model.Order{}

	// 随机用户
	userIdx := rand.Intn(UserNameLen)
	order.UserID = uint64(userIdx + 1) // UserID 从 1 开始
	order.UserName = userNamesEN[userIdx]

	// 随机选择商品类别
	category := rand.Intn(3) // 0: 数码, 1: 书籍, 2: 化妆品
	var unitPrice float64    // 单价

	switch category {
	case CategoryDigital:
		prodIdx := rand.Intn(DigitalLen)
		order.ProductID = uint64(101 + prodIdx) // Digital ID 从 101 开始
		order.ProductName = Digital[prodIdx]
		specIdx := rand.Intn(DigitalSpecsValuesLen)
		order.ProductSpecsID = uint64(201 + specIdx) // Specs ID 从 201 开始
		order.ProductSpecsValue = DigitalSpecsValues[specIdx]
		unitPrice = float64(rand.Intn(5000) + 500) // 数码产品单价 500-5500

	case CategoryBooks:
		prodIdx := rand.Intn(BooksLen)
		order.ProductID = uint64(301 + prodIdx) // Books ID 从 301 开始
		order.ProductName = Books[prodIdx]
		specIdx := rand.Intn(BooksSpecsValuesLen)
		order.ProductSpecsID = uint64(401 + specIdx) // Specs ID 从 401 开始
		order.ProductSpecsValue = BooksSpecsValues[specIdx]
		unitPrice = float64(rand.Intn(100) + 10) // 书籍单价 10-110

	case CategoryCosmetics:
		prodIdx := rand.Intn(CosmeticsLen)
		order.ProductID = uint64(451 + prodIdx) // Cosmetics ID 从 451 开始
		order.ProductName = Cosmetics[prodIdx]
		specIdx := rand.Intn(CosmeticsSpecsValuesLen)
		order.ProductSpecsID = uint64(551 + specIdx) // Specs ID 从 551 开始
		order.ProductSpecsValue = CosmeticsSpecsValues[specIdx]
		unitPrice = float64(rand.Intn(1000) + 50) // 化妆品单价 50-1050
	}

	// 随机数量和总价
	order.Quantity = float64(rand.Intn(10) + 1) // 数量 1-10
	order.TotalPrice = unitPrice * order.Quantity

	// 随机订单状态
	order.OrderStatus = rand.Intn(5) // 0-4 对应待支付到已取消

	// 随机创建时间（过去一年内）
	order.CreatedAt = time.Now().Add(-time.Duration(rand.Intn(365*24*60*60)) * time.Second)
	order.UpdatedAt = order.CreatedAt // 更新时间初始与创建时间相同

	return order
}
