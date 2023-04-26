package email

//func getSTMPClient() (*mail.SMTPClient, error) {
//	conf := configure.GetConf()
//	temp := conf.SMTP.Encryption
//	encryption := mail.EncryptionTLS
//	if temp == "SSL" {
//		encryption = mail.EncryptionSSL
//	} else if temp != "TLS" && temp != "SSL" {
//		glgf.Error("wrong SMTP encryption")
//	}
//	sv := &mail.SMTPServer{
//		// Authentication: mail.AuthPlain,
//		Encryption:     encryption,
//		Username:       conf.SMTP.Username,
//		Password:       conf.SMTP.UserPassword,
//		ConnectTimeout: 10 * time.Second,
//		SendTimeout:    20 * time.Second,
//		Host:           conf.SMTP.Host,
//		Port:           conf.SMTP.Port,
//		KeepAlive:      false,
//	}
//	cl, err := sv.Connect()
//	return cl, err
//}
