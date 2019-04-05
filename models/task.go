package models

type(
	MfwTask struct {
		MainAreaID,
		MainAreaName,
		Title,
		Url string      // URL area tag detail的url
		UrlSource uint8 // 1 主页 2 tags
		UrlType   uint8 // 1 列表页面 2 详情页面
	}
)
