"use client";

import { useState } from "react";

import { getSession } from "@/lib/auth/session";
import { useI18n } from "@/lib/i18n/useI18n";
import { dataClient } from "@/lib/services/dataClient";

export default function ProfilePage() {
    const { t } = useI18n();
    const session = getSession();
    const [error, setError] = useState("");
    const [success, setSuccess] = useState(false);

    if (!session || session.user.role !== "SuperAdmin") {
        return (
            <section className="rounded-2xl border border-brand-border/70 bg-white/90 p-6 text-brand-fg shadow-soft">
                <h1 className="font-display text-2xl font-semibold">{t("navProfile")}</h1>
                <p className="mt-3 text-sm text-brand-fg/75">{t("forbidden")}</p>
            </section>
        );
    }

    return (
        <section className="space-y-5">
            <header>
                <h1 className="font-display text-3xl font-semibold text-brand-fg">{t("navProfile")}</h1>
                <p className="text-sm text-brand-fg/70">Gérez les informations de votre boutique.</p>
            </header>

            {error ? <p className="rounded-xl bg-red-50 px-3 py-2 text-sm text-red-700">{error}</p> : null}
            {success ? <p className="rounded-xl bg-green-50 px-3 py-2 text-sm text-green-700">Modifications enregistrées avec succès !</p> : null}

            <div className="rounded-2xl border border-brand-border/70 bg-white/90 p-6 shadow-soft">
                <h2 className="font-display text-xl font-semibold text-brand-fg">Paramètres de la Boutique</h2>
                <p className="mb-4 text-sm text-brand-fg/70">Modifiez le nom et le numéro WhatsApp de votre boutique.</p>
                <form
                    onSubmit={async (e) => {
                        e.preventDefault();
                        setError("");
                        setSuccess(false);
                        const formData = new FormData(e.currentTarget);
                        const name = formData.get("shopName") as string;
                        const whatsAppNumber = formData.get("whatsAppNumber") as string;

                        try {
                            await dataClient.updateShop(session.token, { name, whatsAppNumber });

                            const newSession = {
                                ...session,
                                user: {
                                    ...session.user,
                                    shopName: name,
                                    whatsAppNumber: whatsAppNumber,
                                }
                            };
                            localStorage.setItem("arem_user", JSON.stringify(newSession));
                            setSuccess(true);
                            // Refresh after a short delay so user sees the success message
                            setTimeout(() => window.location.reload(), 800);
                        } catch (err) {
                            setError(err instanceof Error ? err.message : "Erreur lors de la mise à jour");
                        }
                    }}
                    className="flex flex-col max-w-sm gap-4"
                >
                    <div>
                        <label className="mb-1 block text-sm font-medium text-brand-fg">Nom de la boutique</label>
                        <input
                            type="text"
                            name="shopName"
                            defaultValue={session.user.shopName || ""}
                            placeholder="Ex: Ma Boutique"
                            className="w-full rounded-xl border border-brand-border bg-brand-bg px-3 py-2 text-sm text-brand-fg focus:border-brand-primary focus:outline-none focus:ring-1 focus:ring-brand-primary"
                        />
                    </div>
                    <div>
                        <label className="mb-1 block text-sm font-medium text-brand-fg">Numéro WhatsApp</label>
                        <input
                            type="text"
                            name="whatsAppNumber"
                            defaultValue={session.user.whatsAppNumber || ""}
                            placeholder="Ex: +33612345678"
                            className="w-full rounded-xl border border-brand-border bg-brand-bg px-3 py-2 text-sm text-brand-fg focus:border-brand-primary focus:outline-none focus:ring-1 focus:ring-brand-primary"
                        />
                    </div>
                    <button
                        type="submit"
                        className="rounded-xl bg-brand-primary px-4 py-2 text-sm font-medium text-white transition hover:bg-brand-primary/90 mt-2"
                    >
                        Enregistrer
                    </button>
                </form>
            </div>

            <div className="rounded-2xl border border-brand-border/70 bg-white/90 p-6 shadow-soft">
                <h2 className="font-display text-xl font-semibold text-brand-fg">Informations du compte</h2>
                <div className="mt-3 space-y-2 text-sm text-brand-fg/80">
                    <p><span className="font-medium text-brand-fg">Nom :</span> {session.user.name}</p>
                    <p><span className="font-medium text-brand-fg">Email :</span> {session.user.email}</p>
                    <p><span className="font-medium text-brand-fg">Rôle :</span> {session.user.role}</p>
                    <p><span className="font-medium text-brand-fg">Shop ID :</span> <code className="text-xs bg-brand-secondary/30 px-1.5 py-0.5 rounded">{session.user.shopID}</code></p>
                </div>
            </div>
        </section>
    );
}
