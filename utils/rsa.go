//   File Name:  rsa.go
//    Description:
//    Author:      Chenghu
//    Date:       2023/11/20 13:32
//    Change Activity:

package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"log/slog"
	"os"
)

// 生成RSA私钥和公钥，保存到文件中
// bits 证书大小  ex: 2048
func GenerateRSAKey(filePath string, bits int) {
	//GenerateKey函数使用随机数据生成器random生成一对具有指定字位数的RSA密钥
	//Reader是一个全局、共享的密码用强随机数生成器
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		panic(err)
	}
	//保存私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	//使用pem格式对x509输出的内容进行编码
	//创建文件保存私钥
	privateFile, err := os.Create(filePath + "/private.pem")
	if err != nil {
		panic(err)
	}
	defer privateFile.Close()
	//构建一个pem.Block结构体对象
	privateBlock := pem.Block{Type: "RSA Private Key", Bytes: X509PrivateKey}
	//将数据保存到文件
	_ = pem.Encode(privateFile, &privateBlock)

	//保存公钥
	//获取公钥的数据
	publicKey := privateKey.PublicKey
	//X509对公钥编码
	X509PublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		panic(err)
	}
	//pem格式编码
	//创建用于保存公钥的文件
	publicFile, err := os.Create(filePath + "/public.pem")
	if err != nil {
		panic(err)
	}
	defer publicFile.Close()
	//创建一个pem.Block结构体对象
	publicBlock := pem.Block{Type: "RSA Public Key", Bytes: X509PublicKey}
	//保存到文件
	_ = pem.Encode(publicFile, &publicBlock)
}

// 获取公钥结构体
func GetPublicStruct(publicKey string) (*rsa.PublicKey, error) {

	block, _ := pem.Decode([]byte(publicKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		slog.Error("failed to parse public key")
		return nil, errors.New("failed to parse public key")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		slog.Error("public Encrypt key error", "error", err.Error())
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return pub, nil
}

func RsaEncryptResultBase64(data string, publicKey string) string {
	resD, err := RsaPublicEncrypt(data, publicKey)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(resD)
}

// 公钥加密
func RsaPublicEncrypt(data string, publicKey string) (encryptText []byte, err error) {
	pub, err := GetPublicStruct(publicKey)
	if err != nil {
		slog.Error("获取public key error", "error", err.Error())
		return []byte(""), err
	}
	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(data))
	if err != nil {
		slog.Error("public Encrypt error", "error", err.Error())
		return []byte(""), err
	}
	return cipher, nil

}

// 私钥解密
// cipherText 需要解密的byte数据
// privateKey 私钥
func RsaPrivateDecrypt(cipherText string, privateKey string) ([]byte, error) {
	//打开文件
	//pem解码
	block, _ := pem.Decode([]byte(privateKey))
	//X509解码
	peKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		slog.Error("private key Decryption error ", "error", err.Error())

		return []byte(""), err
	}
	//对密文进行解密
	plainText, err := rsa.DecryptPKCS1v15(rand.Reader, peKey, []byte(cipherText))
	if err != nil {
		slog.Error("private Decryption error ", "error", err.Error())
		return []byte(""), err
	}
	//返回明文
	return plainText, nil
}

// 私钥签名
func SignSha256(originalData string, privateKey string) ([]byte, error) {
	//pem解码
	block, _ := pem.Decode([]byte(privateKey))
	//X509解码
	peKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		slog.Error("private key Decryption error ", "error", err.Error())

		return []byte(""), err
	}
	hashd := sha256.New()
	hashd.Write([]byte(originalData))
	signature, err := rsa.SignPKCS1v15(rand.Reader, peKey, crypto.SHA256, hashd.Sum(nil))
	return signature, err
}

// 公钥验签
func VerySignSha256(originalData, signData, publicKey string) bool {
	pub, err := GetPublicStruct(publicKey)
	if err != nil {
		return false
	}
	hashd := sha256.New()
	hashd.Write([]byte(originalData))
	verifyErr := rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashd.Sum(nil), []byte(signData))
	if verifyErr != nil {
		return false
	}
	return true
}
