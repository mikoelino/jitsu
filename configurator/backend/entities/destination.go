package entities

import "github.com/jitsucom/jitsu/server/enrichment"

//Destination entity is stored in main storage (Firebase or Redis)
type Destination struct {
	ID               string                   `firestore:"_id" json:"_id"`
	UID              string                   `firestore:"_uid" json:"_uid"`
	Type             string                   `firestore:"_type"  json:"_type"`
	Data             interface{}              `firestore:"_formData" json:"_formData"`
	Mappings         *Mappings                `firestore:"_mappings" json:"_mappings"`
	Enrichment       []*enrichment.RuleConfig `firestore:"_enrichment" json:"_enrichment"`
	UsersRecognition *UsersRecognition        `firestore:"_users_recognition" json:"_users_recognition"`
	OnlyKeys         []string                 `firestore:"_onlyKeys" json:"_onlyKeys"`
	PrimaryKeyFields []string                 `firestore:"_primary_key_fields" json:"_primary_key_fields"`
}

//Destinations entity is stored in main storage (Firebase or Redis)
type Destinations struct {
	Destinations []*Destination `json:"destinations" firestore:"destinations"`
}

//Mappings entity is stored in main storage (Firebase or Redis)
type Mappings struct {
	KeepFields bool      `firestore:"_keepUnmappedFields" json:"_keepUnmappedFields"`
	Rules      []MapRule `firestore:"_mappings" json:"_mappings"`
}

//UsersRecognition entity is stored in main storage (Firebase or Redis)
type UsersRecognition struct {
	Enabled         bool   `firestore:"_enabled" json:"_enabled"`
	AnonymousIDNode string `firestore:"_anonymous_id_node" json:"_anonymous_id_node"`
	UserIDJSONNode  string `firestore:"_user_id_node" json:"_user_id_node"`
}

//IsEmpty returns true if mappings is empty
func (m *Mappings) IsEmpty() bool {
	return m == nil || len(m.Rules) == 0
}

//MapRule entity is stored in main storage (Firebase or Redis)
type MapRule struct {
	Action           string `firestore:"_action" json:"_action"`
	SourceField      string `firestore:"_srcField" json:"_srcField"`
	DestinationField string `firestore:"_dstField" json:"_dstField"`
}
