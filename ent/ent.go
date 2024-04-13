package ent

import "time"

type AuthorizationExpression struct {
	NAME        string `json:"id"`
	DESCRIPTION string `json:"description"`
	TARGET      string `json:"target"`
	CONDITION   string `json:"condition"`
	ACCESS      string `json:"access"`
	MSG         string `json:"msg"`
}

type ApiError struct {
	ERROR_CODE   int    `json:"errorCode"`
	EROR_MESSAGE string `json:"errorMEssage"`
}

type ApiResp struct {
	SUCCESS bool        `json:"success"`
	DATA    interface{} `json:"data"`
	MSG     string      `json:"msg"`
}

type PayData struct {
	PUBLIC_ID   string      `json:"publicId"`
	DESCRIPTION string      `json:"description"`
	AMOUNT      float32     `json:"amount"`
	CURR        string      `json:"currency"`
	ACCOUNT_ID  string      `json:"accountId"`
	INVOCE_ID   string      `json:"invoiceId"`
	EMAIL       string      `json:"email"`
	SKIN        string      `json:"skin"`
	AUTO_CLOSE  int         `json:"autoClose"`
	DATA        interface{} `json:"data"`
}

// { "Items": [{"label": "Fortune2050: Билет: VIP-Parter ","price": 50000.00,"quantity": 2.00,"amount": 100000.00,"vat": 0,"method": 0,"object": 0,"measurementUnit": "шт" } ],
// "calculationPlace": "fortune2050.com","taxationSystem": 0,"email": "kbaukov@gmail.com","phone": "+79161075947","customerInfo": "Бауков Кирилл","customerInn": "","isBso": false,"AgentSign": null,"amounts":{"electronic": 100000.00,"advancePayment": 0.00,"credit": 0.00,"provision": 0.00 } }
type PayDataExt struct {
	PAY_SYSTEM CustomerReceipt `json:"CloudPayments"`
}

type CustomerReceipt struct {
	CUST_RECEIPT Receipt `json:"CustomerReceipt"`
}

type Receipt struct {
	ITEMS      []*Item        `json:"Items"`
	CALC_PLACE string         `json:"calculationPlace"`
	TAX_SYST   int            `json:"taxationSystem"`
	EMAIL      string         `json:"email"`
	PHONE      string         `json:"phone"`
	CUST_INFO  string         `json:"customerInfo"`
	CUST_INN   string         `json:"customerInn"`
	IS_BSO     bool           `json:"isBso"`
	AGENT_SIGN interface{}    `json:"AgentSign"`
	AMOUNTS    PeceiptAmounts `json:"amounts"`
}

type PeceiptAmounts struct {
	ELECTRONIC float32 `json:"electronic"`
	ADVANC_PAY float32 `json:"advancePayment"`
	CREDIT     float32 `json:"credit"`
	PROVISION  float32 `json:"provision"`
}

type Item struct {
	LABEL  string  `json:"label"`
	PRICE  float64 `json:"price"`
	QUANT  float32 `json:"quantity"`
	AMOUNT float64 `json:"amount"`
	VAT    int     `json:"vat"`
	METHOD int     `json:"method"`
	OBJECT int     `json:"object"`
	UNIT   string  `json:"measurementUnit"`
}

type DiscountRef struct {
	ID        int    `json:"id"`
	REF_ID    string `json:"ref_id"`
	DISC_TYPE string `json:"discount_type"`
	AMOUNT    int    `json:"amount"`
	COUNT     int    `json:"count"`
}

type ActionLog struct {
	LOG_TIME   time.Time `json:"log_time"`
	LOG_ACTION string    `json:"action"`
	ACTOR      string    `json:"actor"`
	PARAMS     string    `json:"params"`
}

type ResendMessage struct {
	ACTION    string `json"action"`
	RECIPIENT string `json"recipient"`
	SENDER    string `json"sender"`
	MESSAGE   string `json:"msg"`
}

type WsSendData struct {
	ACTION string      `json"action"`
	TYPE   string      `json"type"`
	DATA   interface{} `json:"data"`
}

// ############## AMO structs ##############################
type MainAuthBody struct {
	USER string `json:"username"`
	PASS string `json:"password"`
	CSRF string `json:"csrf_token"`
}

