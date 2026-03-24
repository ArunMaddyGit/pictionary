import { useEffect, useRef } from "react";

export type GuessMessage = {
  playerName: string;
  text: string;
  isCorrect: boolean;
  timestamp: string;
};

type ChatBoxProps = {
  messages: GuessMessage[];
};

export default function ChatBox({ messages }: ChatBoxProps) {
  const ref = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (ref.current) {
      ref.current.scrollTop = ref.current.scrollHeight;
    }
  }, [messages]);

  return (
    <div className="chat-box" ref={ref}>
      {messages.map((msg, index) => (
        <p key={`${msg.timestamp}-${index}`} className={`chat-message ${msg.isCorrect ? "correct" : ""}`}>
          <strong>{msg.playerName}:</strong> {msg.text}
        </p>
      ))}
    </div>
  );
}
