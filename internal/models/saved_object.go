package models

type SavedObject struct {
	ID      string
	SpaceID string
	/*
		There is no well-defined API for managing dashboards in Kibana, so we interact
		directly with the database object when interacting with the resource
	*/
	Attributes map[string]interface{}
}
