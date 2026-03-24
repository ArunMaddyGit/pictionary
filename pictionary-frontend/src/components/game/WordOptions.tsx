import { useEffect, useState } from "react";

type WordOptionsProps = {
  words: string[];
  onSelect: (word: string) => void;
};

export default function WordOptions({ words, onSelect }: WordOptionsProps) {
  const [visible, setVisible] = useState(true);
  const [secondsLeft, setSecondsLeft] = useState(15);

  useEffect(() => {
    setVisible(true);
    setSecondsLeft(15);
  }, [words]);

  useEffect(() => {
    if (!visible) {
      return;
    }
    if (secondsLeft <= 0) {
      if (words.length > 0) {
        onSelect(words[0]);
      }
      setVisible(false);
      return;
    }
    const id = window.setTimeout(() => {
      setSecondsLeft((prev) => prev - 1);
    }, 1000);
    return () => window.clearTimeout(id);
  }, [secondsLeft, visible, words, onSelect]);

  if (!visible || words.length === 0) {
    return null;
  }

  return (
    <div className="word-options-overlay" role="dialog" aria-label="Word options">
      <div className="word-options-modal">
        <h3>Choose a word</h3>
        <p className="word-options-hint">Pick a word to draw. You have {secondsLeft}s.</p>
        <div className="word-options-buttons">
          {words.map((word) => (
            <button
              key={word}
              type="button"
              onClick={() => {
                onSelect(word);
                setVisible(false);
              }}
            >
              {word}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}
