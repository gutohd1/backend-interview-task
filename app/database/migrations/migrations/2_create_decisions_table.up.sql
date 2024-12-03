CREATE TABLE IF NOT EXISTS decisions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    actor_id int,
    recipient_id int,
    liked BOOL,
    is_new BOOL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY unique_decision (actor_id, recipient_id),
    FOREIGN KEY (actor_id) REFERENCES Users(id),
	FOREIGN KEY (recipient_id) REFERENCES Users(id)
);

CREATE INDEX idx_actor_recipient ON decisions (actor_id, recipient_id);
CREATE INDEX idx_recipient ON decisions (recipient_id);