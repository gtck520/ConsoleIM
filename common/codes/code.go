package codes

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	ERROR_EXIST_TAG         = 10001
	ERROR_NOT_EXIST_TAG     = 10002
	ERROR_NOT_EXIST_ARTICLE = 10003

	ERROR_AUTH_CHECK_TOKEN_FAIL    = 20001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 20002
	ERROR_AUTH_TOKEN               = 20003
	ERROR_AUTH                     = 20004
	ERROR_EXISTS_USER              = 30001
	PAGE_NOT_FOUND                 = 40001

	MD5_PREFIX            = "jkfldfsf" //MD5加密前缀字符串
	TOKEN_KEY             = "X-Token"  //页面token键名
	USER_ID_Key           = "X-USERID" //页面用户ID键名
	USER_UUID_Key         = "X-UUID"   //页面UUID键名
	SUPER_ADMIN_ID uint64 = 956986     // 超级管理员账号ID

	//模拟枚举
	SENDTYPE_SERVER = 1 //服务端消息
	SENDTYPE_CLIENT = 2 //客户端消息

	MESSAGETYPE_TEXT = 1 //文本消息

)
