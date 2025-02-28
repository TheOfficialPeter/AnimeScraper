# Anime Scraper

A modern and sleek anime scraping application built with Go and a responsive web interface. This application uses the Jikan API to fetch anime data and provides a beautiful interface for browsing and searching anime.

## Features

- 🎯 Modern and responsive UI using Tailwind CSS
- 🔍 Advanced search functionality with filters
- 📺 Trending anime section
- 🆕 Upcoming anime section
- 🎭 Genre-specific sections (Isekai, Ecchi, Romance)
- 👍 Personalized anime recommendations
- 🏷️ Filter by genre, season, and status
- ⭐ View ratings and details
- 📱 Mobile-friendly design
- 🔖 Bookmark favorite anime titles
- 💾 Local storage for persisting bookmarks

## Prerequisites

- Go 1.21 or higher
- Internet connection (for API calls)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/anime-scraper.git
cd anime-scraper
```

2. Install dependencies:

```bash
go mod download
```

## Running the Application

1. Start the server:

```bash
go run main.go
```

2. Open your browser and navigate to:

```
http://localhost:8080
```

## Usage

- **Search**: Use the search bar to find specific anime titles
- **Filters**: Apply filters for genre, season, and status
- **Trending**: Click the "Trending" button to see popular anime
- **Upcoming**: Click the "Upcoming" button to see upcoming releases
- **Genre Sections**: Browse specific genres like Isekai, Ecchi, and Romance
- **Recommendations**: Find similar anime based on your selected title
  - Search for an anime (minimum 3 characters)
  - Select from the suggestions dropdown
  - View personalized recommendations with user comments
- **Bookmarks**: Save and manage your favorite anime
  - Click the bookmark icon on any anime card to save it
  - Access your bookmarked anime from the Bookmarks tab
  - Bookmarks are saved locally in your browser
  - Toggle bookmarks on/off by clicking the bookmark icon again

## API Endpoints

- `GET /api/trending` - Get trending anime
- `GET /api/upcoming` - Get upcoming anime
- `GET /api/search` - Search anime with filters
  - Query parameters:
    - `q`: Search query
    - `genre`: Genre ID
    - `season`: Season (spring, summer, fall, winter)
    - `status`: Anime status (airing, complete, upcoming)
- `GET /api/isekai` - Get popular Isekai anime
- `GET /api/ecchi` - Get popular Ecchi anime
- `GET /api/romance` - Get popular Romance anime
- `GET /api/suggestions` - Get anime search suggestions
  - Query parameters:
    - `q`: Search query (minimum 3 characters)
- `GET /api/recommendations/:id` - Get recommendations for a specific anime
  - Parameters:
    - `id`: MAL ID of the anime

## Contributing

Feel free to submit issues and enhancement requests!
