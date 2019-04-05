package queuebiz

import "webapi/models"

// 入列数据
type (
	Task struct {
		MainAreaID,    // 主区域的id
		MainAreaName,  // 主区域的名称
		Title,         // 详情页面的title
		Url string      // 包含 area tag detail的url
		UrlSource uint8 // 1 主页 2 tags
		UrlType   uint8 // 1 列表页面 2 详情页面
	}
)

func (t Task) Exec() error {

	cBusiness.CommonBusiness(
		models.MfwTask{
			t.MainAreaID,
			t.MainAreaName,
			t.Title,
			t.Url,
			t.UrlSource,
			t.UrlType,
		})
	return nil
}

// 特定数据
type GivenListTask struct {
	AreaID    string
	Name      string
	Title     string
	Url       string
	UrlSource uint8 // 1 主页 2 tags
	UrlType   uint8 // 1 列表页面 2 详情页面
}

func (g GivenListTask) Exec() error {
	gBusiness.GivenBusiness(g.AreaID, g.Name, g.Title, g.Url, g.UrlSource, g.UrlType)
	return nil
}
