# recaptcha multiple server solver

##可以一次使用多個captcha solver, captcha solver account 或進行提前預約的Restful Api

## /api/account/new 
1. Post
2. {Account:"",Password:""}
3. create new account
4. return apikey
## /api/account/{apikey}
1. Get
2. get acconut status
3. return account status
## /api/book
1. Post
2. {  
   BookingTime INT                    &emsp;&emsp;&emsp;&emsp;//use unix time  
   Key TEXT                           &emsp;&emsp;&emsp;&emsp;&emsp;&emsp;//user Key  
   Type TEXT                          &emsp;&emsp;&emsp;&emsp;//task type "recaptchav2","recaptchav3"  
   Number INT                         &emsp;&emsp;&emsp;&emsp;//the number of captcha you want to book  
   Cumulative BOOL                    &emsp;&emsp;&emsp;&emsp;//if true you can increase the booking captcha number by post many time.  
   TaskInfo {KEY:VALUE}               &emsp;&emsp;&emsp;&emsp;//task parameter  
   Rate FLOAT                         &emsp;&emsp;&emsp;&emsp;//not necessary(has default value),the rate of captcha number.when captcha is error,you can change new answer fast  
   }
3. booking a task
4. return taskid
5. TaskInfo :  

Recaptcha v2  

|Name|Type|Necessary|Description|
|:-----:|:-----:|:-----:|:-----:|
|url|String|Yes|Address of a target web page. Can be located anywhere on the web site, even in a member area. Our workers don't navigate there but simulate the visit instead.|
|websiteKey|String|Yes|Recaptcha website key. Learn how to find it in this article.|
|recaptchaDataSValue|String|No|Value of 'data-s' parameter. Applies only to Recaptchas on Google web sites.|
|isInvisible|Boolean|No|Specify whether or not Recaptcha is invisible. This will render an appropriate widget for our workers.|
Recaptcha v3

|Name|Type|Necessary|Description|
|:-----:|:-----:|:-----:|:-----:|
|url|String|Yes|Address of a target web page. Can be located anywhere on the web site, even in a member area. Our workers don't navigate there but simulate the visit instead.|
|websiteKey|String|Yes|Recaptcha website key. Learn how to find it in this article.|
|minScore|Double|Yes|Filters workers with a particular score. It can have one of the following values: 0.3 0.7 0.9
|pageAction|String|No|Recaptcha's "action" value. Website owners use this parameter to define what users are doing on the page. Example: grecaptcha.execute('site_key', {action:'login_test'})|
|isEnterprise|Boolean|No|Set this flag to "true" if you need this V3 solved with Enterprise API. Default value is "false" and Recaptcha is solved with non-enterprise API. Can be determined by a javascript call like in the following example: grecaptcha.enterprise.execute('site_key', {..})|




## /api/task/{taskid}
1. Get
2. get captcha result
3. return result


## /api/report
1. Post
2. {"Answer":"","Correct":True}
3. report answer to statis
4. return apikey





