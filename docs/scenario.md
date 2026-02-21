# Scénario et Arborescence de l'Application

Ce fichier décrit la structure des pages du frontend (Next.js) et le rôle de chacune selon le parcours de l'utilisateur.

```mermaid
flowchart TD
    %% Couleurs pour différencier les grandes zones
    classDef public fill:#e0f7fa,stroke:#006064,stroke-width:2px,color:#004d40
    classDef auth fill:#fff3e0,stroke:#e65100,stroke-width:2px,color:#bf360c
    classDef dashboard fill:#f3e5f5,stroke:#4a148c,stroke-width:2px,color:#311b92,font-weight:bold
    classDef crud fill:#e8eaf6,stroke:#1a237e,stroke-width:1px,color:#1a237e
    classDef settings fill:#e8f5e9,stroke:#1b5e20,stroke-width:2px,color:#1b5e20

    %% Pages Publiques
    Root["/"]:::public
    PublicCatalog["/public/[shopID]/page.tsx\n(Catalogue Vendeur Public)"]:::public
    
    %% Authentification
    Login["/login/page.tsx\n(Connexion des Employés)"]:::auth

    %% Zone Protégée (Dashboard & Gestion)
    Dashboard["/dashboard/page.tsx\n(Résumé Financier & Alertes)"]:::dashboard
    
    ProductsList["/products/page.tsx\n(Liste de tous les Produits)"]:::crud
    ProductCreate["/products/new/page.tsx\n(Ajouter un nouveau Produit)"]:::crud
    ProductEdit["/products/[id]/edit/page.tsx\n(Modifier un Produit Existant)"]:::crud
    
    TransactionCreate["/transactions/new/page.tsx\n(Nouvelle Vente / Dépense)"]:::crud

    Profile["/profile/page.tsx\n(Paramètres Boutique & Compte)"]:::settings

    %% Les chemins
    Root -->|"Redirige vers"| Login
    Login -->|"Si succès"| Dashboard

    Dashboard --> ProductsList
    Dashboard --> TransactionCreate
    Dashboard --> Profile

    ProductsList -->|"Bouton Créer"| ProductCreate
    ProductsList -->|"Bouton Éditer"| ProductEdit

    %% Catalogues indépendants
    Root -.->|"Lien partagé sur WhatsApp"| PublicCatalog
```

## Détail des pages :

1. **`/` (Racine)** : Un simple point d'entrée qui redirige automatiquement l'utilisateur vers son tableau de bord s'il est connecté, ou vers `/login` sinon.
2. **`/login`** : La page permettant la connexion des *Admins* ou *SuperAdmins* à leur boutique respective.
3. **`/dashboard`** : Le cœur de l'application sécurisée. Il affiche les métriques clés (chiffre d'affaires, dépenses, bénéfices) et alerte sur les ruptures de stock.
4. **`/products`** : Un tableau affichant l'intégralité de l'inventaire du magasin (prix, stock, catégorie).
5. **`/products/new`** & **`/products/[id]/edit`** : Les formulaires (avec upload d'image local) pour ajouter ou mettre à jour un produit de l'inventaire.
6. **`/transactions/new`** : Le point de vente. C'est ici que l'employé déclare une *Vente* (qui réduit le stock d'un produit) ou une *Dépense / Retrait* d'argent de la caisse.
7. **`/profile`** : La page de profil SuperAdmin. Permet de renommer la boutique et de modifier le numéro WhatsApp utilisé pour les liens de contact sur la vitrine publique.
8. **`/public/[shopID]`** : La vitrine orientée client. Elle est générée automatiquement et permet aux visiteurs de voir le catalogue des produits et de contacter le vendeur via un bouton WhatsApp, sans avoir besoin d'un compte.
