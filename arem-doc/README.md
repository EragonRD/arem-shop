# Arem Shop Documentation Vault

Ce dossier est un mini-vault Obsidian pour comprendre l'API Arem Shop de A a Z.

## Point d'entree

- [[00-Plan-Explication]]

## Parcours recommande

1. [[01-Vue-Produit]]
2. [[02-Architecture-API]]
3. [[03-Modele-Donnees]]
4. [[04-Flux-Metier-Critiques]]
5. [[05-Contrats-API]]
6. [[06-Runbook-Technique]]

## Schema de navigation

![[schemas/plan-explication.png]]

## Generation des schemas PNG HD

```bash
bash arem-doc/scripts/generate-diagrams.sh
```

Le script convertit les fichiers Markdown Mermaid du dossier `arem-doc/mermaid` en PNG HD dans `arem-doc/schemas`.
