package email

import (
	"fmt"
	"time"

	"github.com/bellis-daemon/bellis/common/models"
	"github.com/matcornic/hermes"
)

func base() *hermes.Hermes {
	return &hermes.Hermes{
		Theme: new(hermes.Default),
		Product: hermes.Product{
			Name:      "Bellis Envoy",
			Link:      "https://github.com/bellis-daemon/bellis",
			Copyright: fmt.Sprintf("Copyright © 2020 - %d minoic. All rights reserved.", time.Now().Year()),
		},
	}
}

func offlineEmail(user *models.User, entity *models.Application, offlineLog *models.OfflineLog) hermes.Email {
	email := hermes.Email{
		Body: hermes.Body{
			Name: user.Email,
			Intros: []string{
				"We won't send further reminders until your app is back online or the email cools down",
			},
			Title: fmt.Sprintf("Your Entity <%s> just went offline!", entity.Name),
			Dictionary: []hermes.Entry{
				{
					Key:   "Entity name",
					Value: entity.Name,
				},
				{
					Key:   "Timezone",
					Value: user.Timezone.Location().String(),
				},
				{
					Key:   "Entity create time",
					Value: entity.CreatedAt.In(user.Timezone.Location()).Format(time.DateTime),
				},
				{
					Key:   "Offline time",
					Value: offlineLog.OfflineTime.In(user.Timezone.Location()).Format(time.DateTime),
				},
			},
			Table: hermes.Table{
				Data: [][]hermes.Entry{},
			},
			Outros: []string{
				"This should be a note worthy and validating message.",
			},
		},
	}
	if entity.Description != "" {
		email.Body.Dictionary = append(email.Body.Dictionary, hermes.Entry{
			Key:   "Entity Description",
			Value: entity.Description,
		})
	}
	for _, log := range offlineLog.SentryLogs {
		email.Body.Table.Data = append(email.Body.Table.Data, []hermes.Entry{
			{
				Key:   "Time",
				Value: log.SentryTime.In(user.Timezone.Location()).Format(time.DateTime),
			},
			{
				Key:   "Sentry",
				Value: log.SentryName,
			},
			{
				Key:   "Error",
				Value: log.ErrorMessage,
			},
		})
	}
	return email
}

func captchaEmail(captcha string) hermes.Email {
	email := hermes.Email{
		Body: hermes.Body{
			Intros: []string{
				"Bellis account manager",
				"Your captcha is " + captcha,
			},
		},
	}
	return email
}

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
