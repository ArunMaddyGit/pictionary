import { renderHook, act } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { useThrottle } from "./useThrottle";

describe("useThrottle", () => {
  it("throttles updates by delay", () => {
    vi.useFakeTimers();
    const { result, rerender } = renderHook(
      ({ value, delay }) => useThrottle(value, delay),
      { initialProps: { value: "a", delay: 100 } }
    );

    expect(result.current).toBe("a");

    rerender({ value: "b", delay: 100 });
    expect(result.current).toBe("a");

    act(() => {
      vi.advanceTimersByTime(99);
    });
    expect(result.current).toBe("a");

    act(() => {
      vi.advanceTimersByTime(1);
    });
    expect(result.current).toBe("b");

    vi.useRealTimers();
  });
});
