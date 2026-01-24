CREATE TABLE IF NOT EXISTS games (
  id SERIAL PRIMARY KEY,
  home_team TEXT NOT NULL,
  away_team TEXT NOT NULL,
  home_score INT NOT NULL DEFAULT 0,
  away_score INT NOT NULL DEFAULT 0,
  timestamp TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS goals (
  id SERIAL PRIMARY KEY,
  game_id INT REFERENCES games(id) ON DELETE CASCADE,
  description TEXT,
  goalscorer TEXT,
  minute TEXT,
  url TEXT,
  reddit_url TEXT,
  mirrors TEXT,
  away BOOLEAN DEFAULT false,
  home_score INT DEFAULT 0,
  away_score INT DEFAULT 0
);

ALTER TABLE goals ADD CONSTRAINT unique_goal_url UNIQUE (url);
