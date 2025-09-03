package lark

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yostardev/gf"
	"github.com/bytedance/sonic"
)

type MessageContentForTask struct {
	Schema string `json:"schema"`
	Config struct {
		UpdateMulti bool `json:"update_multi"`
		Style       struct {
			TextSize struct {
				NormalV2 struct {
					Default string `json:"default"`
					Pc      string `json:"pc"`
					Mobile  string `json:"mobile"`
				} `json:"normal_v2"`
			} `json:"text_size"`
		} `json:"style"`
	} `json:"config"`
	Body struct {
		Direction string `json:"direction"`
		Padding   string `json:"padding"`
		Elements  []struct {
			Tag       string `json:"tag"`
			Content   string `json:"content,omitempty"`
			TextAlign string `json:"text_align,omitempty"`
			TextSize  string `json:"text_size,omitempty"`
			Margin    string `json:"margin"`
			Text      *struct {
				Tag       string `json:"tag"`
				Content   string `json:"content"`
				TextSize  string `json:"text_size,omitempty"`
				TextAlign string `json:"text_align,omitempty"`
				TextColor string `json:"text_color,omitempty"`
			} `json:"text,omitempty"`
			Type      string `json:"type,omitempty"`
			Width     string `json:"width,omitempty"`
			Size      string `json:"size,omitempty"`
			Behaviors []struct {
				Type       string `json:"type"`
				DefaultUrl string `json:"default_url"`
				PcUrl      string `json:"pc_url"`
				IosUrl     string `json:"ios_url"`
				AndroidUrl string `json:"android_url"`
			} `json:"behaviors,omitempty"`
		} `json:"elements"`
	} `json:"body"`
	Header struct {
		Title struct {
			Tag     string `json:"tag"`
			Content string `json:"content"`
		} `json:"title"`
		Subtitle struct {
			Tag     string `json:"tag"`
			Content string `json:"content"`
		} `json:"subtitle"`
		Template string `json:"template"`
		Padding  string `json:"padding"`
	} `json:"header"`
}

func newMessageContentForTask() *MessageContentForTask {
	defaultMsg := `{"schema":"2.0","config":{"update_multi":true,"style":{"text_size":{"normal_v2":{"default":"normal","pc":"normal","mobile":"heading"}}}},"body":{"direction":"vertical","padding":"12px 12px 12px 12px","elements":[{"tag":"markdown","content":"","text_align":"left","text_size":"normal_v2","margin":"0px 0px 0px 0px"},{"tag":"hr","margin":"0px 0px 0px 0px"},{"tag":"div","text":{"tag":"plain_text","content":"任务进度","text_size":"heading","text_align":"left","text_color":"default"},"margin":"0px 0px 0px 0px"},{"tag":"markdown","content":"","text_align":"left","text_size":"normal_v2","margin":"0px 0px 0px 0px"},{"tag":"button","text":{"tag":"plain_text","content":"查看详情"},"type":"primary_filled","width":"fill","size":"large","behaviors":[{"type":"open_url","default_url":"","pc_url":"","ios_url":"","android_url":""}],"margin":"0px 0px 0px 0px"}]},"header":{"title":{"tag":"plain_text","content":""},"subtitle":{"tag":"plain_text","content":""},"template":"","padding":"12px 12px 12px 12px"}}`
	var mc MessageContentForTask
	_ = sonic.Unmarshal(gf.StringToBytes(defaultMsg), &mc)
	return &mc
}

func (mc *MessageContentForTask) setTaskName(taskName string) *MessageContentForTask {
	mc.Header.Title.Content = gf.StringJoin(taskName)
	return mc
}

func (mc *MessageContentForTask) setStatus(status string) *MessageContentForTask {
	switch status {
	case "running":
		mc.Header.Template = "blue"
		mc.Header.Subtitle.Content = "状态：执行中"
	case "success":
		mc.Header.Template = "green"
		mc.Header.Subtitle.Content = "状态：执行成功"
	case "failed":
		mc.Header.Template = "red"
		mc.Header.Subtitle.Content = "状态：执行失败"
	}
	return mc
}

func (mc *MessageContentForTask) setBasicInfo(startTime, endTime *time.Time, createdUser string) *MessageContentForTask {
	var contentList []string
	if startTime != nil {
		contentList = append(contentList, gf.StringJoin("**开始时间**: ", startTime.Format("2006-01-02 15:04:05")))
	}

	if endTime != nil {
		contentList = append(contentList, gf.StringJoin("**结束时间**: ", endTime.Format("2006-01-02 15:04:05")))
	}

	contentList = append(contentList, gf.StringJoin("**发起人**: ", createdUser))

	mc.Body.Elements[0].Content = strings.Join(contentList, "\n")

	return mc
}

func (mc *MessageContentForTask) addParameter(key, value string) *MessageContentForTask {
	if mc.Body.Elements[0].Content != "" {
		mc.Body.Elements[0].Content += "\n"
	}

	mc.Body.Elements[0].Content += fmt.Sprintf("**%s**: %s", key, value)
	return mc
}

func (mc *MessageContentForTask) setTaskID(taskID uint64) *MessageContentForTask {
	mc.Body.Elements[4].Behaviors[0].DefaultUrl = fmt.Sprintf(superLinkUrlFmt, taskID)
	return mc
}

func (mc *MessageContentForTask) addSubtasks(deep int, taskName, status string, startTime, endTime *time.Time, ignoreFailed bool) *MessageContentForTask {
	if mc.Body.Elements[3].Content != "" {
		mc.Body.Elements[3].Content += "\n"
	}

	switch status {
	case "success":
		mc.Body.Elements[3].Content += fmt.Sprintf("%"+strconv.Itoa(deep*4)+"s- %s (状态：<font color='green'>成功</font>) -- <font color='grey'>耗时： %.2f秒</font>", "", taskName, float64(endTime.Sub(*startTime))/1000/1000/1000)
	case "failed":
		if ignoreFailed {
			mc.Body.Elements[3].Content += fmt.Sprintf("%"+strconv.Itoa(deep*4)+"s- %s (状态：<font color='red'>失败（可忽略）</font>) -- <font color='grey'>耗时： %.2f秒</font>", "", taskName, float64(endTime.Sub(*startTime))/1000/1000/1000)
		} else {
			mc.Body.Elements[3].Content += fmt.Sprintf("%"+strconv.Itoa(deep*4)+"s- %s (状态：<font color='red'>失败</font>) -- <font color='grey'>耗时： %.2f秒</font>", "", taskName, float64(endTime.Sub(*startTime))/1000/1000/1000)
		}
	case "running":
		mc.Body.Elements[3].Content += fmt.Sprintf("%"+strconv.Itoa(deep*4)+"s- %s (状态：<font color='blue'>执行中</font>)", "", taskName)
	case "waiting":
		mc.Body.Elements[3].Content += fmt.Sprintf("%"+strconv.Itoa(deep*4)+"s- %s (状态：<font color='grey'>等待执行</font>)", "", taskName)
	case "abandon":
		mc.Body.Elements[3].Content += fmt.Sprintf("%"+strconv.Itoa(deep*4)+"s- %s (状态：<font color='orange'>放弃执行</font>)", "", taskName)
	}

	return mc
}

func (mc *MessageContentForTask) tooManySubtasks() *MessageContentForTask {
	mc.Body.Elements[3].Content = "- 子任务过多，请点击下方按钮查看详情。"

	return mc
}

func (mc *MessageContentForTask) string() string {
	b, _ := json.Marshal(mc)
	return gf.BytesToString(b)
}
