package protocol

type RSS3Index031 struct {
	Version        string `bson:"version"`      // "rss3.io/version/v0.3.1"
	ID             string `bson:"id"`           // EVM+ address
	DateCreated    string `bson:"date_created"` // "2021-08-15T03:00:57.449Z"
	DateUpdated    string `bson:"date_updated"` // "2022-02-21T01:56:48.880Z"
	Signature      string `bson:"signature"`
	AgentID        string `bson:"agent_id"`
	AgentSignature string `bson:"agent_signature"`
	Profile        struct {
		Name     string   `bson:"name"`
		Bio      string   `bson:"bio"`
		Avatar   []string `bson:"avatar"`
		Accounts []struct {
			ID        string   `bson:"id"`
			Tags      []string `bson:"tags"`
			Signature string   `bson:"signature"`
		} `bson:"accounts"`
	} `bson:"profile"`
	Links []struct {
		ID   string `bson:"id"`
		List string `bson:"list"`
	} `bson:"links"`
	Backlinks []struct {
		Auto bool   `bson:"auto"` // true
		ID   string `bson:"id"`
		List string `bson:"list"`
	} `bson:"backlinks"`
	Items struct {
		ListAuto string `bson:"list_auto"`
	} `bson:"items"`
	Assets struct {
		ListAuto string `bson:"list_auto"`
	} `bson:"assets"`

	CustomFiled_Pass struct {
		Assets []struct {
			ID    string `bson:"id"`
			Class string `bson:"class"`
			Order int    `bson:"order"`
		} `bson:"assets"`
	} `bson:"_pass"`
}

type RSS3AutoAssets031 struct {
	Version     string   `bson:"version"` // "rss3.io/version/v0.3.1"
	ID          string   `bson:"id"`
	DateCreated string   `bson:"date_created"`
	DateUpdated string   `bson:"date_updated"`
	Auto        bool     `bson:"auto"` // true
	List        []string `bson:"list"`
}

type RSS3Links031 struct {
	Version        string   `bson:"version"` // "rss3.io/version/v0.3.1"
	ID             string   `bson:"id"`
	DateCreated    string   `bson:"date_created"`
	DateUpdated    string   `bson:"date_updated"`
	List           []string `bson:"list"`
	Signature      string   `bson:"signature"`
	AgentID        string   `bson:"agent_id"`
	AgentSignature string   `bson:"agent_signature"`
}

type RSS3Backlinks031 struct {
	Version     string   `bson:"version"` // "rss3.io/version/v0.3.1"
	ID          string   `bson:"id"`
	DateCreated string   `bson:"date_created"`
	DateUpdated string   `bson:"date_updated"`
	Auto        bool     `bson:"auto"` // true
	List        []string `bson:"list"`
}
