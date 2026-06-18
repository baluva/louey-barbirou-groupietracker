# 📈 Crypto Tracker (Go)

Application web en **Go** (`net/http`) pour suivre les cryptomonnaies : recherche,
affichage des prix / capitalisation / variation 24h, et gestion d'une liste de favoris.

> Projet Ynov. Le dépôt s'appelle historiquement *groupietracker*, mais le code
> implémente un **tracker de cryptomonnaies** (modèle de données `CryptoData`,
> prix, market cap, `favoriteCoins`).

## ✨ Fonctionnalités

- Page d'accueil et recherche de cryptomonnaies
- Affichage prix / capitalisation / variation sur 24h
- Filtres de recherche
- Gestion de favoris (persistés dans `favorite.json`)

## 🏗️ Structure

```
.
├── main.go                 # Point d'entrée (init templates + serveur)
├── routeur/routeur.go      # Routes HTTP
├── controller/handlers.go  # Contrôleurs (pages, filtres)
├── backend/
│   ├── functions.go        # Appels API + logique
│   └── struct.go           # Modèles (CryptoData, Favorite…)
├── templates/              # Pages HTML (accueil, search, favorite, error)
├── assets/                 # CSS et images
└── favorite.json           # Favoris persistés
```

## 🚀 Lancer le projet

```bash
go run .
```
Puis ouvrir **http://localhost:8000/accueil**.

## 🔑 Clé API

L'application consomme l'API **CoinMarketCap** (`pro-api.coinmarketcap.com`), qui
nécessite une clé. La clé est lue depuis la variable d'environnement
`CMC_PRO_API_KEY` (voir `.env.example`) :

```bash
# PowerShell
$env:CMC_PRO_API_KEY="votre_cle"; go run .
# Linux/macOS
CMC_PRO_API_KEY="votre_cle" go run .
```

> ⚠️ Une ancienne clé était auparavant codée en dur dans le code. Elle a été
> retirée, mais comme elle reste visible dans l'historique Git, **régénérez-la**
> sur le dashboard CoinMarketCap.

## 🛣️ Routes

| Route | Description |
|-------|-------------|
| `/accueil` | Page d'accueil |
| `/search` | Recherche |
| `/favorite` | Favoris |
| `/filter_submit` | Application des filtres |
| `/static/` | Fichiers statiques (CSS, images) |

## 👨‍💻 Auteur

Louey Barbirou — Ynov.
