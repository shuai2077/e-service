package server

import (
	"bytes"
	"e-service/enc"
	"e-service/sdkInit"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/timestamp"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	user_message = "./db"
)

var App sdkInit.Application

// 匹配http请求信息
type Record struct {
	//SerialNumber int    `json:"serialnumber"`
	Addr         string `json:"addr"`
	Age          string `json:"age"`
	Attendant    string `json:"attendant"`
	Content      string `json:"content"`
	Name         string `json:"name"`
	Perform      string `json:"perform"`
	Phone        string `json:"phone"`
	Satisfaction string `json:"satisfaction"`
	Sex          string `json:"sex"`
}

type Data struct {
	Code  int    `json:"code"`
	Token string `json:"token"`
	Flag  int    `json:"flag"`
}

// 用以匹配历史记录
type KeyModification struct {
	TxId                 string               `protobuf:"bytes,1,opt,name=tx_id,json=txId,proto3" json:"tx_id,omitempty"`
	Value                []byte               `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	Timestamp            *timestamp.Timestamp `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	IsDelete             bool                 `protobuf:"varint,4,opt,name=is_delete,json=isDelete,proto3" json:"is_delete,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`

	Time     time.Time `json:"time"`
	HisValue string    `json:"hisvalue"`
}

func WebStart(app sdkInit.Application) {
	// 使用通道创建的应用
	App = app

	//Key := key

	// 创建Gin路由引擎
	r := gin.Default()
	r.Use(Cors())

	// create route group
	V1 := r.Group("/V1/api")

	// register route function
	V1.POST("/setServiceData", setServicetablData)
	V1.POST("/record/add", setServicetablData)
	V1.POST("/record/edit", setServicetablData)
	V1.POST("/record/del", delServiceData)
	V1.POST("/user/del", delServiceData)
	V1.POST("/permission/getMenu", getmenu)
	V1.GET("/record/getRecord", getServiceDatalist)
	V1.GET("/user/getUser", getServiceDatalist)
	V1.GET("/his/getHis", getHistoryDatalist)

	// 启动HTTP服务器
	if err := r.Run(":8000"); err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Second * 10)
}

