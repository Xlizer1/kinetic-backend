# Kinetic Backend

A real-time communication platform backend built with Go, featuring WebSocket support, voice channels with WebRTC signaling, and a RESTful API.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green)
![Swagger](https://img.shields.io/badge/Swagger-Enabled-blue)

## Features

### Authentication
- JWT-based authentication with access and refresh tokens
- User registration and login
- Password reset flow
- Token refresh endpoint

### Servers & Channels
- Create, update, delete servers
- Invite system with unique codes
- Text and voice channels
- Server member management

### Messaging
- Real-time message delivery via WebSocket
- REST API for message history
- Pagination support

### Real-time Features
- WebSocket-based messaging
- Typing indicators
- Online presence tracking
- Voice channel support with WebRTC signaling
- Multi-server deployment with Redis pub/sub

### Voice & Video
- Voice channel rooms
- WebRTC peer-to-peer audio
- SFU integration support (LiveKit)

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP), Gorilla WebSocket
- **Database**: PostgreSQL with GORM
- **Cache/Pub-Sub**: Redis
- **Authentication**: JWT
- **Media**: LiveKit (optional SFU)

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL
- Redis (optional, for multi-server)

### Configuration

Create a `.env` file:

```env
# Server
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=kinetic

# JWT
JWT_SECRET=your-secret-key

# Redis (optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# LiveKit (optional, for scalable voice)
LIVEKIT_API_KEY=
LIVEKIT_API_SECRET=
LIVEKIT_SERVER_URL=
```

### Build & Run

```bash
# Install dependencies
go mod tidy

# Run the server
go run cmd/server/main.go
```

The server will start at `http://localhost:8080`

## API Documentation

### Swagger UI

Access the interactive API docs at: `http://localhost:8080/swagger/index.html`

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | Login user |
| POST | `/api/auth/refresh-token` | Refresh access token |
| POST | `/api/auth/forgot-password` | Request password reset |
| POST | `/api/auth/reset-password` | Reset password |

### User Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users/@me` | Get current user |
| PATCH | `/api/users/@me` | Update current user |
| PATCH | `/api/users/@me/settings` | Update user settings |

### Server Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/servers` | Get user's servers |
| POST | `/api/servers` | Create server |
| POST | `/api/servers/join` | Join server via invite |
| GET | `/api/servers/:id` | Get server |
| PATCH | `/api/servers/:id` | Update server |
| DELETE | `/api/servers/:id` | Delete server |
| POST | `/api/servers/:id/leave` | Leave server |
| GET | `/api/servers/:id/channels` | Get server channels |

### Channel Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/channels` | Create channel |
| GET | `/api/channels/:id` | Get channel |
| PATCH | `/api/channels/:id` | Update channel |
| DELETE | `/api/channels/:id` | Delete channel |

### Message Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/channels/:id/messages` | Get channel messages |
| POST | `/api/channels/:id/messages` | Send message |
| DELETE | `/api/channels/:id/messages/:messageId` | Delete message |

## WebSocket API

Connect to: `ws://localhost:8080/ws`

### Authentication

Send this immediately after connecting:

```json
{
  "type": "AUTHENTICATE",
  "payload": { "token": "your-jwt-token" }
}
```

### Client → Server Events

| Event | Payload |
|-------|---------|
| `JOIN_ROOM` | `{ "channel_id": 1 }` |
| `LEAVE_ROOM` | `{ "channel_id": 1 }` |
| `SEND_MESSAGE` | `{ "channel_id": 1, "content": "Hello" }` |
| `TYPING_START` | `{ "channel_id": 1 }` |
| `TYPING_STOP` | `{ "channel_id": 1 }` |
| `VOICE_JOIN` | `{ "channel_id": 1 }` |
| `VOICE_LEAVE` | `{ "channel_id": 1 }` |
| `VOICE_SIGNAL` | `{ "channel_id": 1, "target_user_id": 2, "sdp": "...", "type": "offer\|answer\|ice" }` |

### Server → Client Events

| Event | Payload |
|-------|---------|
| `NEW_MESSAGE` | Message object |
| `USER_JOINED` | `{ "user_id": 1, "username": "john" }` |
| `USER_LEFT` | `{ "user_id": 1, "username": "john" }` |
| `TYPING` | `{ "channel_id": 1, "user_id": 1, "username": "john" }` |
| `PRESENCE_LIST` | Array of online users |
| `PRESENCE_UPDATE` | `{ "user_id": 1, "status": "online" }` |
| `VOICE_JOIN` | `{ "channel_id": 1, "users": [] }` |
| `VOICE_USER_JOINED` | `{ "channel_id": 1, "user_id": 1, "username": "john" }` |
| `VOICE_OFFER` | WebRTC offer for voice |
| `VOICE_ANSWER` | WebRTC answer for voice |
| `ICE_CANDIDATE` | ICE candidate for NAT traversal |

## Project Structure

```
kinetic-backend/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Auth & logging middleware
│   ├── models/          # Database models
│   ├── repositories/    # Database access layer
│   ├── services/        # Business logic
│   ├── realtime/        # WebSocket & room management
│   └── utils/           # Utilities (JWT, hashing, responses)
├── docs/                # Swagger documentation
└── go.mod               # Dependencies
```

## Development

### Running Tests

```bash
go test ./...
```

### Adding Dependencies

```bash
go get package@version
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Frontend

Looking for the frontend? Check out [Kinetic Frontend](https://github.com/Xlizer1/kinetic-frontend) - a React-based UI that integrates with this backend.

---

Built with ❤️ for real-time communication