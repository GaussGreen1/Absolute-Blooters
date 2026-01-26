import { useState, useEffect } from 'react';
import { API_BASE_URL } from "./config";
import { useTheme } from './contexts/ThemeContext';
import './App.css';

interface Goal {
  id: number;
  game_id: number;
  description: string;
  home_team: string;
  away_team: string;
  goalscorer: string;
  minute: string;
  url: string;
  reddit_url: string;
  mirrors: string;
  home_score: number;
  away_score: number;
  away: boolean;
}

interface Game {
  id: number;
  home_team: string;
  away_team: string;
  home_score: number;
  away_score: number;
  goals: Goal[];
  timestamp: string;
}

interface GamesResponse {
  games: Game[];
  status: number;
}

function App() {
  const { theme, toggleTheme } = useTheme();
  const [gamesData, setGamesData] = useState<Game[] | null>(null);
  const [loadingGames, setLoadingGames] = useState<boolean>(true);
  const [gamesError, setGamesError] = useState<string | null>(null);

  useEffect(() => {
    console.log('Fetching games...');
    fetch(`${API_BASE_URL}/api/games`)
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.json();
      })
      .then((data: GamesResponse) => {
        setGamesData(data.games);
        setLoadingGames(false);
      })
      .catch((err) => {
        console.error("Fetch games error:", err);
        setGamesError(err.message);
        setLoadingGames(false);
      });
  }, []);

  return (
    <div className="app-container">
      <div className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '1rem' }}>
          <h1>Absolute Blooters</h1>
          <button
            className="theme-toggle"
            onClick={toggleTheme}
            aria-label={`Switch to ${theme === 'light' ? 'dark' : 'light'} mode`}
          >
            {theme === 'light' ? '🌙' : '☀️'}
          </button>
        </div>

        {/* Games Section */}
        <section className="section">

          {loadingGames && (
            <div className="status loading">
              <div className="spinner"></div>
              <p>Loading games...</p>
            </div>
          )}

          {gamesError && (
            <div className="status error">
              <p>❌ Error loading games: {gamesError}</p>
            </div>
          )}

          {gamesData && !loadingGames && !gamesError && (
            <div className="games-container">
              {gamesData.length === 0 ? (
                <p>No games available.</p>
              ) : (
                gamesData.map((game) => (
                  <div key={game.id} className="game-card">
                    <div className="game-header">
                      <h3>{game.home_team} vs {game.away_team}</h3>
                      <p className="score">{game.home_score} - {game.away_score}</p>
                      <p className="timestamp">{new Date(game.timestamp).toLocaleString('it-IT', { 
                        year: 'numeric', 
                        month: '2-digit', 
                        day: '2-digit', 
                        hour: '2-digit', 
                        minute: '2-digit',
                        hour12: false 
                      })}</p>
                    </div>
                    {game.goals.length > 0 && (
                      <div className="goals">
                        <h4>Goals:</h4>
                        <ul>
                          {game.goals.map((goal) => (
                            <li key={goal.id}>
                              {goal.url ? (
                                <div className="goal-links">
                                  <a href={goal.url} target="_blank" rel="noopener noreferrer" className="goal-link">
                                    <strong>{goal.goalscorer}</strong> ({goal.minute}') - {goal.home_score}-{goal.away_score}
                                    {goal.away ? ` (${goal.away_team})` : ` (${goal.home_team})`}
                                    <span className="watch-text"> ▶ Watch</span>
                                  </a>
                                  {goal.mirrors && (
                                    <a href={goal.mirrors} target="_blank" rel="noopener noreferrer" className="mirror-link">
                                      🔄 Mirrors
                                    </a>
                                  )}
                                </div>
                              ) : (
                                <span>
                                  <strong>{goal.goalscorer}</strong> ({goal.minute}') - {goal.home_score}-{goal.away_score}
                                  {goal.away ? ` (${goal.away_team})` : ` (${goal.home_team})`}
                                </span>
                              )}
                            </li>
                          ))}
                        </ul>
                      </div>
                    )}
                  </div>
                ))
              )}
            </div>
          )}
        </section>
      </div>
    </div>
  );
}

export default App;
