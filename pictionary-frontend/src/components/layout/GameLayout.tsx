import { RefObject } from "react";
import { DrawPayload, GamePhase, Player } from "../../types/game";
import LeftSidebar from "./LeftSidebar";
import CenterPanel from "./CenterPanel";
import RightPanel from "./RightPanel";
import { GuessMessage } from "../game/ChatBox";
import { DrawingCanvasHandle } from "../canvas/DrawingCanvas";
import "../../styles/game.css";

type GameLayoutProps = {
  players: Player[];
  currentUserId: string;
  phase: GamePhase;
  timer: number;
  maskedWord: string;
  isDrawer: boolean;
  hasGuessed: boolean;
  messages: GuessMessage[];
  onAnswerSubmit: (text: string) => void;
  canvasRef: RefObject<DrawingCanvasHandle>;
  onDraw: (payload: DrawPayload) => void;
  onClearCanvas: () => void;
  answerDisabled: boolean;
};

export default function GameLayout({
  players,
  currentUserId,
  phase,
  timer,
  maskedWord,
  isDrawer,
  hasGuessed,
  messages,
  onAnswerSubmit,
  canvasRef,
  onDraw,
  onClearCanvas,
  answerDisabled
}: GameLayoutProps) {
  return (
    <div className="game-layout">
      <LeftSidebar players={players} currentUserId={currentUserId} />
      <CenterPanel
        phase={phase}
        timer={timer}
        maskedWord={maskedWord}
        isDrawer={isDrawer}
        canvasRef={canvasRef}
        onDraw={onDraw}
        onClearCanvas={onClearCanvas}
      />
      <RightPanel
        phase={phase}
        hasGuessed={hasGuessed}
        messages={messages}
        onAnswerSubmit={onAnswerSubmit}
        answerDisabled={answerDisabled}
      />
    </div>
  );
}
