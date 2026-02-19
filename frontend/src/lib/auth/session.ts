"use client";

import { STORAGE_KEYS } from "@/lib/constants";
import type { AuthSession, AuthUser, JwtClaims, UserRole } from "@/lib/types";

function safeWindow(): Window | null {
  if (typeof window === "undefined") {
    return null;
  }
  return window;
}

export function saveSession(session: AuthSession): void {
  const win = safeWindow();
  if (!win) {
    return;
  }

  win.localStorage.setItem(STORAGE_KEYS.token, session.token);
  win.localStorage.setItem(STORAGE_KEYS.user, JSON.stringify(session.user));
}

export function clearSession(): void {
  const win = safeWindow();
  if (!win) {
    return;
  }

  win.localStorage.removeItem(STORAGE_KEYS.token);
  win.localStorage.removeItem(STORAGE_KEYS.user);
}

export function getToken(): string {
  const win = safeWindow();
  if (!win) {
    return "";
  }
  return win.localStorage.getItem(STORAGE_KEYS.token) ?? "";
}

export function getUser(): AuthUser | null {
  const win = safeWindow();
  if (!win) {
    return null;
  }

  const raw = win.localStorage.getItem(STORAGE_KEYS.user);
  if (!raw) {
    return null;
  }

  try {
    return JSON.parse(raw) as AuthUser;
  } catch {
    return null;
  }
}

export function getSession(): AuthSession | null {
  const token = getToken();
  const user = getUser();

  if (!token || !user) {
    return null;
  }

  return { token, user };
}

export function getCurrentRole(): UserRole | null {
  return getUser()?.role ?? decodeClaims(getToken())?.role ?? null;
}

export function decodeClaims(token: string): JwtClaims | null {
  if (!token) {
    return null;
  }

  const parts = token.split(".");
  if (parts.length < 2) {
    return null;
  }

  try {
    const payload = parts[1].replace(/-/g, "+").replace(/_/g, "/");
    const json =
      typeof window !== "undefined" && typeof window.atob === "function"
        ? window.atob(payload)
        : Buffer.from(payload, "base64").toString("utf-8");

    return JSON.parse(json) as JwtClaims;
  } catch {
    return null;
  }
}
