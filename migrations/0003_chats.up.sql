CREATE TABLE chat
(
    id        UUID PRIMARY KEY,
    user_1_id UUID NOT NULL,
    user_2_id UUID NOT NULL
);

ALTER TABLE chat ADD CONSTRAINT chat_unique_users UNIQUE (user_1_id, user_2_id);
ALTER TABLE chat ADD CONSTRAINT fk_chat_user_1 FOREIGN KEY(user_1_id)  REFERENCES user_tbl (id);
ALTER TABLE chat ADD CONSTRAINT fk_chat_user_2 FOREIGN KEY(user_2_id)  REFERENCES user_tbl (id);