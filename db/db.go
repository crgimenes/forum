package db

type Topic struct {
	ID          int      `db:"id" json:"id"`
	ParentID    int      `db:"parent_id" json:"parent_id,omitempty"`
	Level       int      `db:"level" json:"level,omitempty"`
	ZOrder      int      `db:"zorder" json:"zorder,omitempty"`
	Slug        string   `db:"slug" json:"slug,omitempty"`
	Title       string   `db:"title" json:"title"`
	Description string   `db:"description" json:"description,omitempty"`
	CreatedAt   string   `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt   string   `db:"updated_at" json:"updated_at,omitempty"`
	Topics      []Topic  `db:"-" json:"subtopic,omitempty"`
	Threads     []Thread `db:"-" json:"threads,omitempty"`
}

type Thread struct {
	ID          int    `db:"id" json:"id"`
	TopicID     int    `db:"topic_id" json:"topic_id"`
	Slug        string `db:"slug" json:"slug,omitempty"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description,omitempty"`
	CreatedAt   string `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt   string `db:"updated_at" json:"updated_at,omitempty"`
}

type DB interface {
	Open() error
	Close() error
	GetTopicID(slug []string) (int, error)
	ListTopics(parent_id int) ([]Topic, error)
	ListThreadsPerTopic(topicID int) ([]Thread, error)
}