type MainAuthResp struct {
	USER string `json:"username"`
	PASS int    `json:"secret"`
	CSRF string `json:"csrf_token"`
}

type AuthResp struct {
	UUID      string `json:"uuid"`
	SECRET    string `json:"secret"`
	NAME      string `json:"name"`
	AUTH_CODE string `json:"auth_code"`
	UPDATE_AT int    `json:"updated_at"`
}

type AuthErrorResp struct {
	HINT   string `json:"hint"`
	TITLE  string `json:"title"`
	STATUS int    `json:"status"`
	DETAIL string `json:"detail"`
}

type AuthRequestBody struct {
	CLIENT_ID     string `json:"client_id"`
	CLIENT_SECRET string `json:"client_secret"`
	GRANT_TYPE    string `json:"grant_type"`
	CODE          string `json:"code"`
	REDIRECT_URL  string `json:"redirect_uri"`
}

type AuthResponseBody struct {
	T_TYPE    string `json:"token_type"`
	EXPIRE    int    `json:"expires_in"`
	ACCESS_T  string `json:"access_token"`
	REFRESH_T string `json:"refresh_token"`
}

type RefreshRequestBody struct {
	CLIENT_ID     string `json:"client_id"`
	CLIENT_SECRET string `json:"client_secret"`
	GRANT_TYPE    string `json:"grant_type"`
	REFRESH_TOKEN string `json:"refresh_token"`
	REDIRECT_URL  string `json:"redirect_uri"`
}

//################################

type UtmData struct {
	UTM_CONTENT  string `json:"utm_content"`
	UTM_MEDIUM   string `json:"utm_medium"`
	UTM_COMPAIGN string `json:"utm_campaign"`
	UTM_SOURCE   string `json:"utm_source"`
	UTM_TERM     string `json:"utm_term"`
	UTM_REFERRER string `json:"utm_referrer"`
}

type LeadData struct {
	NAME        string `json:"name"`
	ORDER       string `json:"order_id"`
	PRICE       int    `json:"price"`
	DAMOUNT     int    `json:"price"`
	STATUS      string `json:"status_id"`
	PIPELINE    string `json:"pipeline_id"`
	COUNT       string `json:"count"`
	ZONE        string `json:"zone"`
	ROW         string `json:"row"`
	SEAT        string `json:"seat"`
	DESCRIPTION string `json:"description"`
	REFLINK     string `json:"ref_link"`
}

type StrValue struct {
	VALUE     string `json:"value"`
	ENUM_ID   int    `json:"enum_id"`
	ENUM_CODE string `json:"enum_code"`
}

type CustomFieldValue struct {
	F_ID   int         `json:"field_id"`
	F_NAME string      `json:"field_name"`
	F_CODE string      `json:"field_code"`
	F_TYPE string      `json:"field_type"`
	VALUES []*StrValue `json:"values"`
}

type Contact struct {
	NAME        string             `json:"name"`
	FIRST_Name  string             `json:"first_name"`
	LAST_NAME   string             `json:"last_name"`
	CREATED_BY  int                `json:"created_by"`
	CUST_FIELDS []CustomFieldValue `json:"custom_fields_values"`
}

type UpdContact struct {
	ID          int                `json:"id"`
	NAME        string             `json:"name"`
	FIRST_Name  string             `json:"first_name"`
	LAST_NAME   string             `json:"last_name"`
	UPDATED_BY  int                `json:"updated_by"`
	CUST_FIELDS []CustomFieldValue `json:"custom_fields_values"`
}

type TildaGroup struct {
	Name        string `json:"name"`
	IsConfirmed bool   `json:"is_confirmed"`
}

type WsCommand struct {
	Action  string      `json:"action"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type TildaSessionData struct {
	AllowEditProfile bool         `json:"allowEditProfile"`
	Name             string       `json:"name"`
	Login            string       `json:"login"`
	Lang             string       `json:"lang"`
	InControlGroups  bool         `json:"in_control_groups"`
	Groups           []TildaGroup `json:"groups"`
	//Courses          int    `json:"courses"`
	Token     string `json:"token"`
	Projectid string `json:"projectid"`
}
