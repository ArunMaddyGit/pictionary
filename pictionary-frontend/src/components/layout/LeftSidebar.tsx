import { Player } from "../../types/game";
import AdBanner from "../game/AdBanner";
import Leaderboard from "../game/Leaderboard";

type LeftSidebarProps = {
  players: Player[];
  currentUserId: string;
};

export default function LeftSidebar({ players, currentUserId }: LeftSidebarProps) {
  return (
    <aside className="left-sidebar">
      <AdBanner size="square" />
      <Leaderboard players={players} currentUserId={currentUserId} />
    </aside>
  );
}
