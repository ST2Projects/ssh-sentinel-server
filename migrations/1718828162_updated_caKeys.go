package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("k1ppz0e8q9of3wd")
		if err != nil {
			return err
		}

		// add
		new_default := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "zcycmqbe",
			"name": "default",
			"type": "bool",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {}
		}`), new_default); err != nil {
			return err
		}
		collection.Schema.AddField(new_default)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("k1ppz0e8q9of3wd")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("zcycmqbe")

		return dao.SaveCollection(collection)
	})
}
