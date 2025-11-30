import { describe, it, expect } from "vitest";
import { ref } from "vue";
import { useDebouncedRef } from "../useDebounce";

describe("useDebouncedRef", () => {
  it("should return initial value immediately", () => {
    const source = ref("initial");
    const debouncedValue = useDebouncedRef(source);

    expect(debouncedValue.value).toBe("initial");
  });

  it("should accept options parameter", () => {
    const source = ref("test");
    const debouncedValue = useDebouncedRef(source, { delay: 500, immediate: true });

    expect(debouncedValue.value).toBe("test");
  });

  it("should work with different types", () => {
    const numberSource = ref(42);
    const debouncedNumber = useDebouncedRef(numberSource);
    expect(debouncedNumber.value).toBe(42);

    const objectSource = ref({ name: "test" });
    const debouncedObject = useDebouncedRef(objectSource);
    expect(debouncedObject.value).toEqual({ name: "test" });
  });
});
