import { useEffect, useMemo, useRef, useState } from "react";
import { Navigate, useNavigate } from "react-router-dom";
import GameLayout from "../components/layout/GameLayout";
import WordOptions from "../components/game/WordOptions";
import GameEndScreen from "../components/game/GameEndScreen";
import ConnectionStatus from "../components/game/ConnectionStatus";
import { DrawingCanvasHandle } from "../components/canvas/DrawingCanvas";
import { useGameState } from "../hooks/useGameState";
import "../styles/game.css";

export default function GamePage() {
  const navigate = useNavigate();
  const canvasRef = useRef<DrawingCanvasHandle>(null);
  const [hadConnectedOnce, setHadConnectedOnce] = useState(false);

  const playerId = sessionStorage.getItem("playerId") ?? localStorage.getItem("playerId");
  const roomId = sessionStorage.getItem("roomId") ?? localStorage.getItem("roomId");
  const wsUrl = sessionStorage.getItem("wsUrl") ?? localStorage.getItem("wsUrl");
  const missingSession = !playerId || !roomId || !wsUrl;

  const safePlayerId = playerId ?? "";
  const safeRoomId = roomId ?? "";
  const safeWsUrl = wsUrl ?? "";

  const game = useGameState({
    playerId: safePlayerId,
    roomId: safeRoomId,
    wsUrl: safeWsUrl,
    onDrawBroadcast: (payload) => {
      canvasRef.current?.drawStroke(payload.points, payload.color, payload.brushSize);
    },
    onClearCanvas: () => {
      canvasRef.current?.clearCanvas();
    }
  });

  if (missingSession) {
    return <Navigate to="/" replace />;
  }

  const currentPlayer = useMemo(
    () => game.players.find((p) => p.id === playerId) ?? null,
    [game.players, playerId]
  );
  const hasGuessed = currentPlayer?.hasGuessed ?? false;
  const answerDisabled = hasGuessed || game.isDrawer;
  const showConnectionLost = game.status === "error" || (game.status === "disconnected" && hadConnectedOnce);

  useEffect(() => {
    if (game.status === "connected") {
      setHadConnectedOnce(true);
    }
  }, [game.status]);

  return (
    <>
      <GameLayout
        players={game.players}
        currentUserId={playerId}
        phase={game.phase}
        timer={game.timer}
        maskedWord={game.maskedWord}
        isDrawer={game.isDrawer}
        hasGuessed={hasGuessed}
        messages={game.chatMessages}
        onAnswerSubmit={game.sendGuess}
        canvasRef={canvasRef}
        onDraw={game.sendDraw}
        onClearCanvas={game.sendClearCanvas}
        answerDisabled={answerDisabled}
      />
      <ConnectionStatus status={game.status} />

      {game.phase === "choosing_word" && game.isDrawer ? (
        <WordOptions words={game.wordOptions} onSelect={game.sendSelectWord} />
      ) : null}

      {showConnectionLost ? (
        <div className="connection-lost-overlay">
          <div className="connection-lost-modal">Connection lost. Please refresh.</div>
        </div>
      ) : null}

      {game.gameEnded ? <GameEndScreen leaderboard={game.leaderboard} onPlayAgain={() => navigate("/")} /> : null}
    </>
  );
}
