CREATE INDEX ON chat.users (username);

CREATE TABLE chat.chats (
    id UUID,
    name          text ,
    members       set<text>,
    typing_members set<text>,
-- 	messages      list<Message>,
    image_url      text,
    created_at   timestamp,
    updated_at   timestamp,
    PRIMARY KEY (id),
);