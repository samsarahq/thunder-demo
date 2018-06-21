CREATE TABLE repos (
  id            BIGINT NOT NULL PRIMARY KEY,
  full_name     VARCHAR(255) NOT NULL,
  api_json BLOB NOT NULL
);
CREATE TABLE events (
  at_ms    BIGINT NOT NULL,
  repo_id  BIGINT NOT NULL,
  event_id VARCHAR(255) NOT NULL,
  api_json BLOB NOT NULL,

  PRIMARY KEY(at_ms, repo_id, event_id),
  CONSTRAINT fk_repo_id FOREIGN KEY (repo_id) REFERENCES repos(id)
);
