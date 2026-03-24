import { GamePhase } from "../../types/game";
import AnswerInput from "../game/AnswerInput";
import ChatBox, { GuessMessage } from "../game/ChatBox";

type RightPanelProps = {
  phase: GamePhase;
  hasGuessed: boolean;
  messages: GuessMessage[];
  onAnswerSubmit: (text: string) => void;
  answerDisabled: boolean;
};

export default function RightPanel({ phase, hasGuessed, messages, onAnswerSubmit, answerDisabled }: RightPanelProps) {
  return (
    <aside className="right-panel">
      <AnswerInput onSubmit={onAnswerSubmit} disabled={answerDisabled || hasGuessed} phase={phase} />
      <ChatBox messages={messages} />
    </aside>
  );
}
