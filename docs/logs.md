# Documentation des Logs (Journaux d'événements)

Cette page explique en détail les types de logs générés par l'API backend **Arem-Shop** et comment les interpréter. Ces logs sont affichés dans la console lors de l'exécution (ou via `docker compose logs`).

## 1. Moteur de Logs
L'application utilise deux mécanismes principaux :
- Le module standard Go `log` pour les erreurs critiques de démarrage.
- Le middleware `gin.Logger()` pour tracer toutes les requêtes HTTP entrantes.

---

## 2. Logs de Démarrage (Startup Logs)
Lorsqu'Arem-Shop démarre, il s'assure que toutes les dépendances sont prêtes. En cas de problème, il produit un log formatté via `log.Fatalf()`.

**Exemples :**
```text
2026/02/21 14:00:00 load config: error loading .env file
2026/02/21 14:00:01 connect database: failed to connect to host=...
2026/02/21 14:00:02 run server: address already in use
```
- **Structure** : `[Date] [Heure] [Étape ayant échouée]: [Raison technique]`
- **Conséquence** : Ces erreurs sont "fatales", c'est-à-dire que l'application s'arrête immédiatement. Il faut corriger le problème (ex: base de données éteinte, port bloqué) avant de relancer.

---

## 3. Logs de Trafic HTTP (Gin Logger)
Chaque passe via l'API, que ça soit le frontend React ou un test Postman, génère automatiquement une ligne de log détaillée.

**Exemple typique d'une requête réussie :**
```text
[GIN] 2026/02/21 - 14:05:22 | 200 |    2.541ms |       127.0.0.1 | GET      "/health"
[GIN] 2026/02/21 - 14:06:10 | 201 |   14.301ms |    192.168.1.10 | POST     "/auth/login"
[GIN] 2026/02/21 - 14:10:05 | 200 |    5.112ms |       127.0.0.1 | GET      "/products"
```

**Analyse de la structure d'une ligne :**
1. **[GIN]** : Identifiant du routeur.
2. **Date et Heure** : Le moment exact de la requête (ex: `2026/02/21 - 14:06:10`).
3. **Statut HTTP (`200`, `201`, `401`, `500`...)** : Le code de réponse. 
   - `2xx` = Succès.
   - `4xx` = Erreur utilisateur (ex: mauvais mot de passe = `401`, accès refusé = `403`, page introuvable = `404`).
   - `5xx` = Erreur grave côté serveur (ex: base de données plantée).
4. **Temps d'exécution (`14.301ms`)** : Indique si l'API est performante ou lente.
5. **Adresse IP Client (`192.168.1.10`)** : La provenance de l'appareil ayant effectué la requête.
6. **Méthode HTTP (`GET`, `POST`, `PUT`, `DELETE`)** : L'action effectuée.
7. **La Route ou le Chemin (`"/auth/login"`)** : Quelle partie de l'API a été contactée.

---

## 4. Logs de Crashs / Paniques (Gin Recovery)
L'API est protégée par un middleware `gin.Recovery()`. Si un bug grave (ex: pointeur nil en langage Go) survient dans un Handler, le serveur ne s'éteint pas. Il recupère l'erreur, renvoie un statut `500 Internal Server Error` au client, et écrit une **Trace d'Exécution (Stack Trace)** dans les logs.

**Exemple d'un crash intercepté :**
```text
[GIN] 2026/02/21 - 14:30:00 | 500 |      1.1ms |       127.0.0.1 | GET      "/products"
[Recovery] 2026/02/21 - 14:30:00 panic recovered:
runtime error: invalid memory address or nil pointer dereference
... (Suivi du nom des fichiers et numéros de lignes précis causant l'erreur) ...
```

---

## 5. Comment observer les logs en production (Docker) ?

Puisque le backend fonctionne dans l'écosystème Docker via le fichier `docker-compose.yml`, vous pouvez consulter les logs en temps réel avec la commande suivante à la racine du projet :

```bash
docker compose logs -f arem-shop-api
```
- `-f` (suivi en continu, *follow*) : Permet de voir les nouveaux logs s'afficher en direct.
- `arem-shop-api` : Cible précisément le service backend (pour éviter de mélanger avec ceux de PostgreSQL ou du Frontend React).

Pour quitter cette vue, utilisez le raccourci `Ctrl + C`.
