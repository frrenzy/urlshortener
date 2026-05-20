CREATE TABLE IF NOT EXISTS links (
  uuid         VARCHAR(36) PRIMARY KEY,
  short_url    VARCHAR(6),
  original_url VARCHAR(256)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_links_short ON links (short_url);
