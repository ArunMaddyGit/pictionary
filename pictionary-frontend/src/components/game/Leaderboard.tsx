import { Player } from "../../types/game";

type LeaderboardProps = {
  players: Player[];
  currentUserId?: string;
};

export default function Leaderboard({ players, currentUserId = "" }: LeaderboardProps) {
  const sortedPlayers = [...players].sort((a, b) => b.score - a.score);

  return (
    <section className="leaderboard">
      <h3>Leaderboard</h3>
      <ul>
        {sortedPlayers.map((player, index) => (
          <li
            key={player.id}
            className={`leaderboard-item ${player.id === currentUserId ? "current-user" : ""}`}
          >
            <span>{index + 1}</span>
            <span>{player.name}</span>
            <span>{player.score}</span>
          </li>
        ))}
      </ul>
    </section>
  );
}
