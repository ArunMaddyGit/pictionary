type ConnectionState = "connecting" | "connected" | "disconnected" | "error";

type ConnectionStatusProps = {
  status: ConnectionState;
};

const statusLabel: Record<ConnectionState, string> = {
  connecting: "Connecting",
  connected: "Connected",
  disconnected: "Disconnected",
  error: "Error"
};

export default function ConnectionStatus({ status }: ConnectionStatusProps) {
  return (
    <div className={`connection-status status-${status}`}>
      <span className="status-dot" />
      <span>{statusLabel[status]}</span>
    </div>
  );
}
