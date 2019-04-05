package queuebiz

import (
	"webapi/models"
	"webapi/common/queue"
)

func DeQueueAll(ch chan bool) {
	businessQueue.DeQueueAll(ch)
}

func QuitQueue() {
	businessQueue.Quit()
}

func MfwTaskEnQueue(t models.MfwTask) error {
	var task = &Task{}
	task.MainAreaID = t.MainAreaID
	task.MainAreaName = t.MainAreaName
	task.Title = t.Title
	task.Url = t.Url
	task.UrlSource = t.UrlSource
	task.UrlType = t.UrlType

	_, err := queue.EnQueueTask(businessQueue, task)

	if err != nil {
		return err
	}

	return nil
}

func GivenTaskEnQueue(areaID, name, title, url string, source, uType uint8) error {
	_, err := queue.EnQueueTask(businessQueue, &GivenListTask{AreaID: areaID, Title: title, Url: url, Name: name, UrlSource: source, UrlType: uType})
	if err != nil {
		return err
	}
	return nil
}
