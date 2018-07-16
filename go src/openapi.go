package common

import (
	"encoding/binary"
	"hash/crc32"
	"strconv"
	"strings"
)

const (
	Third_ID_SALT_KEY = "_thd_salt__key_" // 此值不能改变
)

func GetThirdId(userid uint64) string {
	var tmpByte []byte = make([]byte, 12)
	binary.LittleEndian.PutUint64(tmpByte, userid)
	binary.LittleEndian.PutUint32(tmpByte[8:], crc32.ChecksumIEEE(tmpByte[:8]))
	tmpByte[1], tmpByte[8] = tmpByte[8], tmpByte[1]
	tmpByte[3], tmpByte[9] = tmpByte[9], tmpByte[3]
	tmpByte[5], tmpByte[10] = tmpByte[10], tmpByte[5]
	tmpByte[7], tmpByte[11] = tmpByte[11], tmpByte[7]
	tmp := Rc4Encode(tmpByte)
	if strings.Contains(tmp, "_") {
		return strconv.FormatUint(userid, 10)
	}
	return tmp + "_" + strconv.FormatUint(userid%1000, 10)
}

func GetUidFromOpenId(openId string) uint64 {
	thdId := GetThdIdFromOpenId(openId)
	return GetUidFromThdId(thdId)
}

func GetUidFromThdId(thdid string) uint64 {
	datas := strings.Split(thdid, "_")
	if len(datas) == 1 {
		return StringToUint64(thdid)
	}

	tmpByte, _ := Rc4Decode(datas[0])
	if len(tmpByte) < 12 {
		return 0
	}
	tmpByte[8], tmpByte[1] = tmpByte[1], tmpByte[8]
	tmpByte[9], tmpByte[3] = tmpByte[3], tmpByte[9]
	tmpByte[10], tmpByte[5] = tmpByte[5], tmpByte[10]
	tmpByte[11], tmpByte[7] = tmpByte[7], tmpByte[11]

	return binary.LittleEndian.Uint64(tmpByte[:8])
}

func GetOpenIdFromUid(userid uint64, clientid string) string {
	thirdid := GetThirdId(userid)
	return GetOpenIdFromThdId(thirdid, clientid)
}

func GetOpenIdFromThdId(thdid, clientid string) string {
	tmp := EncMd5(thdid + clientid)
	return tmp[:5] + thdid
}

func GetThdIdFromOpenId(openid string) string {
	if len(openid) <= 5 {
		return ""
	}
	return openid[5:]
}
