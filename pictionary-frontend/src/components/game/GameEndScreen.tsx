import { PlayerScore } from "../../types/game";

type GameEndScreenProps = {
  leaderboard: PlayerScore[];
  onPlayAgain: () => void;
};

export default function GameEndScreen({ leaderboard, onPlayAgain }: GameEndScreenProps) {
  const sorted = [...leaderboard].sort((a, b) => b.score - a.score);

  return (
    <div className="game-end-overlay">
      <div className="game-end-modal">
        <h2>Game Over!</h2>
        <ol>
          {sorted.map((entry) => (
            <li key={entry.id}>
              <span>{entry.name}</span>
              <strong>{entry.score}</strong>
            </li>
          ))}
        </ol>
        <button type="button" onClick={onPlayAgain}>
          Play Again
        </button>
      </div>
    </div>
  );
}
