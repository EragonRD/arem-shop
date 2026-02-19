import React from "react";
import userEvent from "@testing-library/user-event";
import { render, screen } from "@testing-library/react";

import { I18nProvider } from "@/lib/i18n/I18nProvider";
import { useI18n } from "@/lib/i18n/useI18n";

function Harness() {
  const { locale, toggleLocale, t } = useI18n();

  return (
    <div>
      <p data-testid="locale">{locale}</p>
      <p>{t("loginTitle")}</p>
      <button type="button" onClick={toggleLocale}>
        toggle
      </button>
    </div>
  );
}

describe("I18nProvider", () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it("switches locale between fr and en", async () => {
    const user = userEvent.setup();

    render(
      <I18nProvider>
        <Harness />
      </I18nProvider>,
    );

    expect(screen.getByText("Connexion")).toBeInTheDocument();
    expect(screen.getByTestId("locale")).toHaveTextContent("fr");

    await user.click(screen.getByRole("button", { name: "toggle" }));

    expect(screen.getByText("Login")).toBeInTheDocument();
    expect(screen.getByTestId("locale")).toHaveTextContent("en");
  });
});
