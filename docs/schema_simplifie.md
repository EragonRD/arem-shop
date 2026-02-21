```mermaid
flowchart TD
    %% Styling pour l'affichage non technique
    classDef shop fill:#4a148c,stroke:#2a0066,stroke-width:2px,color:#fff,font-weight:bold,font-size:16px,padding:15px;
    classDef section fill:#e8eaf6,stroke:#9fa8da,stroke-width:2px,rx:10px,ry:10px,color:#1a237e,font-size:14px,padding:10px;
    classDef relation fill:none,stroke:none,color:#4a148c,font-style:italic;

    %% Nœud central
    Shop["🛒 MA BOUTIQUE\n(Nom, Numéro WhatsApp, Sous-numéro)"]:::shop

    %% Éléments appartenant à la boutique
    Employes["👥 EMPLOYES\n(Gérants, Caissiers)"]:::section
    Produits["🏷️ PRODUITS\n(Nom, Prix, Stock dispo)"]:::section
    Ventes["🧾 VENTES & DEPENSES\n(Historique financier)"]:::section
    Categories["📂 CATEGORIES\n(Familles d'articles)"]:::section

    %% Relations simples et descriptives
    Shop -->|"Emploie des"| Employes
    Shop -->|"Vend des"| Produits
    Shop -->|"Enregistre des"| Ventes
    Shop -->|"Organise en"| Categories

    %% Interactions internes optionnelles
    Produits -.->|"Appartient à une"| Categories
    Ventes -.->|"Concerne un"| Produits
```

### Explication Rapide :
* **Ma Boutique** est le cœur du système. Tout le reste lui appartient.
* Les **Employés** (comme vous ou vos vendeurs) se connectent pour gérer la boutique.
* Les **Produits** sont regroupés par **Catégories** (ex: Téléphones, Ordinateurs).
* Chaque fois qu'un produit est vendu, une ligne s'ajoute dans les **Ventes**, ce qui baisse automatiquement le stock du produit concerné.
