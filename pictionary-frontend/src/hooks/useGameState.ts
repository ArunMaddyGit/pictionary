import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import {
  DrawPayload,
  GamePhase,
  GuessMessagePayload,
  Player,
  PlayerScore
} from "../types/game";
import { GuessMessage } from "../components/game/ChatBox";
import { useWebSocket } from "./useWebSocket";
import clockTickingSound from "../sounds/clockticking.mp3";

type UseGameStateParams = {
  playerId: string;
  roomId: string;
  wsUrl: string;
  onDrawBroadcast?: (payload: DrawPayload) => void;
  onClearCanvas?: () => void;
};

export function useGameState({
  playerId,
  roomId,
  wsUrl,
  onDrawBroadcast,
  onClearCanvas
}: UseGameStateParams) {
  const [players, setPlayers] = useState<Player[]>([]);
  const [phase, setPhase] = useState<GamePhase>("waiting");
  const [timer, setTimer] = useState(60);
  const [maskedWord, setMaskedWord] = useState("");
  const [isDrawer, setIsDrawer] = useState(false);
  const [wordOptions, setWordOptions] = useState<string[]>([]);
  const [chatMessages, setChatMessages] = useState<GuessMessage[]>([]);
  const [leaderboard, setLeaderboard] = useState<PlayerScore[]>([]);
  const [gameEnded, setGameEnded] = useState(false);
  const [round, setRound] = useState(1);
  const tickAudioRef = useRef<HTMLAudioElement | null>(null);
  const tickingActiveRef = useRef(false);

  const handleMessage = useCallback(
    (msg: { type: string; payload: unknown }) => {
      const payload = msg.payload as Record<string, unknown>;

      switch (msg.type) {
        case "ROOM_STATE": {
          const nextPlayers = (payload.players as Player[]) ?? [];
          setPlayers(nextPlayers);
          const me = nextPlayers.find((p) => p.id === playerId);
          if (me) {
            setIsDrawer(Boolean(me.isDrawer));
          }
          setRound((payload.round as number) ?? 1);
          if (typeof payload.timer === "number") {
            setTimer(payload.timer);
          }
          break;
        }
        case "TURN_START": {
          const drawer = Boolean(payload.isDrawer);
          setIsDrawer(drawer);
          setMaskedWord(String(payload.maskedWord ?? ""));
          setWordOptions([]);
          setPhase("drawing");
          onClearCanvas?.();
          break;
        }
        case "WORD_OPTIONS": {
          setWordOptions((payload.words as string[]) ?? []);
          setIsDrawer(true);
          setPhase("choosing_word");
          break;
        }
        case "CORRECT_GUESS": {
          const scorerId = String(payload.playerId ?? "");
          const score = Number(payload.score ?? 0);
          setPlayers((prev) => prev.map((p) => (p.id === scorerId ? { ...p, score, hasGuessed: true } : p)));
          break;
        }
        case "GUESS_MESSAGE": {
          const guess = payload as unknown as GuessMessagePayload;
          const isCorrect = guess.text?.includes("guessed correctly") ?? false;
          setChatMessages((prev) => [
            ...prev,
            {
              playerName: guess.playerName ?? "",
              text: guess.text ?? "",
              isCorrect,
              timestamp: String(Date.now())
            }
          ]);
          break;
        }
        case "ROUND_END": {
          const word = String(payload.word ?? "");
          setPhase("reveal");
          setMaskedWord(word);
          onClearCanvas?.();
          setChatMessages((prev) => [
            ...prev,
            {
              playerName: "System",
              text: `Word was: ${word}`,
              isCorrect: false,
              timestamp: String(Date.now())
            }
          ]);
          break;
        }
        case "GAME_END": {
          setLeaderboard((payload.leaderboard as PlayerScore[]) ?? []);
          setGameEnded(true);
          break;
        }
        case "DRAW_BROADCAST": {
          onDrawBroadcast?.(payload as unknown as DrawPayload);
          break;
        }
        case "CLEAR_CANVAS": {
          onClearCanvas?.();
          break;
        }
        default:
          break;
      }
    },
    [onDrawBroadcast, onClearCanvas, playerId]
  );

  const { status, sendMessage, disconnect } = useWebSocket({
    url: wsUrl || null,
    onMessage: handleMessage
  });

  useEffect(() => {
    if (phase !== "drawing") {
      return;
    }
    const id = window.setInterval(() => {
      setTimer((prev) => (prev > 0 ? prev - 1 : 0));
    }, 1000);
    return () => window.clearInterval(id);
  }, [phase]);

  useEffect(() => {
    const shouldTick = phase === "drawing" && timer <= 5 && timer > 0;
    if (!shouldTick) {
      if (tickAudioRef.current && tickingActiveRef.current) {
        tickAudioRef.current.pause();
        tickAudioRef.current.currentTime = 0;
      }
      tickingActiveRef.current = false;
      return;
    }

    if (typeof Audio === "undefined") {
      return;
    }
    if (!tickAudioRef.current) {
      tickAudioRef.current = new Audio(clockTickingSound);
      tickAudioRef.current.preload = "auto";
      tickAudioRef.current.loop = true;
    }
    if (tickingActiveRef.current) {
      return;
    }

    tickingActiveRef.current = true;
    try {
      const playResult = tickAudioRef.current.play();
      if (playResult && typeof (playResult as Promise<void>).catch === "function") {
        void (playResult as Promise<void>).catch(() => {
          tickingActiveRef.current = false;
          // Ignore autoplay/path failures in MVP; game logic remains unaffected.
        });
      }
    } catch {
      tickingActiveRef.current = false;
      // Ignore audio playback failures (e.g. tests/JSDOM).
    }
  }, [phase, timer]);

  const actions = useMemo(
    () => ({
      sendDraw: (payload: DrawPayload) => sendMessage("DRAW", payload),
      sendGuess: (text: string) => sendMessage("GUESS", { text }),
      sendSelectWord: (word: string) => sendMessage("SELECT_WORD", { word }),
      sendClearCanvas: () => sendMessage("CLEAR_CANVAS", {})
    }),
    [sendMessage]
  );

  return {
    playerId,
    roomId,
    status,
    players,
    phase,
    timer,
    maskedWord,
    isDrawer,
    wordOptions,
    chatMessages,
    leaderboard,
    gameEnded,
    round,
    disconnect,
    ...actions
  };
}
