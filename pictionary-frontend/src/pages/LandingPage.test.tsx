import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, it, vi, beforeEach, afterEach } from "vitest";
import LandingPage from "./LandingPage";

const navigateMock = vi.fn();

vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual<typeof import("react-router-dom")>("react-router-dom");
  return {
    ...actual,
    useNavigate: () => navigateMock
  };
});

describe("LandingPage", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    navigateMock.mockReset();
    localStorage.clear();
    sessionStorage.clear();
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it("renders title Pictionary", () => {
    render(
      <MemoryRouter>
        <LandingPage />
      </MemoryRouter>
    );
    expect(screen.getByText("Pictionary")).toBeInTheDocument();
  });

  it("renders name input and Play Now button", () => {
    render(
      <MemoryRouter>
        <LandingPage />
      </MemoryRouter>
    );
    expect(screen.getByLabelText("Player name")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: "Play Now" })).toBeInTheDocument();
    expect(screen.getByText("0/20")).toBeInTheDocument();
  });

  it("shows error if name is empty on submit", async () => {
    render(
      <MemoryRouter>
        <LandingPage />
      </MemoryRouter>
    );

    fireEvent.click(screen.getByRole("button", { name: "Play Now" }));

    expect(await screen.findByRole("alert")).toHaveTextContent("Name required");
  });

  it("updates character counter as user types", () => {
    render(
      <MemoryRouter>
        <LandingPage />
      </MemoryRouter>
    );
    fireEvent.change(screen.getByLabelText("Player name"), { target: { value: "Alice" } });
    expect(screen.getByText("5/20")).toBeInTheDocument();
  });

  it("shows name too long error from server", async () => {
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue(new Response("name too long", { status: 400 })));
    render(
      <MemoryRouter>
        <LandingPage />
      </MemoryRouter>
    );
    fireEvent.change(screen.getByLabelText("Player name"), { target: { value: "Alice" } });
    fireEvent.click(screen.getByRole("button", { name: "Play Now" }));
    expect(await screen.findByRole("alert")).toHaveTextContent("Name too long");
  });

  it("shows server unavailable on network error", async () => {
    vi.stubGlobal("fetch", vi.fn().mockRejectedValue(new Error("network")));
    render(
      <MemoryRouter>
        <LandingPage />
      </MemoryRouter>
    );
    fireEvent.change(screen.getByLabelText("Player name"), { target: { value: "Alice" } });
    fireEvent.click(screen.getByRole("button", { name: "Play Now" }));
    expect(await screen.findByRole("alert")).toHaveTextContent("Server unavailable");
  });

  it("disables button during loading", async () => {
    let resolveFetch: (value: Response) => void = () => {};
    const fetchMock = vi.fn(
      () =>
        new Promise<Response>((resolve) => {
          resolveFetch = resolve;
        })
    );
    vi.stubGlobal("fetch", fetchMock);

    render(
      <MemoryRouter>
        <LandingPage />
      </MemoryRouter>
    );

    fireEvent.change(screen.getByLabelText("Player name"), { target: { value: "Alice" } });
    fireEvent.click(screen.getByRole("button", { name: "Play Now" }));

    const button = screen.getByRole("button", { name: "Joining..." });
    expect(button).toBeDisabled();

    resolveFetch(
      new Response(JSON.stringify({ playerId: "p1", roomId: "r1", wsUrl: "ws://localhost:8080/ws?playerId=p1&roomId=r1" }), {
        status: 200,
        headers: { "Content-Type": "application/json" }
      })
    );

    await waitFor(() => {
      expect(navigateMock).toHaveBeenCalledWith("/game");
    });
  });
});
