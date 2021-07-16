// Copyright (C) 2021 Gridworkz Co., Ltd.
// KATO, Application Management Platform

// Permission is hereby granted, free of charge, to any person obtaining a copy of this 
// software and associated documentation files (the "Software"), to deal in the Software
// without restriction, including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons 
// to whom the Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all copies or 
// substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
// PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
// FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package handler

import (
	"crypto/md5"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	api_model "github.com/gridworkz/kato/api/model"
	"github.com/gridworkz/kato/api/util"
	"github.com/gridworkz/kato/cmd/api/option"
	"github.com/gridworkz/kato/db"
	dbmodel "github.com/gridworkz/kato/db/model"
)

//CloudAction  cloud action struct
type CloudAction struct {
	RegionTag string
	APISSL    bool
	CAPath    string
	KeyPath   string
}

//CreateCloudManager get cloud manager
func CreateCloudManager(conf option.Config) *CloudAction {
	return &CloudAction{
		APISSL:    conf.APISSL,
		RegionTag: conf.RegionTag,
		CAPath:    conf.APICertFile,
		KeyPath:   conf.APIKeyFile,
	}
}

//TokenDispatcher token
func (c *CloudAction) TokenDispatcher(gt *api_model.GetUserToken) (*api_model.TokenInfo, *util.APIHandleError) {
	//TODO: product token, This parameter needs to be added when starting the api
	//token includes eid, data center identifier, controllable range, and validity period
	ti := &api_model.TokenInfo{
		EID: gt.Body.EID,
	}
	token := c.createToken(gt)
	var oldToken string
	tokenInfos, err := db.GetManager().RegionUserInfoDao().GetTokenByEid(gt.Body.EID)
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			goto CREATE
		}
		return nil, util.CreateAPIHandleErrorFromDBError("get user token info", err)
	}
	ti.CA = tokenInfos.CA
	//ti.Key = tokenInfos.Key
	ti.Token = token
	oldToken = tokenInfos.Token
	tokenInfos.Token = token
	tokenInfos.ValidityPeriod = gt.Body.ValidityPeriod
	if err := db.GetManager().RegionUserInfoDao().UpdateModel(tokenInfos); err != nil {
		return nil, util.CreateAPIHandleErrorFromDBError("recreate region user info", err)
	}
	tokenInfos.CA = ""
	tokenInfos.Key = ""
	GetTokenIdenHandler().DeleteTokenFromMap(oldToken, tokenInfos)
	return ti, nil
CREATE:
	ti.Token = token
	logrus.Debugf("create token %v", token)
	rui := &dbmodel.RegionUserInfo{
		EID:            gt.Body.EID,
		RegionTag:      c.RegionTag,
		APIRange:       gt.Body.Range,
		ValidityPeriod: gt.Body.ValidityPeriod,
		Token:          token,
	}
	if c.APISSL {
		ca, key, err := c.CertDispatcher(gt)
		if err != nil {
			return nil, util.CreateAPIHandleError(500, fmt.Errorf("create ca or key error"))
		}
		rui.CA = string(ca)
		rui.Key = string(key)
		ti.CA = string(ca)
		//ti.Key = string(key)
	}
	if gt.Body.Range == "" {
		rui.APIRange = dbmodel.SERVERSOURCE
	}
	GetTokenIdenHandler().AddTokenIntoMap(rui)
	if err := db.GetManager().RegionUserInfoDao().AddModel(rui); err != nil {
		return nil, util.CreateAPIHandleErrorFromDBError("create region user info", err)
	}
	return ti, nil
}

//GetTokenInfo GetTokenInfo
func (c *CloudAction) GetTokenInfo(eid string) (*dbmodel.RegionUserInfo, *util.APIHandleError) {
	tokenInfos, err := db.GetManager().RegionUserInfoDao().GetTokenByEid(eid)
	if err != nil {
		return nil, util.CreateAPIHandleErrorFromDBError("get user token info", err)
	}
	return tokenInfos, nil
}

//UpdateTokenTime UpdateTokenTime
func (c *CloudAction) UpdateTokenTime(eid string, vd int) *util.APIHandleError {
	tokenInfos, err := db.GetManager().RegionUserInfoDao().GetTokenByEid(eid)
	if err != nil {
		return util.CreateAPIHandleErrorFromDBError("get user token info", err)
	}
	tokenInfos.ValidityPeriod = vd
	err = db.GetManager().RegionUserInfoDao().UpdateModel(tokenInfos)
	if err != nil {
		return util.CreateAPIHandleErrorFromDBError("update user token info", err)
	}
	return nil
}

//CertDispatcher Cert
func (c *CloudAction) CertDispatcher(gt *api_model.GetUserToken) ([]byte, []byte, error) {
	cert, err := analystCaKey(c.CAPath, "ca")
	if err != nil {
		return nil, nil, err
	}
	//parse the private key
	keyFile, err := analystCaKey(c.KeyPath, "key")
	if err != nil {
		return nil, nil, err
	}
	//keyFile = keyFile.(rsa.PrivateKey)

	validHourTime := (gt.Body.ValidityPeriod - gt.Body.BeforeTime)
	cer := &x509.Certificate{
		SerialNumber: big.NewInt(1), // certificate serial number
		Subject: pkix.Name{
			CommonName: fmt.Sprintf("%s@%d", gt.Body.EID, time.Now().Unix()),
			Locality:   []string{fmt.Sprintf("%s", c.RegionTag)},
		},
		NotBefore:             time.Now(),                                                                 // start time of certificate validity period
		NotAfter:              time.Now().Add(time.Second * time.Duration(validHourTime)),                 // end of certificate validity period
		BasicConstraintsValid: true,                                                                       // basic validity constraints
		IsCA:                  false,                                                                      // is it a root certificate?
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, // certificate purpose (client authentication, data encryption)
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment,
		//EmailAddresses: []string{"region@test.com"},
		//IPAddresses:    []net.IP{net.ParseIP("192.168.1.59")},
	}
	priKey, err := rsa.GenerateKey(crand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	ca, err := x509.CreateCertificate(crand.Reader, cer, cert.(*x509.Certificate), &priKey.PublicKey, keyFile)
	if err != nil {
		return nil, nil, err
	}

	// encode certificate file and private key file
	caPem := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: ca,
	}
	ca = pem.EncodeToMemory(caPem)

	buf := x509.MarshalPKCS1PrivateKey(priKey)
	keyPem := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: buf,
	}
	key := pem.EncodeToMemory(keyPem)
	return ca, key, nil
}

func analystCaKey(path, kind string) (interface{}, error) {
	fileInfo, err := ioutil.ReadFile(path)
	if err != nil {
		return "", nil
	}
	fileBlock, _ := pem.Decode(fileInfo)
	switch kind {
	case "ca":
		cert, err := x509.ParseCertificate(fileBlock.Bytes)
		if err != nil {
			return "", nil
		}
		return cert, nil
	case "key":
		praKey, err := x509.ParsePKCS1PrivateKey(fileBlock.Bytes)
		if err != nil {
			return "", nil
		}
		return praKey, nil
	}
	return "", nil
}

func (c *CloudAction) createToken(gt *api_model.GetUserToken) string {
	fullStr := fmt.Sprintf("%s-%s-%s-%d-%d", gt.Body.EID, c.RegionTag, gt.Body.Range, gt.Body.ValidityPeriod, int(time.Now().Unix()))
	h := md5.New()
	h.Write([]byte(fullStr))
	return hex.EncodeToString(h.Sum(nil))
}
