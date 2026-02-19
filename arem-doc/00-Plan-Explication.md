# Plan d'explication du projet

Objectif: expliquer l'API multi-tenant Arem Shop de facon progressive, pour qu'une equipe puisse comprendre le fonctionnel, l'architecture, les invariants de securite, puis contribuer rapidement.

## Fil conducteur

![[schemas/plan-explication.png]]

## Etapes de presentation

1. **Vision produit** ([[01-Vue-Produit]])  
   Probleme resolu, acteurs, cas d'usage principaux.
2. **Architecture API** ([[02-Architecture-API]])  
   Router, middlewares, handlers, services, repositories, DB.
3. **Modele de donnees** ([[03-Modele-Donnees]])  
   Tables, relations, contraintes SQL, invariants metier.
4. **Flux critiques** ([[04-Flux-Metier-Critiques]])  
   Login JWT, transaction de vente, catalogue public.
5. **Contrats API** ([[05-Contrats-API]])  
   Endpoints, roles, payloads, erreurs standards.
6. **Runbook technique** ([[06-Runbook-Technique]])  
   Setup, exploitation, points de surveillance.
7. **Docker & Déploiement** ([[07-Docker-Deployment]])  
   Dockerfile multi-stage, docker-compose, cycle de démarrage.
8. **Guide de tests Docker** ([[08-Guide-Tests-Docker]])  
   Script de test automatisé, flux de test, résultats attendus.

## Resultat attendu

- Chaque intervenant comprend **ou** modifier le code et **comment** le faire sans casser l'isolation multi-tenant.
- Les schemas PNG sont exploitables dans des supports externes (slides, docs PDF, wiki).
