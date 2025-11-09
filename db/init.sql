-- schema for blooters
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
  away BOOLEAN DEFAULT false,
  home_score INT DEFAULT 0,
  away_score INT DEFAULT 0
);

INSERT INTO games (home_team, away_team, home_score, away_score) VALUES ('Arsenal', 'Chelsea', 1, 0) ON CONFLICT DO NOTHING;
INSERT INTO games (home_team, away_team, home_score, away_score) VALUES ('Arsenal', 'Sunderland', 4, 2) ON CONFLICT DO NOTHING;

-- link goals to games
DO $$ DECLARE
 g1 INT;
 g2 INT;
BEGIN
 SELECT id INTO g1 FROM games WHERE home_team='Arsenal' AND away_team='Chelsea' LIMIT 1;
 SELECT id INTO g2 FROM games WHERE home_team='Arsenal' AND away_team='Sunderland' LIMIT 1;

 IF g1 IS NOT NULL THEN
   INSERT INTO goals (game_id, description, goalscorer, minute, url, away, home_score, away_score)
   SELECT g1, 'Arsenal [1]-0 Chelsea - Mikel Merino 90+7''', 'Mikel Merino', '90+7''', 'https://www.youtube.com/watch?v=wnG96bon5IQ', false, 1, 0
   WHERE NOT EXISTS (SELECT 1 FROM goals WHERE game_id=g1 AND goalscorer='Mikel Merino');
 END IF;

 IF g2 IS NOT NULL THEN
   INSERT INTO goals (game_id, description, goalscorer, minute, url, away, home_score, away_score)
   SELECT g2, 'Arsenal [4]-2 Sunderland - Leandro Trossard 23'''' (Great Goal)', 'Leandro Trossard', '23''', 'https://www.youtube.com/watch?v=BRZbmT4yMWI', false, 4, 2
   WHERE NOT EXISTS (SELECT 1 FROM goals WHERE game_id=g2 AND goalscorer='Leandro Trossard');
 END IF;
END$$;