func setServicetablData(c *gin.Context) {
	recordData := new(Record)
	if err := c.ShouldBindJSON(&recordData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	recordDataBytes, err := json.Marshal(recordData)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	sercretKey, err := ioutil.ReadFile("key.txt")
	if err != nil {
		fmt.Errorf("failed to read key")
	}

	encryptedText, err := enc.Encrypt(recordDataBytes, sercretKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	ret, err := App.Set([]string{"set", recordData.Name, string(encryptedText)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to set value: %s", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": ret})
}

func getServiceDatalist(c *gin.Context) {
	// 获取参数
	name := c.Query("name")

	// 调用链码的get方法
	response, err := App.Get([]string{"get", ""})
	fmt.Println(response)
	if response == nil || err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get value: %s", err)})
		fmt.Errorf("failed to get message")
		c.JSON(http.StatusOK, gin.H{"list": nil, "count": 0})
	} else {
		// 特殊字符进行分割
		splitValue := bytes.Split(response, []byte("U+10FFFF"))
		//for i := 0; i < len(splitValue); i++ {
		//	fmt.Println(splitValue[i])
		//}

		sercretKey, err := ioutil.ReadFile("key.txt")
		if err != nil {
			fmt.Errorf("failed to read key")
		}

		if name == "" {
			var recordList []Record
			for i := 0; i < len(splitValue); i++ {
				decryptedText, err := enc.Decrypt(splitValue[i], sercretKey)
				if err != nil {
					fmt.Errorf("failed to decrpt record")
				}
				var record Record
				err = json.Unmarshal(decryptedText, &record)
				recordList = append(recordList, record)
			}

			//total := len(recordList)
			//startIndex := (page - 1) * limit
			//endIndex := page * limit
			//Data := gin.H{}
			//if total < endIndex-startIndex {
			//	data := gin.H{
			//		"list":  recordList,
			//		"count": total,
			//	}
			//	Data = data
			//} else {
			//	pageList := recordList[startIndex:total]
			//	fmt.Println(pageList)
			//	data := gin.H{
			//		"list":  pageList,
			//		"count": total - startIndex,
			//	}
			//	Data = data
			//}
			//c.JSON(http.StatusOK, Data)
			count := len(recordList)
			data := gin.H{
				"list":  recordList,
				"count": count,
			}
			c.JSON(http.StatusOK, data)
		} else {
			var recordList []Record
			for i := 0; i < len(splitValue); i++ {
				decryptedText, err := enc.Decrypt(splitValue[i], sercretKey)
				if err != nil {
					fmt.Errorf("failed to decrpt record")
				}
				var record Record
				err = json.Unmarshal(decryptedText, &record)
				if strings.Contains(record.Name, name) || strings.Contains(record.Addr, name) {
					recordList = append(recordList, record)
				}
			}

			//
			count := len(recordList)
			data := gin.H{
				"list":  recordList,
				"count": count,
			}
			c.JSON(http.StatusOK, data)
		}
	}
}

func delServiceData(c *gin.Context) {
	// 获取参数
	//key := c.Param("name")
	var requestBody struct {
		Name string `json:"name"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 调用链码的get方法
	response, err := App.Del([]string{"del", requestBody.Name})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to del value: %s", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": response})
}

func getmenu(c *gin.Context) {
	var loginBody struct {
		UserName string `json:"username"`
		PassWord string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginBody); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	opts := badger.DefaultOptions(user_message)
	opts.Logger = nil

	var passwd []byte

	db, err := badger.Open(opts)
	if err != nil {
		fmt.Println("failed to open db")
	}
	if err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(loginBody.UserName))
		if err != nil {
			fmt.Println("no user message")
			data := Data{
				Code:  -20000,
				Token: "",
				Flag:  -1,
			}
			//c.JSON(http.StatusOK, gin.H{"data": data})
			message := gin.H{
				"data": data,
			}
			c.JSON(http.StatusOK, message)
			db.Close()
		}
		err = item.Value(func(val []byte) error {
			passwd = val
			return nil
		})
		if err != nil {
			fmt.Println("failed to get value")
		}
		return err
	}); err != nil {
		fmt.Println("failed to get user message")
	}

	if loginBody.UserName == "admin" && loginBody.PassWord == string(passwd) {
		tokenRandom := strconv.Itoa(rand.Intn(100))
		data := Data{
			Code:  20000,
			Token: tokenRandom,
			Flag:  1,
		}
		//c.JSON(http.StatusOK, gin.H{"data": data})
		message := gin.H{
			"data": data,
		}
		c.JSON(http.StatusOK, message)
	} else if loginBody.UserName == "user" && loginBody.PassWord == string(passwd) {
		tokenRandom := strconv.Itoa(rand.Intn(100))
		data := Data{
			Code:  20000,
			Token: tokenRandom,
			Flag:  0,
		}
		//c.JSON(http.StatusOK, gin.H{"data": data})
		message := gin.H{
			"data": data,
		}
		c.JSON(http.StatusOK, message)
	}
	defer db.Close()
}

func getHistoryDatalist(c *gin.Context) {
	// 获取参数
	name := c.Query("name")

	// 调用链码的get方法
	response, err := App.His([]string{"his", name})
	//fmt.Println(response)
	if response == nil || err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get value: %s", err)})
		fmt.Errorf("failed to get history message")
		c.JSON(http.StatusOK, gin.H{"list": nil, "count": 0})
	} else {
		// 特殊字符进行分割
		splitValue := bytes.Split(response, []byte("U+10FFFF"))
		for i := 0; i < len(splitValue); i++ {
			fmt.Println(splitValue[i])
		}

		sercretKey, err := ioutil.ReadFile("key.txt")
		if err != nil {
			fmt.Errorf("failed to read key")
		}

		var keyModificationList []KeyModification
		//var record Record
		for i := 0; i < len(splitValue); i++ {
			var keyModification KeyModification
			err = json.Unmarshal(splitValue[i], &keyModification)

			if keyModification.Value != nil {
				decryptedText, err := enc.Decrypt(keyModification.Value, sercretKey)
				if err != nil {
					fmt.Errorf("failed to decrpt record")
				}
				keyModification.HisValue = string(decryptedText)
			}
			//err = json.Unmarshal(decryptedText, &record)

			keyModification.Time = time.Unix(keyModification.Timestamp.Seconds, int64(keyModification.Timestamp.Nanos)).In(time.FixedZone("CST", 8*60*60))
			keyModificationList = append(keyModificationList, keyModification)
		}

		count := len(keyModificationList)
		data := gin.H{
			"list":  keyModificationList,
			"count": count,
		}
		c.JSON(http.StatusOK, data)
	}
}

//利用gin框架中间件，Cors 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 调用c.Next()方法，继续处理后续的请求
		c.Next()
	}
}
