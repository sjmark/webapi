package gerror

const (
	StatusDefeated                                 = 0     // 失败
	StatusSuccess                                  = 1     // 成功
	StatusNotInitialized                           = 1000  // 初始化异常
	StatusDataNotAvailable                         = 1001  // 数据不可用
	StatusArgumentMissing                          = 1002  // 参数缺失
	StatusInvalidArgument                          = 1003  // 无效参数
	StatusUnauthorized                             = 1004  // 未授权
	StatusDatabaseException                        = 1006  // 数据库操作异常
	StatusParseError                               = 1007  // 解析出错
	StatusDataRepeatError                          = 1010  // 重复提交数据
	StatusUnexpectedError                          = 10000 // 内部错误
)
