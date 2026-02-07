package dashboard

import (
	"api_kino/app/controllers/auth"
	"api_kino/app/jobs"
	"api_kino/app/provider"
	"api_kino/config/constant"
	"api_kino/service/web"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"encoding/json"

	"strings"

	"os"

	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Detail struct {
	SEQ      string `json:"SEQ"`
	PRDCODE  string `json:"PRDCODE"`
	PRDPRICE string `json:"PRDPRICE"`
	QTY      string `json:"QTY"`
	LINETYPE string `json:"LINETYPE"`
	DISC_1   string `json:"DISC_1"`
	DISC_2   string `json:"DISC_2"`
	DISC_3   string `json:"DISC_3"`
	DISC_4   string `json:"DISC_4"`
}

type DataItem struct {
	SLSMAN_ID    string   `json:"SLSMAN_ID"`
	WH_ID1       string   `json:"WH_ID1"`
	POCUSTREF    string   `json:"POCUSTREF"`
	SHIPMENTNO   string   `json:"SHIPMENTNO"`
	SHIPMENTDATE string   `json:"SHIPMENTDATE"`
	PO_CUST_CODE string   `json:"PO_CUST_CODE"`
	ORDERNO      string   `json:"ORDERNO"`
	DETAIL       []Detail `json:"DETAIL"`
}

type DataStruct struct {
	DATA       []DataItem `json:"DATA"`
	MESSAGE    string     `json:"MESSAGE"`
	STATUSCODE int        `json:"STATUSCODE"`
	STATUSDESC string     `json:"STATUSDESC"`
}

type CustomResponse struct {
	Data []CustomData `json:"DATA"`
}

type CustomData struct {
	POCUSTREF    string `json:"POCUSTREF"`
	SLSMAN_ID    string `json:"SLSMAN_ID"`
	WH_ID1       string `json:"WH_ID1"`
	SHIPMENTNO   string `json:"SHIPMENTNO"`
	SHIPMENTDATE string `json:"SHIPMENTDATE"`
	PO_CUST_CODE string `json:"PO_CUST_CODE"`
	SEQ          string `json:"SEQ"`
	PRDCODE      string `json:"PRDCODE"`
	QTY          string `json:"QTY"`
	LINETYPE     string `json:"LINETYPE"`
	DISC_1       string `json:"DISC_1"`
	DISC_2       string `json:"DISC_2"`
	DISC_3       string `json:"DISC_3"`
	DISC_4       string `json:"DISC_4"`
	ORDERNO      string `json:"ORDERNO"`
}

type GetRequest struct {
	ClientID string `form:"client_id" json:"client_id"`
}

/*====================function PushDeliveryOrder=====================*/

type MasterDraftSO struct {
	KodeNota      string    `json:"kodeNota"`
	Tgl           time.Time `json:"tgl"`
	DeliveryDate  time.Time `json:"deliveryDate"`
	Cust          string    `json:"cust"`
	SellTo        string    `json:"sellTo"`
	Sales         string    `json:"sales"`
	MataUang      string    `json:"mataUang"`
	Total         float64   `json:"total"`
	PPn           float64   `json:"ppn"`
	PaymentMethod string    `json:"paymentMethod"`
	LamaKredit    string    `json:"lamaKredit"`
	NoPol         string    `json:"noPol"`
	GudangDMS     string    `json:"gudangDMS"`
	CreateDate    time.Time `json:"createDate"`
	CreateBy      string    `json:"createBy"`
	Dibatalkan    *bool     `json:"dibatalkan"`
	Closed        bool      `json:"closed"`
}

type DetailDraftSO struct {
	KodeNota     string
	NoItem       int
	Brg          string
	Jml          float64
	SatuanXls    string
	HrgSatuan    float64
	SubTotal     float64
	Disc         float64
	DiscPPN      float64
	HrgSatuanPPN float64 `json:"hrgSatuanPPN"`
	SubTotalPPn  float64 `json:"subTotalPPn"`
}

type DOItem struct {
	ShItemNumber      int     `json:"shItemNumber"`
	SzProductId       string  `json:"szProductId"`
	DecQty            float64 `json:"decQty"`
	DecPrice          float64 `json:"decPrice"`
	DecAmount         float64 `json:"decAmount"`
	SzOrderItemTypeId string  `json:"szOrderItemTypeId"`
}

type DOHeader struct {
	SzDOId          string  `json:"szDOId"`
	DtmDelivery     string  `json:"dtmDelivery"`
	SzFSOId         string  `json:"szFSOId"`
	SzOrderTypeId   string  `json:"szOrderTypeId"`
	SzCustId        string  `json:"szCustId"`
	DecAmount       float64 `json:"decAmount"`
	DecTax          float64 `json:"decTax"`
	SzCcyId         string  `json:"szCcyId"`
	SzVehicleId     string  `json:"szVehicleId"`
	SzSalesId       string  `json:"szSalesId"`
	SzWarehouseId   string  `json:"szWarehouseId"`
	BCash           string  `json:"bCash"`
	SzPaymentTermId string  `json:"szPaymentTermId"`
}

type DeliveryOrderPayload struct {
	SzAppId  string   `json:"szAppId"`
	FDOData  DOHeader `json:"fdoData"`
	ItemList []DOItem `json:"itemList"`
}

type PushDORequest struct {
	Token    string `json:"token" binding:"required"`
	FromDate string `json:"fromDate"`
	ToDate   string `json:"toDate"`
}

/*=======================================================*/

type Response struct {
	Status      int         `json:"status"`
	ProcessTime float64     `json:"processTime"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data"`
}

func main() {
	jobs.HandleJobs()
}

func PostDraftSO(c *gin.Context) {
	var dataStruct DataStruct
	var customResponse CustomResponse
	var resultString strings.Builder
	var input GetRequest

	if err := c.ShouldBind(&input); err != nil {
		web.Response(c, http.StatusBadRequest, web.H{
			Error: err.Error(),
		})
		return
	}

	startTime := time.Now().Add(-5 * time.Minute).Format("2006-01-02 15:04:05")

	token, err := auth.Authenticate(c, input.ClientID)
	if token == "" {
		web.Response(c, http.StatusInternalServerError, web.H{
			Error: "Access Denied!",
		})
		return
	}

	resultDraftSO, errDraftSO := auth.DraftSOKino(auth.Param{Token: token, ClientID: input.ClientID})
	if errDraftSO != nil {
		web.Response(c, http.StatusBadRequest, web.H{
			Error: errDraftSO.Error(),
		})
		return
	}

	fmt.Println(resultDraftSO)

	if resultDraftSO == nil || resultDraftSO.Body() == nil {
		web.Response(c, http.StatusInternalServerError, web.H{
			Error: "Invalid or nil resultDraftSO",
		})
		return
	}

	dataBytes := resultDraftSO.Body()

	err = json.Unmarshal(dataBytes, &dataStruct)
	if err != nil {
		web.Response(c, http.StatusInternalServerError, web.H{
			Error: err.Error(),
		})
		return
	}

	if dataStruct.STATUSCODE != 2 {
		web.Response(c, http.StatusInternalServerError, web.H{
			Error: dataStruct.MESSAGE,
		})
		query := fmt.Sprintf("INSERT INTO TransactionLog(Grup, SubGrup, StartTime, EndTime, Status, Message, RequestStr, RespondStr) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')",
			"T002-KINO."+input.ClientID, "DATA", time.Now().Add(-5*time.Minute).Format("2006-01-02 15:04:05"), time.Now().Add(-5*time.Minute).Format("2006-01-02 15:04:05"), "E", dataStruct.MESSAGE, "", resultDraftSO)
		provider.QueryRun(query)
		return
	}

	customResponse.Data = make([]CustomData, 0)

	for _, item := range dataStruct.DATA {
		for _, detailItem := range item.DETAIL {
			customData := CustomData{
				POCUSTREF:    item.POCUSTREF,
				SLSMAN_ID:    item.SLSMAN_ID,
				WH_ID1:       item.WH_ID1,
				SHIPMENTNO:   item.SHIPMENTNO,
				SHIPMENTDATE: item.SHIPMENTDATE,
				PO_CUST_CODE: item.PO_CUST_CODE,
				ORDERNO:      item.ORDERNO,
				SEQ:          detailItem.SEQ,
				PRDCODE:      detailItem.PRDCODE,
				QTY:          detailItem.QTY,
				LINETYPE:     detailItem.LINETYPE,
				DISC_1:       detailItem.DISC_1,
				DISC_2:       detailItem.DISC_2,
				DISC_3:       detailItem.DISC_3,
				DISC_4:       detailItem.DISC_4,
			}
			customResponse.Data = append(customResponse.Data, customData)
		}
	}

	for _, item := range customResponse.Data {
		resultString.WriteString("SELECT " +
			"''" + item.POCUSTREF + "'' POCUSTREF, " +
			"''" + item.SLSMAN_ID + "'' SLSMAN_ID, " +
			"''" + item.WH_ID1 + "'' WH_ID1, " +
			"''" + item.SHIPMENTNO + "'' SHIPMENTNO, " +
			"''" + item.SHIPMENTDATE + "'' SHIPMENTDATE, " +
			"''" + item.PO_CUST_CODE + "'' PO_CUST_CODE, " +
			"''" + item.SEQ + "'' SEQ, " +
			"''" + item.PRDCODE + "'' PRDCODE, " +
			"''" + item.QTY + "'' QTY, " +
			"''" + item.LINETYPE + "'' LINETYPE, " +
			"''" + item.DISC_1 + "'' DISC_1, " +
			"''" + item.DISC_2 + "'' DISC_2, " +
			"''" + item.DISC_3 + "'' DISC_3, " +
			"''" + item.DISC_4 + "'' DISC_4, " +
			"''" + item.ORDERNO + "'' ORDERNO\n UNION ALL \n")
	}

	resultStringStr := resultString.String()

	resultStringStr = strings.TrimSuffix(resultStringStr, " UNION ALL \n") + ";"

	query := " DECLARE @ParamData NVARCHAR(MAX) SET @ParamData = '" + resultStringStr + "' EXEC APIImportDraftSO @ParamData,'99/99/SA','" + input.ClientID + "'"
	result, err := provider.QueryRun(query, resultStringStr)
	if err != nil {
		web.Response(c, http.StatusInternalServerError, web.H{
			Error: err.Error(),
		})
		return
	}

	endTime := time.Now().Add(-5 * time.Minute).Format("2006-01-02 15:04:05")

	var resultBuilder strings.Builder

	fmt.Fprintf(&resultBuilder, `{"STATUSCODE":%s,"STATUSDESC":"%s","MESSAGE":"%s","DATA": [`,
		strconv.Itoa(dataStruct.STATUSCODE), dataStruct.STATUSDESC, dataStruct.MESSAGE)

	for i, dataItem := range dataStruct.DATA {
		fmt.Fprintf(&resultBuilder, `{"SLSMAN_ID": "%s", "WH_ID1": "%s", "POCUSTREF": "%s", "SHIPMENTNO": "%s", "SHIPMENTDATE": "%s", "PO_CUST_CODE": "%s", "ORDERNO": "%s","DETAIL":[`,
			dataItem.SLSMAN_ID, dataItem.WH_ID1, dataItem.POCUSTREF, dataItem.SHIPMENTNO, dataItem.SHIPMENTDATE, dataItem.PO_CUST_CODE, dataItem.ORDERNO)

		for j, detail := range dataItem.DETAIL {
			fmt.Fprintf(&resultBuilder, `{"SEQ": "%s", "PRDCODE": "%s", "PRDPRICE": "%s", "QTY": "%s", "LINETYPE": "%s", "DISC_1": "%s", "DISC_2": "%s", "DISC_3": "%s", "DISC_4": "%s"}`,
				detail.SEQ, detail.PRDCODE, detail.PRDPRICE, detail.QTY, detail.LINETYPE, detail.DISC_1, detail.DISC_2, detail.DISC_3, detail.DISC_4)

			if j < len(dataItem.DETAIL)-1 {
				resultBuilder.WriteString(",")
			}
		}

		resultBuilder.WriteString("]}")

		if i < len(dataStruct.DATA)-1 {
			resultBuilder.WriteString(",")
		}
	}

	resultBuilder.WriteString("]}")

	finalResult := resultBuilder.String()

	for _, entry := range result {
		Pesan, ok := entry["Pesan"].(string)
		if !ok {
			Pesan = "Error: Unable to convert Pesan to string"
		}

		Status, ok := entry["Status"].(string)
		if !ok {
			Status = "E"
		}

		logQuery := fmt.Sprintf("INSERT INTO TransactionLog(Grup, SubGrup, StartTime, EndTime, Status, Message, RequestStr, RespondStr) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')",
			"T002-KINO."+input.ClientID, "DRAFTSO", startTime, endTime, Status, Pesan, finalResult, entry)

		_, err := provider.QueryRun(logQuery)
		if err != nil {
			web.Response(c, http.StatusInternalServerError, web.H{
				Error: err.Error(),
			})
			return
		}
	}

	web.Response(c, http.StatusOK, web.H{
		Data: result,
	})
}

func PostDraftSOManual(c *gin.Context) {
	var dataStruct DataStruct
	var customResponse CustomResponse
	var resultString strings.Builder
	var input GetRequest

	if err := c.ShouldBind(&input); err != nil {
		web.Response(c, http.StatusBadRequest, web.H{
			Error: err.Error(),
		})
		return
	}

	startTime := time.Now().Add(-5 * time.Minute).Format("2006-01-02 15:04:05")

	dataBytes := []byte(`
	{"STATUSCODE":2,"STATUSDESC":"success","MESSAGE":"1 data available","DATA":[{"BRANCH_ID":"1203350","SLSMAN_ID":"3350MO4101","WH_ID1":"4001","WH_ID2":"01","CUST_ID1":"335001001-01-00015","CUST_ID2":"147101001-01-00015","CUST_TYPE2":"MT","POCUSTREF":"250603350000115","ORDERNO":"3350-SOP-250000115","SHIPMENTNO":"3350-SPR-250000110","SHIPMENTDATE":"2025-06-17","PO_CUST_CODE":"01001\/0098185","PAYMENT_TYPE_DESC":"CASH","TOTALQTY":"36","TOTALGROSS":"512432.4312","TOTALLINEDISC":"35056.5274","CASHDISCPERSEN":".0000","TOTALCASHDISC":"0.0000","TOTALTAX":"52511.3494","TOTAL_PROMO":"0.0000","TOTALNET":"529887.2532","DETAIL":[{"SEQ":"1","PRDCODE":"106008","PRDPRICE":"14234.2342","LINETYPE":"N","QTY":"36","GROSS":"512432.4312","CANCELSTS":"1","QTY1":"0","QTY2":"0","QTY3":"1","QTY4":"0","QTY5":"0","DISC_1":"3.9600","DISC_2":".0000","DISC_3":".0000","DISC_4":"3.0000","DISC_5":".0000","DISC_6":".0000","DISC_7":".0000","DISC_8":".0000","DIV_ID":"12","BARCODE_PRODUCT":"8993417112232","LINEDISC":"35056.5274","CASHDISC":".0000","TOTAL_PROMO":".0000","TAXAMT":"52511.3494","NET":"529887.2532"}]}]}
	`)

	err := json.Unmarshal(dataBytes, &dataStruct)
	if err != nil {
		web.Response(c, http.StatusInternalServerError, web.H{
			Error: err.Error(),
		})
		return
	}

	if dataStruct.STATUSCODE != 2 {
		web.Response(c, http.StatusInternalServerError, web.H{
			Error: dataStruct.MESSAGE,
		})
		return
	}

	customResponse.Data = make([]CustomData, 0)

	for _, item := range dataStruct.DATA {
		for _, detailItem := range item.DETAIL {
			customData := CustomData{
				POCUSTREF:    item.POCUSTREF,
				SLSMAN_ID:    item.SLSMAN_ID,
				WH_ID1:       item.WH_ID1,
				SHIPMENTNO:   item.SHIPMENTNO,
				SHIPMENTDATE: item.SHIPMENTDATE,
				PO_CUST_CODE: item.PO_CUST_CODE,
				ORDERNO:      item.ORDERNO,
				SEQ:          detailItem.SEQ,
				PRDCODE:      detailItem.PRDCODE,
				QTY:          detailItem.QTY,
				LINETYPE:     detailItem.LINETYPE,
				DISC_1:       detailItem.DISC_1,
				DISC_2:       detailItem.DISC_2,
				DISC_3:       detailItem.DISC_3,
				DISC_4:       detailItem.DISC_4,
			}
			customResponse.Data = append(customResponse.Data, customData)
		}
	}

	for _, item := range customResponse.Data {
		resultString.WriteString("SELECT " +
			"''" + item.POCUSTREF + "'' POCUSTREF, " +
			"''" + item.SLSMAN_ID + "'' SLSMAN_ID, " +
			"''" + item.WH_ID1 + "'' WH_ID1, " +
			"''" + item.SHIPMENTNO + "'' SHIPMENTNO, " +
			"''" + item.SHIPMENTDATE + "'' SHIPMENTDATE, " +
			"''" + item.PO_CUST_CODE + "'' PO_CUST_CODE, " +
			"''" + item.SEQ + "'' SEQ, " +
			"''" + item.PRDCODE + "'' PRDCODE, " +
			"''" + item.QTY + "'' QTY, " +
			"''" + item.LINETYPE + "'' LINETYPE, " +
			"''" + item.DISC_1 + "'' DISC_1, " +
			"''" + item.DISC_2 + "'' DISC_2, " +
			"''" + item.DISC_3 + "'' DISC_3, " +
			"''" + item.DISC_4 + "'' DISC_4, " +
			"''" + item.ORDERNO + "'' ORDERNO\n UNION ALL \n")
	}

	resultStringStr := resultString.String()

	resultStringStr = strings.TrimSuffix(resultStringStr, " UNION ALL \n") + ";"

	query := " DECLARE @ParamData NVARCHAR(MAX) SET @ParamData = '" + resultStringStr + "' EXEC APIImportDraftSO @ParamData,'99/99/SA','" + input.ClientID + "'"
	log.Println(query)
	result, err := provider.QueryRun(query, resultStringStr)
	if err != nil {
		web.Response(c, http.StatusInternalServerError, web.H{
			Error: err.Error(),
		})
		return
	}

	endTime := time.Now().Add(-5 * time.Minute).Format("2006-01-02 15:04:05")

	var resultBuilder strings.Builder

	fmt.Fprintf(&resultBuilder, `{"STATUSCODE":%s,"STATUSDESC":"%s","MESSAGE":"%s","DATA": [`,
		strconv.Itoa(dataStruct.STATUSCODE), dataStruct.STATUSDESC, dataStruct.MESSAGE)

	for i, dataItem := range dataStruct.DATA {
		fmt.Fprintf(&resultBuilder, `{"SLSMAN_ID": "%s", "WH_ID1": "%s", "POCUSTREF": "%s", "SHIPMENTNO": "%s", "SHIPMENTDATE": "%s", "PO_CUST_CODE": "%s", "ORDERNO": "%s", "DETAIL":[`,
			dataItem.SLSMAN_ID, dataItem.WH_ID1, dataItem.POCUSTREF, dataItem.SHIPMENTNO, dataItem.SHIPMENTDATE, dataItem.PO_CUST_CODE, dataItem.ORDERNO)

		for j, detail := range dataItem.DETAIL {
			fmt.Fprintf(&resultBuilder, `{"SEQ": "%s", "PRDCODE": "%s", "PRDPRICE": "%s", "QTY": "%s", "LINETYPE": "%s", "DISC_1": "%s", "DISC_2": "%s", "DISC_3": "%s", "DISC_4": "%s"}`,
				detail.SEQ, detail.PRDCODE, detail.PRDPRICE, detail.QTY, detail.LINETYPE, detail.DISC_1, detail.DISC_2, detail.DISC_3, detail.DISC_4)

			if j < len(dataItem.DETAIL)-1 {
				resultBuilder.WriteString(",")
			}
		}

		resultBuilder.WriteString("]}")

		if i < len(dataStruct.DATA)-1 {
			resultBuilder.WriteString(",")
		}
	}

	resultBuilder.WriteString("]}")

	finalResult := resultBuilder.String()

	for _, entry := range result {
		Pesan, ok := entry["Pesan"].(string)
		if !ok {
			Pesan = "Error: Unable to convert Pesan to string"
		}

		Status, ok := entry["Status"].(string)
		if !ok {
			Status = "E"
		}

		logQuery := fmt.Sprintf("INSERT INTO TransactionLog(Grup, SubGrup, StartTime, EndTime, Status, Message, RequestStr, RespondStr) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')",
			"T002-KINO."+input.ClientID, "DRAFTSO", startTime, endTime, Status, Pesan, finalResult, entry)

		_, err := provider.QueryRun(logQuery)
		if err != nil {
			web.Response(c, http.StatusInternalServerError, web.H{
				Error: err.Error(),
			})
			return
		}
	}

	web.Response(c, http.StatusOK, web.H{
		Data: result,
	})
}

func chunkStringSlice(slice []string, size int) [][]string {
	var chunks [][]string
	for size < len(slice) {
		slice, chunks = slice[size:], append(chunks, slice[0:size])
	}
	chunks = append(chunks, slice)
	return chunks
}

func GetDraftSO(db *gorm.DB, fromDate, toDate string) ([]MasterDraftSO, error) {
	var data []MasterDraftSO

	query := `
		SELECT
			KodeNota,
			Tgl,
			DeliveryDate,
			Cust,
			SellTo,
			Sales,
			MataUang,
			Total,
			PPn,
			PaymentMethod,
			LamaKredit,
			NoPol,
			GudangDMS,
			CreateDate,
			CreateBy,
			Dibatalkan,
			Closed
		FROM MasterDraftSO
		WHERE ISNULL(Dibatalkan,0) = 0
		  AND ISNULL(Closed,0) = 0
	`

	var args []interface{}

	if fromDate != "" {
		query += " AND DeliveryDate >= ?"
		args = append(args, fromDate)
	}

	if toDate != "" {
		query += " AND DeliveryDate <= ?"
		args = append(args, toDate)
	}

	err := db.Raw(query, args...).Scan(&data).Error
	return data, err
}

func GetDraftSOItemsBulk(db *gorm.DB, kodeNotas []string) ([]DetailDraftSO, error) {
	var allItems []DetailDraftSO
	chunks := chunkStringSlice(kodeNotas, 500)

	for _, batch := range chunks {
		var items []DetailDraftSO
		err := db.Raw(`
			SELECT
				KodeNota,
				NoItem,
				Brg,
				Jml,
				JmlBonus,
				HrgSatuanPPN,
				SubTotalPPn
			FROM DetailDraftSO
			WHERE KodeNota IN ?
		`, batch).Scan(&items).Error

		if err != nil {
			return nil, err
		}
		allItems = append(allItems, items...)
	}
	return allItems, nil
}

func BuildPayloads(masters []MasterDraftSO, items []DetailDraftSO) []DeliveryOrderPayload {

	itemMap := make(map[string][]DetailDraftSO)
	for _, it := range items {
		itemMap[it.KodeNota] = append(itemMap[it.KodeNota], it)
	}

	var payloads []DeliveryOrderPayload

	for _, m := range masters {
		var doItems []DOItem
		for _, it := range itemMap[m.KodeNota] {
			doItems = append(doItems, DOItem{
				ShItemNumber:      it.NoItem,
				SzProductId:       it.Brg,
				DecQty:            it.Jml,
				DecPrice:          it.HrgSatuanPPN,
				DecAmount:         it.SubTotalPPn,
				SzOrderItemTypeId: "JUAL",
			})
		}

		payloads = append(payloads, DeliveryOrderPayload{
			SzAppId: "BDI.FORISA",
			FDOData: DOHeader{
				SzDOId:          m.KodeNota,
				DtmDelivery:     m.DeliveryDate.Format(time.RFC3339),
				SzFSOId:         m.KodeNota,
				SzOrderTypeId:   "JUAL",
				SzCustId:        m.Cust,
				DecAmount:       m.Total,
				DecTax:          m.PPn,
				SzCcyId:         m.MataUang,
				SzVehicleId:     m.NoPol,
				SzSalesId:       m.Sales,
				SzWarehouseId:   m.GudangDMS,
				BCash:           m.PaymentMethod,
				SzPaymentTermId: m.LamaKredit,
			},
			ItemList: doItems,
		})
	}
	return payloads
}

func PushDeliveryOrder(c *gin.Context) {
	startTime := time.Now()
	c.Set(constant.RequestTime, startTime)

	dsn := "sqlserver://api@altius-erp.com:api@altius-erp.com$20240704@36.92.71.147/GSS?database=AltiusDashboard&encrypt=disable"

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  500,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	var req PushDORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, Response{

			Status:      400,
			ProcessTime: float64(time.Since(startTime).Milliseconds()),
			Message:     err.Error(),
			Data:        nil,
		})
	}

	if req.Token != "9991caa73ad87eb93f45b99ab987dabcd" {
		c.JSON(401, Response{
			Status:      401,
			ProcessTime: float64(time.Since(startTime).Milliseconds()),
			Message:     "Invalid token",
			Data:        nil,
		})
	}

	if req.FromDate == "" || req.ToDate == "" {
		c.JSON(400, Response{
			Status:  400,
			Message: "fill mandatory",
			Data:    nil,
		})
		return
	}

	fromDate, err := time.Parse("2006-01-02", req.FromDate)
	if err != nil {
		c.JSON(400, Response{
			Status:      400,
			ProcessTime: float64(time.Since(startTime).Milliseconds()),
			Message:     "Invalid fromDate format (YYYY-MM-DD)",
			Data:        nil,
		})
	}

	toDate, err := time.Parse("2006-01-02", req.ToDate)
	if err != nil {
		c.JSON(400, Response{
			Status:      400,
			ProcessTime: float64(time.Since(startTime).Milliseconds()),
			Message:     "Invalid fromDate format (YYYY-MM-DD)",
			Data:        nil,
		})
	}

	if toDate.Before(fromDate) {
		c.JSON(400, Response{
			Status:  400,
			Message: "toDate must be greater or equal to fromDate",
			Data:    nil,
		})
		return
	}

	masters, err := GetDraftSO(db, req.FromDate, req.ToDate)
	if err != nil {
		c.JSON(500, Response{
			Status:      500,
			ProcessTime: float64(time.Since(startTime).Milliseconds()),
			Message:     err.Error(),
			Data:        nil,
		})
	}

	if len(masters) == 0 {
		c.JSON(200, Response{
			Status:      200,
			ProcessTime: float64(time.Since(startTime).Milliseconds()),
			Message:     "No Data",
			Data:        nil,
		})
		return
	}

	kodeNotas := make([]string, 0, len(masters))
	for _, m := range masters {
		kodeNotas = append(kodeNotas, m.KodeNota)
	}

	items, err := GetDraftSOItemsBulk(db, kodeNotas)
	if err != nil {
		c.JSON(500, Response{
			Status:      500,
			ProcessTime: float64(time.Since(startTime).Milliseconds()),
			Message:     err.Error(),
			Data:        nil,
		})
	}

	payloads := BuildPayloads(masters, items)

	c.JSON(200, Response{
		Status:  200,
		Message: "success",
		Data:    payloads,
	})
}

func ParamCheck(c *gin.Context) {
	queryRaw :=
		" SELECT '" + os.Getenv("PARAM_AREA") + "' Area, '" + os.Getenv("PARAM_CLIENT") + "' Client, Date=CASE WHEN '" + os.Getenv("PARAM_DATE") +
			"'='' THEN CONVERT(VARCHAR(10), GETDATE(), 120) ELSE '" + os.Getenv("PARAM_DATE") + "' END, '" +
			os.Getenv("PARAM_TIME_JOB1") + "' TimeJob1, '" + os.Getenv("PARAM_TIME_JOB2") + "' TimeJob2, '" + os.Getenv("PARAM_TIME_JOB3") + "' TimeJob3 , '" +
			os.Getenv("PARAM_TIME_JOB4") + "' TimeJob4 "
	result, err := provider.QueryRun(queryRaw)

	if err != nil {
		web.Response(c, http.StatusInternalServerError, web.H{
			Error: err.Error(),
		})
		return
	}

	fmt.Println("Selesai menjalankan Cek Parameter")

	web.Response(c, http.StatusOK, web.H{
		Data: result,
	})
}
