package realname

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

// VerifyRealName 实名认证三要素核验方法（姓名+身份证号+手机号）
// 参数说明：
//
//	secretId: 云市场分配的密钥Id
//	secretKey: 云市场分配的密钥Key
//	name: 姓名
//	idCard: 身份证号
//	mobile: 手机号
//
// 返回值：
//
//	string: 接口返回的原始响应内容
//	error: 错误信息（参数错误/请求错误/网络错误等）
func VerifyRealName(secretId, secretKey, name, idCard, mobile string) (string, error) {
	// 1. 参数校验
	if secretId == "" || secretKey == "" {
		return "", fmt.Errorf("密钥ID和密钥Key不能为空")
	}
	if name == "" {
		return "", fmt.Errorf("姓名不能为空")
	}
	if idCard == "" {
		return "", fmt.Errorf("身份证号不能为空")
	}
	if mobile == "" {
		return "", fmt.Errorf("手机号不能为空")
	}

	// 2. 构建请求体参数
	bodyParams := map[string]string{
		"idcard": idCard,
		"mobile": mobile,
		"name":   name,
	}
	bodyParamStr := urlencode(bodyParams)

	// 3. 构建请求URL
	reqURL := "https://ap-shanghai.cloudmarket-apigw.com/service-kbq119nt/web/interface/yyssysyz"

	// 4. 发送请求（若因签名时间超时，则尝试用服务端 Date 头重试一次）
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	body, status, respDate, err := doVerifyRealNameOnce(client, secretId, secretKey, reqURL, bodyParamStr, time.Now().UTC())
	if err != nil {
		return "", err
	}
	if status >= 200 && status < 300 {
		return body, nil
	}

	if status == 403 && strings.Contains(body, "签名时间") && strings.Contains(body, "超时") && strings.TrimSpace(respDate) != "" {
		if t, parseErr := time.Parse(http.TimeFormat, strings.TrimSpace(respDate)); parseErr == nil {
			body2, status2, _, err2 := doVerifyRealNameOnce(client, secretId, secretKey, reqURL, bodyParamStr, t.Add(1*time.Second))
			if err2 != nil {
				return "", err2
			}
			if status2 >= 200 && status2 < 300 {
				return body2, nil
			}
			return "", fmt.Errorf("实名认证接口HTTP状态码异常: %d, body=%s", status2, body2)
		}
	}

	return "", fmt.Errorf("实名认证接口HTTP状态码异常: %d, body=%s", status, body)
}

func doVerifyRealNameOnce(client *http.Client, secretId, secretKey, reqURL, bodyParamStr string, now time.Time) (body string, statusCode int, respDate string, err error) {
	auth, datetime, err := calcAuthorizationAt(secretId, secretKey, now)
	if err != nil {
		return "", 0, "", fmt.Errorf("计算签名失败: %v", err)
	}
	reqID := uuid.New().String()

	request, err := http.NewRequest("POST", reqURL, strings.NewReader(bodyParamStr))
	if err != nil {
		return "", 0, "", fmt.Errorf("构建请求失败: %v", err)
	}
	request.Header.Set("Authorization", auth)
	request.Header.Set("request-id", reqID)
	request.Header.Set("X-Date", datetime)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		return "", 0, "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", 0, "", fmt.Errorf("读取响应内容失败: %v", err)
	}
	return string(bodyBytes), response.StatusCode, response.Header.Get("Date"), nil
}

func calcAuthorizationAt(secretId string, secretKey string, now time.Time) (auth string, datetime string, err error) {
	datetime = now.UTC().Format(http.TimeFormat)
	signStr := fmt.Sprintf("x-date: %s", datetime)

	// hmac-sha1 加密
	mac := hmac.New(sha1.New, []byte(secretKey))
	mac.Write([]byte(signStr))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	auth = fmt.Sprintf("{\"id\":\"%s\", \"x-date\":\"%s\", \"signature\":\"%s\"}",
		secretId, datetime, sign)

	return auth, datetime, nil
}

// urlencode URL编码（原方法保留，作为内部辅助函数）
func urlencode(params map[string]string) string {
	var p = url.Values{}
	for k, v := range params {
		p.Add(k, v)
	}
	return p.Encode()
}

func ParseVerifyResult(raw string) (bool, string) {
	if raw == "" {
		return false, "实名认证返回为空"
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		if strings.Contains(raw, "一致") || strings.Contains(raw, "匹配") {
			return true, "实名认证通过"
		}
		return false, "实名认证返回非JSON: " + raw
	}

	if v, ok := m["success"].(bool); ok && v {
		return true, "实名认证通过"
	}
	if v, ok := m["res"].(string); ok && (v == "1" || strings.EqualFold(v, "true")) {
		return true, "实名认证通过"
	}
	if v, ok := m["res"].(float64); ok && v == 1 {
		return true, "实名认证通过"
	}
	if v, ok := m["result"].(string); ok && (v == "1" || strings.EqualFold(v, "true")) {
		return true, "实名认证通过"
	}
	if v, ok := m["result"].(float64); ok && v == 1 {
		return true, "实名认证通过"
	}

	if v, ok := m["code"].(float64); ok && (v == 0 || v == 200) {
		return true, "实名认证通过"
	}
	if v, ok := m["code"].(string); ok && (v == "0" || v == "200") {
		return true, "实名认证通过"
	}

	if msg, ok := m["msg"].(string); ok && msg != "" {
		return false, msg
	}
	if msg, ok := m["message"].(string); ok && msg != "" {
		return false, msg
	}
	return false, raw
}

func ExtractSafeInfo(raw string) string {
	if raw == "" {
		return ""
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		if len(raw) > 200 {
			return raw[:200]
		}
		return raw
	}

	requestID, _ := m["requestid"].(string)
	if requestID == "" {
		requestID, _ = m["requestId"].(string)
	}
	message, _ := m["message"].(string)
	if message == "" {
		message, _ = m["msg"].(string)
	}

	code := ""
	if v, ok := m["code"].(string); ok {
		code = v
	} else if v, ok := m["code"].(float64); ok {
		code = fmt.Sprintf("%.0f", v)
	}

	parts := make([]string, 0, 3)
	if requestID != "" {
		parts = append(parts, "requestid="+requestID)
	}
	if code != "" {
		parts = append(parts, "code="+code)
	}
	if message != "" {
		parts = append(parts, "message="+message)
	}
	return strings.Join(parts, ", ")
}
