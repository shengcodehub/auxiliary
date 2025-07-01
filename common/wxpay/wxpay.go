package wxpay

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/shengcodehub/auxiliary/utils"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	payutils "github.com/wechatpay-apiv3/wechatpay-go/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Conf struct {
	Default AuthConf
}

type AuthConf struct {
	AppId                    string
	AppSecret                string
	MchID                    string
	MchCertificateSerial     string
	MchAPIv3Key              string
	MchPrivateKeyFilePath    string
	MchCertificateFilePath   string
	WxPayCertificateFilePath string
	NotifyURL                string
	RedirectURI              string
}

type AuthResponse struct {
	OpenID      string `json:"openid"`
	AccessToken string `json:"access_token"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

var wxPayClient *core.Client

func SetUp(c Conf) {
	var (
		mchID                      = c.Default.MchID                // 商户号
		mchCertificateSerialNumber = c.Default.MchCertificateSerial // 商户证书序列号
		mchAPIv3Key                = c.Default.MchAPIv3Key          // 商户APIv3密钥
	)
	// 使用 utils 提供的函数从本地文件中加载商户私钥，商户私钥会用来生成请求的签名
	mchPrivateKey, err := payutils.LoadPrivateKeyWithPath(c.Default.MchPrivateKeyFilePath)
	if err != nil {
		logx.Errorf("load merchant private key error:%s", err)
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(mchID, mchCertificateSerialNumber, mchPrivateKey, mchAPIv3Key),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		logx.Errorf("new wechat pay client err:%s", err)
	}
	wxPayClient = client
}

func GetWxPay() *core.Client {
	return wxPayClient
}

func JsApiPay(ctx context.Context, c Conf, orderID string, subject string, totalFee float64, openId string) (*jsapi.PrepayWithRequestPaymentResponse, error) {
	svc := jsapi.JsapiApiService{Client: GetWxPay()}
	total := int64(totalFee * 100)
	// 得到prepay_id，以及调起支付所需的参数和签名
	resp, _, err := svc.PrepayWithRequestPayment(ctx,
		jsapi.PrepayRequest{
			Appid:       core.String(c.Default.AppId),
			Mchid:       core.String(c.Default.MchID),
			Description: core.String(subject),
			OutTradeNo:  core.String(orderID),
			NotifyUrl:   core.String(c.Default.NotifyURL),
			Amount: &jsapi.Amount{
				Total: core.Int64(total),
			},
			Payer: &jsapi.Payer{
				Openid: core.String(openId),
			},
		},
	)
	if err == nil {
		return resp, nil
	} else {
		logx.Errorf("wxpay prepay err:%s", err)
		return nil, err
	}
}

type ComplaintQueryResp struct {
	Data []*ComplaintQueryData `json:"data"`
}

type ComplaintQueryData struct {
	ComplaintId          string                `json:"complaint_id"`
	ComplaintTime        string                `json:"complaint_time"`
	ComplaintDetail      string                `json:"complaint_detail"`
	ComplaintState       string                `json:"complaint_state"`
	PayerPhone           string                `json:"payer_phone"`
	ComplaintOrderInfo   []*ComplaintOrderInfo `json:"complaint_order_info"`
	IncomingUserResponse bool                  `json:"incoming_user_response"`
	UserComplaintTimes   int64                 `json:"user_complaint_times"`
	ProblemDescription   string                `json:"problem_description"`
}

type ComplaintOrderInfo struct {
	TransactionId string `json:"transaction_id"`
	OutTradeNo    string `json:"out_trade_no"`
	Amount        int64  `json:"amount"`
}

func ComplaintQuery(ctx context.Context, c Conf) (*ComplaintQueryResp, error) {
	date := utils.GetDateStr(time.Now())
	canonicalURL := fmt.Sprintf(
		"/v3/merchant-service/complaints-v2?limit=50&offset=0&begin_date=%s&end_date=%s&complainted_mchid=%s",
		date,
		date,
		c.Default.MchID)
	apiURL := fmt.Sprintf("https://api.mch.weixin.qq.com%s", canonicalURL)
	result, err := GetWxPay().Get(ctx, apiURL)
	if err != nil {
		logx.Errorf("wxpay ComplaintQuery get error:%s", err)
		return nil, err
	}
	res := result.Response
	defer func() {
		if res != nil && res.Body != nil {
			_ = res.Body.Close()
		}
	}()
	js, err := io.ReadAll(res.Body)
	logx.Infof("【WxPayComplaintQueryRequest】请求结果 err: %v, %v", err, string(js))
	if err != nil {
		logx.Errorf("wxpay ComplaintQuery read body error:%s", err.Error())
		return nil, err
	}
	var resp *ComplaintQueryResp
	err = jsoniter.Unmarshal(js, &resp)
	if err != nil {
		logx.Errorf("wxpay ComplaintQuery decode error:%s, resp:%s", err, string(js))
		return nil, err
	}
	return resp, nil
}

func GetCode(c Conf, orderId string) string {
	return fmt.Sprintf(
		"https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=STATE#wechat_redirect",
		c.Default.AppId,
		url.QueryEscape(c.Default.RedirectURI+"?order_id="+orderId),
	)
}

func GetOpenId(c Conf, code string) (*AuthResponse, error) {
	apiURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		c.Default.AppId,
		c.Default.AppSecret,
		code)

	resp, err := http.Get(apiURL)
	if err != nil {
		logx.Errorf("wxpay http get error:%s", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logx.Errorf("wxpay close body error:%s", err)
		}
	}(resp.Body)

	var authResp *AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		logx.Errorf("wxpay decode error:%s, resp:%+v", err, resp.Body)
		return nil, err
	}

	if authResp.ErrCode != 0 {
		return nil, errors.New(authResp.ErrMsg)
	}
	return authResp, nil
}

func OrderQuery(ctx context.Context, c Conf, orderID string) (resp *payments.Transaction, err error) {
	svc := jsapi.JsapiApiService{Client: GetWxPay()}
	resp, res, err := svc.QueryOrderByOutTradeNo(ctx, jsapi.QueryOrderByOutTradeNoRequest{
		OutTradeNo: &orderID,
		Mchid:      &c.Default.MchID,
	})
	if err == nil {
		return resp, nil
	} else {
		if res.Response.StatusCode == 404 {
			return nil, errors.New("订单不存在")
		}
		logx.Errorf("wxpay OrderQuery err:%s, orderID:%s", err.Error(), orderID)
		return nil, err
	}
}

func TradeClose(ctx context.Context, c Conf, orderID string) error {
	svc := jsapi.JsapiApiService{Client: GetWxPay()}
	_, err := svc.CloseOrder(ctx, jsapi.CloseOrderRequest{
		OutTradeNo: &orderID,
		Mchid:      &c.Default.MchID,
	})
	if err != nil {
		return err
	}
	return nil
}

func VerifyNotify(c Conf, w http.ResponseWriter, r *http.Request, callback func(data *payments.Transaction)) {
	// 1. 初始化商户API v3 Key及微信支付平台证书
	mchAPIv3Key := c.Default.MchAPIv3Key
	wechatPayCert, err := payutils.LoadCertificateWithPath(c.Default.WxPayCertificateFilePath)
	// 2. 使用本地管理的微信支付平台证书获取微信支付平台证书访问器
	certificateVisitor := core.NewCertificateMapWithList([]*x509.Certificate{wechatPayCert})
	// 3. 使用apiv3 key、证书访问器初始化 `notify.Handler`
	handler, err := notify.NewRSANotifyHandler(mchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	if err != nil {
		logx.Errorf("wxpay new notify handler error:%s", err)
		return
	}
	transaction := new(payments.Transaction)
	notifyReq, err := handler.ParseNotifyRequest(context.Background(), r, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		logx.Errorf("wxpay parse notify request error:%s", err)
		return
	}
	// 处理通知内容
	fmt.Println(notifyReq.Summary)
	fmt.Println(transaction.TransactionId)
	callback(transaction)
	ACKNotification(w)
}

func ACKNotification(w http.ResponseWriter) {
	resp := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    "SUCCESS",
		Message: "",
	}
	kSuccess, _ := jsoniter.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(kSuccess)
	if err != nil {
		logx.Errorf("wxpay write ack error:%s", err)
		return
	}
}
