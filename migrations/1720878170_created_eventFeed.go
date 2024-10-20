package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "wia3swzvun16e9l",
			"created": "2024-07-13 13:42:50.115Z",
			"updated": "2024-07-13 13:42:50.115Z",
			"name": "eventFeed",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "j7lu4q17",
					"name": "event",
					"type": "select",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"maxSelect": 1,
						"values": [
							"login",
							"sign",
							"fetch_ca"
						]
					}
				},
				{
					"system": false,
					"id": "vzr5txjz",
					"name": "user",
					"type": "relation",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"collectionId": "_pb_users_auth_",
						"cascadeDelete": false,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": null
					}
				},
				{
					"system": false,
					"id": "qzs8u6sc",
					"name": "eventTime",
					"type": "date",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": "",
						"max": ""
					}
				}
			],
			"indexes": [],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("wia3swzvun16e9l")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
