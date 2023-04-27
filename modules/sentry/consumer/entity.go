package consumer

import (
	"encoding/json"
	"github.com/bellis-daemon/bellis/common/models"
	"github.com/bellis-daemon/bellis/common/redistream"
	"github.com/bellis-daemon/bellis/modules/sentry/factory"
	"github.com/minoic/glgf"
	"github.com/spf13/cast"
	"time"
)

func registerEntityUpdate() {
	redistream.Instance().Register("EntityUpdate", func(message *redistream.Message) error {
		entity, err := factory.GetEntity(cast.ToString(message.Values["EntityID"]))
		if err != nil {
			return nil
		}
		var options models.Application
		err = json.Unmarshal([]byte(cast.ToString(message.Values["Entity"])), &entity)
		if err != nil {
			return err
		}
		err = entity.UpdateOptions(options)
		if err != nil {
			return err
		}
		return nil
	})
}

func registerEntityClaim() {
	redistream.Instance().Register("EntityClaim", func(message *redistream.Message) error {
		glgf.Debug(message.Values)
		ddl, err := time.Parse(time.RFC3339, cast.ToString(message.Values["Deadline"]))
		if err != nil {
			return err
		}
		if ddl.Before(time.Now()) {
			return nil
		}
		var entity models.Application
		err = json.Unmarshal([]byte(cast.ToString(message.Values["Entity"])), &entity)
		if err != nil {
			return err
		}
		err = factory.RunEntity(cast.ToString(message.Values["EntityID"]), ddl, entity)
		if err != nil {
			return err
		}
		return nil
	})
}

func init() {
	//registerEntityUpdate()
	registerEntityClaim()
}
