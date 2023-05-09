package pg

import (
	"encoding/json"
	"fmt"

	"crg.eti.br/go/forum/config"
	"crg.eti.br/go/forum/db"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Pg struct {
	db *sqlx.DB
}

func New() db.DB {
	return &Pg{}
}

func (p *Pg) Open() error {
	var err error
	p.db, err = sqlx.Open("postgres", config.DBURL)
	if err != nil {
		return err
	}

	return p.db.Ping()
}

func (p *Pg) Close() error {
	return p.db.Close()
}

func (p *Pg) GetTopicID(slug []string) (int, error) {
	var (
		id        int
		parent_id int
		err       error
	)

	for _, s := range slug {
		sqlStatement := `SELECT id FROM forum_topics WHERE slug = $1 AND parent_id = $2`
		err = p.db.QueryRow(sqlStatement, s, parent_id).Scan(&id)
		if err != nil {
			return 0, err
		}
		parent_id = id
	}

	return id, nil
}

func (p *Pg) ListTopics(parent_id int) ([]db.Topic, error) {
	sqlStatement := `WITH RECURSIVE topics_tree AS (
		SELECT 
			id, 
			parent_id, 
			zorder,
			slug,
			title, 
			description, 
			0 as level, 
			created_at
			FROM forum_topics
			WHERE parent_id = $1
		UNION ALL
		SELECT 
			ft.id, 
			ft.parent_id,
			ft.zorder,
			ft.slug,
			ft.title, 
			ft.description, 
			tt.level +1 as level, 
			ft.created_at 
			FROM forum_topics ft
			JOIN topics_tree tt ON ft.parent_id = tt.id
		)
		SELECT 
			id,
			parent_id,
			zorder,
			slug,
			title, 
			description, 
			level, 
			created_at 
		FROM topics_tree
		ORDER BY level,zorder, created_at;`
	query, err := p.db.Queryx(sqlStatement, parent_id)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var topics []db.Topic
	for query.Next() {
		var topic db.Topic
		err := query.StructScan(&topic)
		if err != nil {
			return nil, err
		}
		if topic.ParentID == 0 || topic.ParentID == parent_id {

			topic.Threads, err = p.ListThreadsPerTopic(topic.ID)
			if err != nil {
				return nil, err
			}

			topics = append(topics, topic)
			continue
		}
		for i, t := range topics {
			if t.ID == topic.ParentID {

				topic.Threads, err = p.ListThreadsPerTopic(topic.ID)
				if err != nil {
					return nil, err
				}

				topics[i].Topics = append(topics[i].Topics, topic) // add subtopic
				break
			}
		}
	}

	b, _ := json.MarshalIndent(topics, "", "  ")
	fmt.Println(string(b))

	return topics, nil
}

func (p *Pg) ListThreadsPerTopic(topicID int) ([]db.Thread, error) {
	sqlStatement := `SELECT 
	id,
	title,
	slug,
	description,
	created_at,
	updated_at
	FROM forum_threads 
	WHERE topic_id = $1
	ORDER BY zorder, updated_at DESC;`

	query, err := p.db.Queryx(sqlStatement, topicID)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	var threads []db.Thread
	for query.Next() {
		var thread db.Thread
		err := query.StructScan(&thread)
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	return threads, nil
}
