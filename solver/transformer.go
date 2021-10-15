package solver



//transfor captcha information to server params

type RecaptchaV2 struct {
	Url string `json:"url"`
	WebsiteKey string `json:"websiteKey"`
	RecaptchaDataSValue string `json:"recaptchaDataSValue"`
	IsInvisible bool `json:"isInvisible"`
}
func (recaptchav2 RecaptchaV2) Parse2captcha() map[string]interface{} {
	output := map[string]interface{}{}
	output["method"] = "userrecaptcha"
	output["pageurl"] = recaptchav2.Url
	output["googlekey"] = recaptchav2.WebsiteKey
	if recaptchav2.RecaptchaDataSValue != "" {
		output["data-s"] = recaptchav2.RecaptchaDataSValue
	}
	if recaptchav2.IsInvisible{
		output["invisible"] = "1"
	}
	return output
}
func (recaptchav2 RecaptchaV2) ParseAntiCaptcha() (map[string]interface{}) {
	output := map[string]interface{}{}
	output["type"] = "RecaptchaV2TaskProxyless"
	output["websiteURL"] = recaptchav2.Url
	output["websiteKey"] = recaptchav2.WebsiteKey
	if recaptchav2.RecaptchaDataSValue != "" {
		output["recaptchaDataSValue"] = recaptchav2.RecaptchaDataSValue
	}
	if recaptchav2.IsInvisible{
		output["isInvisible"] = true
	}
	return output
}

type RecaptchaV3 struct {
	Url string `json:"url"`
	WebsiteKey string `json:"websiteKey"`
	MinScore float64 `json:"minScore"`
	PageAction string `json:"pageAction"`
	IsEnterprise bool `json:"isEnterprise"`
}
func (recaptchav3 RecaptchaV3) Parse2captcha() (map[string]interface{}) {
	output := map[string]interface{}{}
	output["method"] = "userrecaptcha"
	output["version"] = "v3"
	output["pageurl"] = recaptchav3.Url
	output["googlekey"] = recaptchav3.WebsiteKey
	output["min_score"] = recaptchav3.MinScore
	if recaptchav3.PageAction != "" {
		output["action"] = recaptchav3.PageAction
	}
	if recaptchav3.IsEnterprise{
		output["enterprise"] = "1"
	}
	return output

}
func (recaptchav3 RecaptchaV3) ParseAntiCaptcha() (map[string]interface{}) {
	output := map[string]interface{}{}
	output["type"] = "RecaptchaV3TaskProxyless"
	output["websiteURL"] = recaptchav3.Url
	output["websiteKey"] = recaptchav3.WebsiteKey
	output["minScore"] = recaptchav3.MinScore
	if recaptchav3.PageAction != "" {
		output["pageAction"] = recaptchav3.PageAction
	}
	if recaptchav3.IsEnterprise{
		output["isEnterprise"] = true
	}
	return output

}
