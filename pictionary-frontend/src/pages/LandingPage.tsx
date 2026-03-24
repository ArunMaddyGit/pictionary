import { FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";

type JoinResponse = {
  playerId: string;
  roomId: string;
  wsUrl: string;
};

export default function LandingPage() {
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [isHovering, setIsHovering] = useState(false);

  const onSubmit = async (event: FormEvent) => {
    event.preventDefault();
    const trimmedName = name.trim();
    if (!trimmedName) {
      setError("Name required");
      return;
    }
    if (trimmedName.length > 20) {
      setError("Name too long");
      return;
    }

    setLoading(true);
    setError("");
    try {
      const response = await fetch("/api/join", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: trimmedName })
      });
      if (!response.ok) {
        if (response.status === 400) {
          const message = (await response.text()).toLowerCase();
          if (message.includes("name too long")) {
            setError("Name too long");
            return;
          }
          if (message.includes("name is required")) {
            setError("Name required");
            return;
          }
        }
        setError("Server unavailable");
        return;
      }
      const data: JoinResponse = await response.json();
      sessionStorage.setItem("playerId", data.playerId);
      sessionStorage.setItem("roomId", data.roomId);
      sessionStorage.setItem("wsUrl", data.wsUrl);
      navigate("/game");
    } catch {
      setError("Server unavailable");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      style={{
        minHeight: "100vh",
        backgroundColor: "#008080",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        padding: "24px"
      }}
    >
      <form
        onSubmit={onSubmit}
        style={{
          backgroundColor: "rgba(255,255,255,0.12)",
          borderRadius: "12px",
          padding: "24px",
          width: "100%",
          maxWidth: "360px",
          display: "flex",
          flexDirection: "column",
          gap: "12px",
          color: "#ffffff"
        }}
      >
        <h1 style={{ margin: 0, textAlign: "center" }}>Pictionary</h1>
        <label htmlFor="name">Your name</label>
        <input
          id="name"
          aria-label="Player name"
          value={name}
          maxLength={20}
          onChange={(e) => setName(e.target.value)}
          style={{
            border: "none",
            borderRadius: "8px",
            padding: "10px 12px"
          }}
        />
        <p style={{ margin: 0, opacity: 0.9, fontSize: "0.9rem", textAlign: "right" }}>{name.length}/20</p>
        <button
          type="submit"
          disabled={loading}
          onMouseEnter={() => setIsHovering(true)}
          onMouseLeave={() => setIsHovering(false)}
          style={{
            border: "none",
            borderRadius: "8px",
            padding: "10px 12px",
            backgroundColor: loading ? "#9aa7a7" : "#ffffff",
            color: "#005c5c",
            fontWeight: 600,
            cursor: loading ? "not-allowed" : "pointer",
            transition: "transform 180ms ease, box-shadow 180ms ease",
            transform: !loading && isHovering ? "translateY(-1px) scale(1.02)" : "none",
            boxShadow: !loading && isHovering ? "0 8px 18px rgba(0,0,0,0.2)" : "none"
          }}
        >
          {loading ? "Joining..." : "Play Now"}
        </button>
        {error ? <p role="alert">{error}</p> : null}
      </form>
    </div>
  );
}
