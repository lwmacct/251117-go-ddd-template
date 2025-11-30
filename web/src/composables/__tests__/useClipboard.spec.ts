import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { useClipboard } from "../useClipboard";

describe("useClipboard", () => {
  const originalClipboard = navigator.clipboard;
  const originalIsSecureContext = window.isSecureContext;

  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    // Restore original values
    Object.defineProperty(navigator, "clipboard", {
      value: originalClipboard,
      configurable: true,
    });
  });

  it("should initialize with default values", () => {
    const { copied, error } = useClipboard();

    expect(copied.value).toBe(false);
    expect(error.value).toBeNull();
  });

  it("should copy text successfully using Clipboard API", async () => {
    // Mock Clipboard API properly
    const writeTextMock = vi.fn().mockResolvedValue(undefined);
    Object.defineProperty(navigator, "clipboard", {
      value: { writeText: writeTextMock },
      configurable: true,
    });
    Object.defineProperty(window, "isSecureContext", {
      value: true,
      configurable: true,
    });

    const { copy, copied } = useClipboard();
    const result = await copy("test text");

    expect(result).toBe(true);
    expect(copied.value).toBe(true);
    expect(writeTextMock).toHaveBeenCalledWith("test text");
  });

  it("should reset copied state after successDuration", async () => {
    vi.useFakeTimers();

    const writeTextMock = vi.fn().mockResolvedValue(undefined);
    Object.defineProperty(navigator, "clipboard", {
      value: { writeText: writeTextMock },
      configurable: true,
    });
    Object.defineProperty(window, "isSecureContext", {
      value: true,
      configurable: true,
    });

    const { copy, copied } = useClipboard({ successDuration: 1000 });

    await copy("test");
    expect(copied.value).toBe(true);

    vi.advanceTimersByTime(1000);
    expect(copied.value).toBe(false);

    vi.useRealTimers();
  });

  it("should handle copy failure", async () => {
    // Mock Clipboard API failure
    const writeTextMock = vi.fn().mockRejectedValue(new Error("Copy failed"));
    Object.defineProperty(navigator, "clipboard", {
      value: { writeText: writeTextMock },
      configurable: true,
    });
    Object.defineProperty(window, "isSecureContext", {
      value: true,
      configurable: true,
    });

    const { copy, copied, error } = useClipboard();
    const result = await copy("test");

    expect(result).toBe(false);
    expect(copied.value).toBe(false);
    expect(error.value).toBe("Copy failed");
  });
});
