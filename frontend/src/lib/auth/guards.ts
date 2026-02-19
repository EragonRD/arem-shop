"use client";

import { clearSession, getSession } from "@/lib/auth/session";

export function isAuthenticated(): boolean {
  return getSession() !== null;
}

export function shouldRedirectToLogin(pathname: string): boolean {
  if (pathname === "/login") {
    return false;
  }
  return !isAuthenticated();
}

export function handleUnauthorized(): void {
  clearSession();
}
