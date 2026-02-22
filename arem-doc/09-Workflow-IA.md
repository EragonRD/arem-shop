# Workflow d'Intégration des Intelligences Artificielles

**Membres du groupe :**
* Alex-Gabriel CHITU
* Mattew VENET
* Enzo Toni REALE
* Rabah Amir DEBBAH

Ce document détaille la méthodologie hybride adoptée par l'équipe pour la conception, le développement, la documentation et le déploiement du projet, en utilisant une synergie de plusieurs Intelligences Artificielles (IA).

## 1. Schéma Global du Workflow

Le développement s'est déroulé en 5 phases séquentielles. Chaque phase a été pilotée par l'IA la plus performante dans son domaine spécifique.

*(Voir l'image HD : `arem-doc/images/workflow-hd.png`)*
![Schéma du Workflow](images/workflow-hd.png)

```mermaid
flowchart TD
    %% Styles
    classDef etape fill:#e3f2fd,stroke:#1565c0,stroke-width:2px,color:#0d47a1,font-weight:bold
    classDef ia fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c
    classDef val fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20,font-style:italic

    %% Noeuds
    A[Étape 1 : Conception & Prompt Engineering]:::etape
    A_IA(Claude Opus / Sonnet 4.6):::ia

    B[Étape 2 : Développement Backend & Tests]:::etape
    B_IA(Codex 5.3 intégré à VS Code):::ia
    B_Val[Validation croisée OS : Windows / Linux]:::val

    C[Étape 3 : Développement Frontend]:::etape
    C_IA(Codex 5.3):::ia

    D[Étape 4 : Documentation & Architecture globale]:::etape
    D_IA(Gemini 3.0 Pro):::ia

    E[Étape 5 : Hotfixes complexes & Dockerisation]:::etape
    E_IA(Gemini 3.1 Pro & Claude Opus 4.6):::ia
    E_Val[Image Docker optimisée et validée]:::val

    %% Flow et relations
    A -->|Génère le prompt strict| B
    A -.->|Piloté par| A_IA
    B --> B_Val
    B -.->|Piloté par| B_IA
    B_Val --> C
    C --> D
    C -.->|Piloté par| C_IA
    D --> E
    D -.->|Piloté par| D_IA
    E --> E_Val
    E -.->|Débuggé par| E_IA
```

## 2. Répartition des Rôles IA et Temps Gagné

Aucune IA n'étant parfaite pour toutes les tâches, nous avons utilisé un écosystème d'IA afin de pallier les faiblesses individuelles de chacune (exemple : manque de contexte de Codex, compensé par Gemini). L'automatisation des tâches rébarbatives (Boilerplate, écriture des tests, formatage textuel) a permis au groupe de se concentrer sur l'architecture, la liaison Front/Back, l'orchestration des données et les tests croisés.

### Résumé Tabulaire classique

| Intelligence Artificielle | Rôle et périmètre d'action | Gain de temps estimé |
| :--- | :--- | :--- |
| **Codex 5.3 (VS Code)** | Création des interfaces (Front), logique (Back) et tests unitaires | **~ 25 heures** *(~15h Back, ~10h Front)* |
| **Gemini 3.0 Pro** | Documentation Markdown, intégration des schémas Mermaid pour Obsidian | **~ 6 à 8 heures** |
| **Claude 4.6 & Gemini 3.1 Pro** | Débogage lourd multi-OS, hotfixes complexes et optimisation Docker | **~ 3 à 5 heures** |
| **Total Estimé** | Gain sur l'ensemble du cycle de vie du projet | **Entre 34 et 38 heures** |

---

### Résumé sous forme de Schéma Diagrammatique 

*(Voir l'image HD : `arem-doc/images/repartition-hd.png`)*
![Schéma des temps](images/repartition-hd.png)

```mermaid
flowchart LR
    classDef ia fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c,font-weight:bold
    classDef role fill:#e3f2fd,stroke:#1565c0,stroke-width:2px,color:#0d47a1
    classDef time fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px,color:#1b5e20,font-weight:bold

    subgraph "Écosystème IA : Rôles & Gains de Productivité"
        direction LR
        Codex[Codex 5.3]:::ia --- R_Codex[Backend & Base de Données<br/>Tests Unitaires<br/>Frontend]:::role
        R_Codex --- T_Codex(Gain de temps : ~25 heures):::time
        
        G3[Gemini 3.0 Pro]:::ia --- R_G3[Documentation Markdown<br/>Schémas Mermaid]:::role
        R_G3 --- T_G3(Gain de temps : ~6 à 8 heures):::time
        
        G31[Gemini 3.1 Pro<br/>& Claude 4.6]:::ia --- R_G31[Hotfixes Complexes<br/>Optimisation Docker]:::role
        R_G31 --- T_G31(Gain de temps : ~3 à 5 heures):::time
    end
```

## 3. Détail des Étapes de Développement

1. **Prompt Engineering (Init) :** Utilisation de Claude 4.6 pour rédiger des consignes de code d'une rigueur absolue en préparation du terrain.
2. **Setup & Backend :** Codex 5.3 prend le relais dans VS Code pour générer le code source étape par étape. Validation manuelle du backend via l'exécution des tests générés, sur des machines **Windows** et **Linux**.
3. **Frontend :** Génération des interfaces et connexion aux routes de l'API par Codex 5.3. Étant techniquement simple, le but a surtout été de relier le Back et le Front.
4. **Documentation (Obsidian) :** L'intégralité du code source validé a été fourni à Gemini 3.0 Pro. Ce denier, grâce à sa très large fenêtre de contexte, a processé le projet complet pour générer la documentation enrichie expliquant l'architecture.
5. **Optimisation finale & Conteneurisation :** Des conflits systèmes (ex: Podman vs Docker sur certains ordinateurs) ont entraîné des bugs. Gemini 3.1 et Claude 4.6 ont été utilisés de concert pour rivaliser d'ingéniosité afin d'alléger les conteneurs (Dockerfile) et résoudre les blocages multi-plateformes.

