CREATE TABLE direct_message
(
    id                   UUID PRIMARY KEY,
    sender_id            UUID  NOT NULL,
    receiver_id          UUID  NOT NULL,
    cipher_text          BYTEA NOT NULL,
    sender_identity_id   UUID  NOT NULL,
    receiver_identity_id UUID  NOT NULL,
    created_at           TIMESTAMP WITH TIME ZONE DEFAULT now()
);

ALTER TABLE direct_message ADD CONSTRAINT fk_dm_sender  FOREIGN KEY(sender_id) REFERENCES user_tbl (id);
ALTER TABLE direct_message ADD CONSTRAINT fk_dm_receiver FOREIGN KEY(receiver_id)  REFERENCES user_tbl (id);

ALTER TABLE direct_message ADD CONSTRAINT fk_dm_sender_identity_id FOREIGN KEY(sender_identity_id) REFERENCES identity_tbl (id);
ALTER TABLE direct_message ADD CONSTRAINT fk_dm_receiver_identity_id FOREIGN KEY(receiver_identity_id) REFERENCES identity_tbl (id);