export type GamePhase = "waiting" | "choosing_word" | "drawing" | "reveal";

export type RoomStatus = "waiting" | "playing" | "finished";

export type Player = {
  id: string;
  name: string;
  score: number;
  isDrawer: boolean;
  hasGuessed: boolean;
};

export type Room = {
  id: string;
  players: Player[];
  status: RoomStatus;
  round: number;
  maxRounds: number;
  currentDrawerIndex: number;
  phase: GamePhase;
};

export type Message<T> = {
  type: string;
  payload: T;
};

export type DrawPayload = {
  points: [number, number][];
  color: string;
  brushSize: number;
};

export type GuessPayload = {
  text: string;
};

export type SelectWordPayload = {
  word: string;
};

export type RoomStatePayload = {
  players: Player[];
  round: number;
  drawerId: string;
  timer: number;
};

export type DrawBroadcastPayload = {
  points: [number, number][];
  color: string;
  brushSize: number;
};

export type TurnStartPayload = {
  drawerId: string;
  maskedWord: string;
  isDrawer: boolean;
};

export type WordOptionsPayload = {
  words: string[];
};

export type CorrectGuessPayload = {
  playerId: string;
  score: number;
};

export type RoundEndPayload = {
  word: string;
};

export type PlayerScore = {
  id: string;
  name: string;
  score: number;
};

export type GameEndPayload = {
  leaderboard: PlayerScore[];
};

export type GuessMessagePayload = {
  playerId: string;
  playerName: string;
  text: string;
};
