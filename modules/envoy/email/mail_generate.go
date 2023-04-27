package email

//
//import (
//	"fmt"
//	"github.com/minoic/glgf"
//	"strconv"
//)
//
//func getProd() hermes.Hermes {
//	return hermes.Hermes{
//		Theme: new(hermes.Default),
//		Product: hermes.Product{
//			Name:        "Mail",
//			Copyright:   "Copyright © 2020 - 2023 minoic. All rights reserved.",
//			TroubleText: "如果点击链接无效，请复制下列链接并在浏览器中打开：",
//		},
//	}
//}
//
//func genOfflineEmail(username string, entityUuid int64, entityName string, introduction string, offlineTime string) (string, error) {
//	h := getProd()
//	email := hermes.Email{
//		Body: hermes.Body{
//			Name: username,
//			Intros: []string{
//				introduction,
//				"我们在您的应用再次上线或邮件冷却之前不会发送更多的提醒",
//			},
//			Title: "您的应用 <" + entityName + "> 刚刚离线了",
//			Table: hermes.Table{
//				Data: [][]hermes.Entry{
//					{
//						{
//							Key:   "应用名称",
//							Value: entityName,
//						},
//						{
//							Key:   "应用 UUID",
//							Value: strconv.FormatInt(entityUuid, 10),
//						},
//						{
//							Key:   "离线时间",
//							Value: offlineTime,
//						},
//					},
//				},
//			},
//			Outros: []string{
//				"这应该是一条值得注意和验证的消息，如果有错误请联系 " + configure.GetConf().SMTP.Receiver,
//			},
//		},
//	}
//	mailBody, err := h.GenerateHTML(email)
//	return mailBody, err
//}
//
//func genForgetPasswordEmail(userName string, key string) (string, error) {
//	h := getProd()
//	email := hermes.Email{
//		Body: hermes.Body{
//			Name: userName,
//			Intros: []string{
//				" 账户管理",
//				"您正在修改密码，验证码为：" + key,
//			},
//			Outros: []string{
//				"来自",
//				"需要帮助请发邮件至 ",
//			},
//		}}
//	mailBody, err := h.GenerateHTML(email)
//	return mailBody, err
//}
//
//type Weekly map[int64]struct {
//	Name        string
//	SuccessRate float32
//}
//
//type OFL struct {
//	EntityName   string
//	StartTime    string
//	DurationTime string
//}
//
//func genWeeklyEmail(userName string, introduction string, weekly Weekly, ofl []OFL) (string, error) {
//	var (
//		avgSuccessRate  float32
//		totalEntities   int
//		healthyEntities int
//	)
//	for _, v := range weekly {
//		if v.SuccessRate >= 90 {
//			healthyEntities++
//		}
//		totalEntities++
//		avgSuccessRate += v.SuccessRate
//	}
//	avgSuccessRate /= float32(totalEntities)
//	h := getProd()
//	email := hermes.Email{
//		Body: hermes.Body{
//			Name:     userName,
//			Greeting: "这是您的每周应用在线率报告 ",
//			Intros: []string{
//				"应用离线记录：",
//			},
//			Table: hermes.Table{
//				Data: [][]hermes.Entry{},
//			},
//			Outros: []string{
//				introduction,
//				"来自",
//				"需要帮助请发邮件至 ",
//			},
//		},
//	}
//	for k, v := range weekly {
//		email.Body.Table.Data = append(email.Body.Table.Data, []hermes.Entry{
//			{Key: "UUID", Value: strconv.FormatInt(k, 10)},
//			{Key: "Name", Value: v.Name},
//			{Key: "SuccessRate", Value: strconv.FormatFloat(float64(v.SuccessRate*100), 'f', 2, 32) + "%"},
//		})
//	}
//	for i := range ofl {
//		email.Body.Intros = append(email.Body.Intros, fmt.Sprintf("%s 于 %s 离线了持续 %s", ofl[i].EntityName, ofl[i].StartTime, ofl[i].DurationTime))
//	}
//	mailBody, err := h.GenerateHTML(email)
//	return mailBody, err
//}
//
//func genAnyEmail(text string) (string, string) {
//	h := getProd()
//	email := hermes.Email{
//		Body: hermes.Body{
//			Intros: []string{
//				text,
//			},
//			Outros: []string{
//				"来自",
//				"需要帮助请发邮件至 ",
//			},
//		},
//	}
//	mailBody, err := h.GenerateHTML(email)
//	if err != nil {
//		glgf.Error(err)
//		return "", ""
//	}
//	mailText, err := h.GeneratePlainText(email)
//	if err != nil {
//		glgf.Error(err)
//		return "", ""
//	}
//	return mailBody, mailText
//}
