package email

//func sendEmail(email *mail.Email) error {
//	cl, err := getSTMPClient()
//	if err != nil {
//		return err
//	}
//	if err := email.Send(cl); err != nil {
//		return err
//	}
//	return nil
//}
//
//func SendOfflineEmail(receiveAddr string, username string, entityUuid int64, entityName string, introduction string, offlineTime string) error {
//	mailHtml, err := genOfflineEmail(username, entityUuid, entityName, introduction, offlineTime)
//	if err != nil {
//		return err
//	}
//	email := mail.NewMSG()
//	email.SetFrom(configure.GetConf().SMTP.SendFrom).
//		AddTo(receiveAddr).
//		SetSubject(configure.GetConf().ApplicationName+" 应用离线通知").
//		SetBody(mail.TextHTML, mailHtml)
//	err = sendEmail(email)
//	if err != nil {
//		glgf.Error(err)
//	}
//	return sendEmail(email)
//}
//
//func SendCaptcha(userName string, receiver string, key string) error {
//	// glg.Info(receiver)
//	mailHtml, err := genForgetPasswordEmail(userName, key)
//	if err != nil {
//		return err
//	}
//	email := mail.NewMSG()
//	email.SetFrom(configure.GetConf().SMTP.SendFrom).
//		AddTo(receiver).
//		SetSubject(configure.GetConf().ApplicationName+" 验证码").
//		SetBody(mail.TextHTML, mailHtml)
//	return sendEmail(email)
//}
//
//func SendWeeklyMail(userName string, receiver string, introduction string, weekly Weekly, ofl []OFL) error {
//	mailHtml, err := genWeeklyEmail(userName, introduction, weekly, ofl)
//	if err != nil {
//		return err
//	}
//	email := mail.NewMSG()
//	email.SetFrom(configure.GetConf().SMTP.SendFrom).
//		AddTo(receiver).
//		SetSubject(configure.GetConf().ApplicationName+" 每周汇总邮件").
//		SetBody(mail.TextHTML, mailHtml)
//	return sendEmail(email)
//}
