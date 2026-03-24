import { FormEvent, useState } from "react";
import { GamePhase } from "../../types/game";

type AnswerInputProps = {
  onSubmit: (text: string) => void;
  disabled: boolean;
  phase: GamePhase;
};

export default function AnswerInput({ onSubmit, disabled, phase }: AnswerInputProps) {
  const [text, setText] = useState("");
  const computedDisabled = disabled || phase !== "drawing";

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault();
    const trimmed = text.trim();
    if (!trimmed || computedDisabled) {
      return;
    }
    onSubmit(trimmed);
    setText("");
  };

  return (
    <form className="answer-input-wrap" onSubmit={handleSubmit}>
      <input
        className="answer-input"
        placeholder="Type here"
        disabled={computedDisabled}
        value={text}
        onChange={(e) => setText(e.target.value)}
      />
      <div className="answer-dashes">_ _ _ _ _ _ _ _ _</div>
    </form>
  );
}
