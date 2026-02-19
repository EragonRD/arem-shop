import { clearSession, saveSession } from "@/lib/auth/session";
import { isAuthenticated, shouldRedirectToLogin } from "@/lib/auth/guards";

describe("auth guards", () => {
  beforeEach(() => {
    localStorage.clear();
  });

  afterEach(() => {
    clearSession();
  });

  it("redirects to login when no session exists", () => {
    expect(shouldRedirectToLogin("/dashboard")).toBe(true);
    expect(isAuthenticated()).toBe(false);
  });

  it("does not redirect when session exists", () => {
    saveSession({
      token: "mock.token.payload",
      user: {
        id: "u-1",
        name: "Owner",
        email: "owner@shop.com",
        role: "SuperAdmin",
        shopID: "11111111-1111-1111-1111-111111111111",
        createdAt: "2026-02-19T00:00:00Z",
      },
    });

    expect(shouldRedirectToLogin("/products")).toBe(false);
    expect(isAuthenticated()).toBe(true);
  });
});
