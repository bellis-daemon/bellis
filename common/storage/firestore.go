package storage

import (
	"context"
	"sync"

	firebase "firebase.google.com/go"
	"github.com/spf13/cast"
	"google.golang.org/api/option"
)

type FirebaseInterface interface {
	Config() map[string]interface{}
	ConfigGetString(key string) string
	App() *firebase.App
}

type fb struct {
	ctx  context.Context
	app  *firebase.App
	conf map[string]interface{}
}

var firebaseInstance fb
var firebaseInitOnce sync.Once

func Firebase() FirebaseInterface {
	firebaseInitOnce.Do(func() {
		firebaseInstance.ctx = context.Background()
		opt := option.WithCredentialsJSON([]byte(Secret("firebase_config")))
		var err error
		firebaseInstance.app, err = firebase.NewApp(firebaseInstance.ctx, nil, opt)
		if err != nil {
			panic(err)
		}
		cl, err := firebaseInstance.app.Firestore(firebaseInstance.ctx)
		if err != nil {
			panic(err)
		}
		doc, err := cl.Collection("backend").Doc("config").Get(firebaseInstance.ctx)
		if err != nil {
			panic(err)
		}
		firebaseInstance.conf = doc.Data()
	})
	return &firebaseInstance
}

func (this *fb) Config() map[string]interface{} {
	return this.conf
}

func (this *fb) App() *firebase.App {
	return this.app
}

func (this *fb) ConfigGetString(key string) string {
	if v, ok := this.conf[key]; ok {
		return cast.ToString(v)
	} else {
		return ""
	}
}
