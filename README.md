# Polls & Games Application

An interactive polling application built with React.js (frontend) and Go (backend), featuring user authentication, poll creation, voting, and bookmarking capabilities.

## Features

- ✅ User authentication (signup/login) with SQLite database
- ✅ Browse all polls with pagination
- ✅ View individual polls with voting functionality
- ✅ Bookmark polls and view bookmarked polls
- ✅ User dashboard to manage created polls
- ✅ Create new polls with multiple options
- ✅ Toast notifications for success/error messages
- ✅ Responsive design using Tailwind CSS and DaisyUI

## Project Structure

```
polls-and-games/
├── backend/
│   ├── main.go          # Go backend server
│   ├── go.mod           # Go dependencies
│   └── polls.db         # SQLite database (created automatically)
├── frontend/
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── context/     # Auth context
│   │   ├── pages/       # Page components
│   │   ├── App.jsx      # Main app component
│   │   └── main.jsx     # Entry point
│   ├── package.json     # Frontend dependencies
│   ├── vite.config.js   # Vite configuration
│   └── tailwind.config.js
└── README.md
```

## Prerequisites

Before running the application, make sure you have the following installed:

1. **Go** (version 1.21 or higher)
   - Download from: https://golang.org/dl/
   - Verify installation: `go version`

2. **Node.js** (version 18 or higher) and **npm**
   - Download from: https://nodejs.org/
   - Verify installation: `node --version` and `npm --version`

## Installation & Setup

### Step 1: Set up the Backend (Go)

1. Navigate to the backend directory:
   ```bash
   cd polls-and-games/backend
   ```

2. Install Go dependencies:
   ```bash
   go mod download
   ```

   This will install:
   - `github.com/gorilla/mux` - HTTP router
   - `github.com/mattn/go-sqlite3` - SQLite driver
   - `golang.org/x/crypto` - Password hashing

### Step 2: Set up the Frontend (React)

1. Navigate to the frontend directory:
   ```bash
   cd polls-and-games/frontend
   ```

2. Install npm dependencies:
   ```bash
   npm install
   ```

   This will install:
   - React and React DOM
   - React Router DOM
   - Axios (HTTP client)
   - React Toastify (notifications)
   - Tailwind CSS and DaisyUI (styling)
   - Vite (build tool)

## Running the Application

### Option 1: Run Backend and Frontend Separately (Recommended for Development)

#### Terminal 1 - Start the Backend Server:

1. Navigate to backend directory:
   ```bash
   cd polls-and-games/backend
   ```

2. Run the Go server:
   ```bash
   go run main.go
   ```

   The backend server will start on `http://localhost:8080`

   You should see:
   ```
   Server running on http://localhost:8080
   ```

   **Note:** The SQLite database (`polls.db`) will be created automatically in the backend directory on first run.

#### Terminal 2 - Start the Frontend Development Server:

1. Navigate to frontend directory:
   ```bash
   cd polls-and-games/frontend
   ```

2. Start the Vite development server:
   ```bash
   npm run dev
   ```

   The frontend will start on `http://localhost:5173`

   You should see:
   ```
   VITE v5.x.x  ready in xxx ms

   ➜  Local:   http://localhost:5173/
   ➜  Network: use --host to expose
   ```

3. Open your browser and navigate to:
   ```
   http://localhost:5173
   ```

### Option 2: Build and Run Production Version

#### Build Frontend:

1. Navigate to frontend directory:
   ```bash
   cd polls-and-games/frontend
   ```

2. Build the React app:
   ```bash
   npm run build
   ```

   This creates a `dist` folder with production-ready files.

3. Preview the production build:
   ```bash
   npm run preview
   ```

## API Endpoints

The backend provides the following REST API endpoints:

### Authentication
- `POST /api/register` - Register a new user
- `POST /api/login` - Login user

### Polls
- `GET /api/polls?page=1` - Get paginated list of polls
- `GET /api/polls/{id}` - Get a specific poll
- `POST /api/polls` - Create a new poll
- `DELETE /api/polls/{id}` - Delete a poll

### Voting
- `POST /api/vote` - Vote on a poll option

### Bookmarks
- `POST /api/bookmark` - Bookmark a poll
- `DELETE /api/bookmark` - Remove bookmark
- `GET /api/bookmarks/{user_id}` - Get user's bookmarked polls
- `GET /api/check-bookmark?user_id=X&poll_id=Y` - Check if poll is bookmarked

### User
- `GET /api/user/{user_id}/polls` - Get polls created by a user

## Usage Guide

1. **Sign Up**: Click "Sign Up" to create a new account
2. **Login**: Use your credentials to login
3. **Browse Polls**: View all polls on the home page with pagination
4. **View Poll**: Click "View Poll" to see poll details and vote
5. **Vote**: Select an option and click "Vote" (requires login)
6. **Bookmark**: Click the bookmark icon to save polls for later
7. **Create Poll**: Click "Create Poll" to create your own poll
8. **Dashboard**: View and manage your created polls
9. **Bookmarks Page**: View all your bookmarked polls

## Troubleshooting

### Backend Issues

1. **Port 8080 already in use**:
   - Change the port in `backend/main.go` or set `PORT` environment variable:
     ```bash
     PORT=8081 go run main.go
     ```

2. **Database errors**:
   - Delete `polls.db` and restart the server to recreate the database
   - Ensure you have write permissions in the backend directory

3. **Go module errors**:
   ```bash
   cd backend
   go mod tidy
   go mod download
   ```

### Frontend Issues

1. **Port 5173 already in use**:
   - Vite will automatically use the next available port
   - Or change it in `frontend/vite.config.js`

2. **Dependencies not installing**:
   ```bash
   cd frontend
   rm -rf node_modules package-lock.json
   npm install
   ```

3. **CORS errors**:
   - Ensure backend is running on port 8080
   - Check that the proxy configuration in `vite.config.js` is correct

4. **Build errors**:
   ```bash
   cd frontend
   npm run build
   ```
   Check the error messages and ensure all dependencies are installed.

## Database Schema

The SQLite database contains the following tables:

- **users**: User accounts (id, username, email, password, created_at)
- **polls**: Polls (id, title, description, user_id, created_at)
- **options**: Poll options (id, poll_id, text)
- **votes**: User votes (id, user_id, poll_id, option_id, created_at)
- **bookmarks**: User bookmarks (id, user_id, poll_id, created_at)

## Technologies Used

### Backend
- Go 1.21+
- Gorilla Mux (HTTP router)
- SQLite3 (database)
- bcrypt (password hashing)

### Frontend
- React 18
- React Router DOM
- Axios
- React Toastify
- Tailwind CSS
- DaisyUI
- Vite

## Development Notes

- The backend uses CORS to allow requests from the frontend
- Passwords are hashed using bcrypt before storage
- Users can only vote once per poll (votes can be updated)
- Polls are paginated with 10 items per page
- The database file (`polls.db`) is created automatically in the backend directory

## License

This project is open source and available for educational purposes.

