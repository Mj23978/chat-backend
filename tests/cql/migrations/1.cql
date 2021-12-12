CREATE TABLE chat.users (
  id UUID,
	username    text,     
	email       text,     
	display_name text,     
  password    text,     
	"token"       text,
	photo_url    text,     
	role        text,     
	last_active  text,     
	status      varint, 
	friends     set<text>,   
	created_at   timestamp,     
	updated_at   timestamp,     
  PRIMARY KEY ((username, email), id)
);

CREATE TABLE chat.user_messages (
                                    user_id INT,
                                    other_user_id INT,
                                    message_id TIMEUUID,
                                    message_author_id INT,
                                    message_body TEXT,
                                    PRIMARY KEY ((user_id), other_user_id, message_id)
) WITH CLUSTERING ORDER BY (other_user_id ASC, message_id DESC);