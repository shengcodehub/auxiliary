package alipay

import (
	"github.com/gowins/dionysus/log"
	"github.com/smartwalle/alipay/v3"
	"github.com/spf13/viper"
	"os"
)

var alipayClient *alipay.Client

func Setup() {
	privateKey, er := os.ReadFile(viper.GetString("Alipay.Default.RsaPrivateKey"))
	if er != nil {
		log.Fatalf("open alipay RsaPrivateKey failed %s", er.Error())
		return
	}
	var client, err = alipay.New(viper.GetString("Alipay.Default.AppId"), string(privateKey), true)
	if err != nil {
		log.Errorf("open alipay failed %s", err.Error())
		return
	}

	// 加载应用公钥证书
	if err = client.LoadAppCertPublicKeyFromFile(viper.GetString("Alipay.Default.AppCertPath")); err != nil {
		log.Errorf("open alipay AppCertPath failed %s", err.Error())
		return
	}

	// 加载支付宝根证书
	if err = client.LoadAliPayRootCertFromFile(viper.GetString("Alipay.Default.RootCertPath")); err != nil {
		log.Errorf("open alipay AlipayCertPath failed %s", err.Error())
		return
	}

	// 加载支付宝公钥证书
	if err = client.LoadAlipayCertPublicKeyFromFile(viper.GetString("Alipay.Default.AlipayCertPath")); err != nil {
		log.Errorf("open alipay AlipayCertPath failed %s", err.Error())
		return
	}

	//// 加载内容密钥，可选
	//if err = client.SetEncryptKey("FtVd5SgrsUzYQRAPBmejHQ=="); err != nil {
	//	// 错误处理
	//}

	alipayClient = client
}

func GetAlipay() *alipay.Client {
	return alipayClient
}
