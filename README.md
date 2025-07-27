# 🔔 Notification Service

A production-ready notification service built in Go that supports multiple notification channels including email, Slack, and in-app notifications. Features include notification scheduling, template management, and a robust API.

## 🏗️ Architecture Overview

### System Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Gateway   │    │  Notification   │    │   PostgreSQL    │
│   (Gin Router)  │◄──►│    Service      │◄──►│    Database     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │    │   Scheduler     │    │   Data Volume   │
│   (REST API)    │    │  (Background)   │    │   (Persistent)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Service Architecture
```
notification-service/
├── main.go                          # Application entry point
├── internal/
│   ├── config/                      # Configuration management
│   │   └── config.go               # Environment variables & defaults
│   ├── database/                    # Database layer
│   │   └── database.go             # PostgreSQL connection & migrations
│   ├── models/                      # Data models
│   │   └── models.go               # GORM models & JSON handling
│   ├── services/                    # Business logic layer
│   │   ├── notification_service.go # Core notification orchestration
│   │   ├── email_sender.go         # SMTP email delivery
│   │   ├── slack_sender.go         # Slack API integration
│   │   └── in_app_sender.go        # In-app notification handler
│   ├── handlers/                    # HTTP request handlers
│   │   └── handlers.go             # REST API endpoints
│   └── scheduler/                   # Background job processing
│       └── scheduler.go            # Scheduled notification processor
├── docker-compose.yml              # Multi-service orchestration
├── Dockerfile                      # Multi-stage container build
└── .env                           # Environment configuration
```

## 🛠️ Tech Stack

### Backend
- **Language:** Go 1.21
- **Framework:** Gin (HTTP router)
- **ORM:** GORM v2
- **Database:** PostgreSQL 15
- **Scheduler:** GoCron (background jobs)
- **Configuration:** godotenv

### External Services
- **Email:** SMTP (Gmail, Outlook, SendGrid)
- **Slack:** Slack Web API
- **Templates:** Go Templates

### Infrastructure
- **Containerization:** Docker & Docker Compose
- **Database:** PostgreSQL with persistent volumes
- **Health Checks:** Built-in health monitoring
- **Logging:** Structured logging with GORM

### Development Tools
- **Testing:** Go testing framework
- **Dependency Management:** Go modules
- **Code Quality:** Go linting standards

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)
- Git

### 1. Clone the Repository
```bash
git clone <repository-url>
cd notification-service
```

### 2. Configure Environment
```bash
# Copy environment template
cp env.example .env

# Edit .env with your credentials
nano .env
```

### 3. Start Services
```bash
# Start all services
docker-compose up -d

# Check service status
docker-compose ps
```

### 4. Verify Installation
```bash
# Health check
curl http://localhost:8080/health

# Expected response:
# {"service":"notification-service","status":"ok"}
```

## ⚙️ Configuration

### Environment Variables

Create a `.env` file with the following configuration:

```bash
# Database Configuration
DATABASE_URL=postgres://notification_user:notification_password@localhost:5432/notifications?sslmode=disable

# Email Configuration
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-app-password

# Slack Configuration
SLACK_TOKEN=xoxb-your-slack-bot-token
SLACK_CHANNEL=#general

# JWT Configuration
JWT_SECRET=your-jwt-secret-key

# Environment
ENVIRONMENT=development

# Server Configuration
PORT=8080
```

### Email Provider Setup

#### Gmail
1. Enable 2-Step Verification
2. Generate App Password: https://myaccount.google.com/apppasswords
3. Use App Password instead of regular password

#### Outlook/Hotmail
```bash
EMAIL_HOST=smtp-mail.outlook.com
EMAIL_PORT=587
EMAIL_USERNAME=your-email@outlook.com
EMAIL_PASSWORD=your-password
```

#### SendGrid
```bash
EMAIL_HOST=smtp.sendgrid.net
EMAIL_PORT=587
EMAIL_USERNAME=apikey
EMAIL_PASSWORD=your-sendgrid-api-key
```

### Slack Setup
1. Create Slack App: https://api.slack.com/apps
2. Add Bot Token Scopes: `chat:write`, `chat:write.public`
3. Install app to workspace
4. Copy Bot User OAuth Token

## 📡 API Reference

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication
Currently, the API doesn't require authentication for simplicity. In production, implement JWT-based authentication.

### Endpoints

#### Health Check
```http
GET /health
```

#### Notifications

**Send Notification**
```http
POST /api/v1/notifications
Content-Type: application/json

{
  "type": "email|slack|in_app",
  "title": "Notification Title",
  "message": "Notification message",
  "recipient": "user@example.com|#channel|user123",
  "template_id": 1,
  "metadata": {
    "key": "value"
  }
}
```

**Schedule Notification**
```http
POST /api/v1/notifications/schedule
Content-Type: application/json

{
  "type": "email",
  "title": "Scheduled Notification",
  "message": "This will be sent later",
  "recipient": "user@example.com",
  "scheduled_at": "2024-01-15T10:30:00Z"
}
```

**Get Notifications**
```http
GET /api/v1/notifications?page=1&limit=10
```

**Get Notification by ID**
```http
GET /api/v1/notifications/{id}
```

**Update Notification**
```http
PUT /api/v1/notifications/{id}
Content-Type: application/json

{
  "title": "Updated Title",
  "message": "Updated message"
}
```

**Delete Notification**
```http
DELETE /api/v1/notifications/{id}
```

#### Templates

**Create Template**
```http
POST /api/v1/templates
Content-Type: application/json

{
  "name": "welcome_email",
  "type": "email",
  "subject": "Welcome!",
  "content": "Hello {{.Name}}, welcome to our platform!",
  "variables": {
    "Name": "string"
  },
  "is_active": true
}
```

**Get Templates**
```http
GET /api/v1/templates
```

**Get Template by ID**
```http
GET /api/v1/templates/{id}
```

**Update Template**
```http
PUT /api/v1/templates/{id}
```

**Delete Template**
```http
DELETE /api/v1/templates/{id}
```

#### Channels

**Get Channels**
```http
GET /api/v1/channels
```

**Test Channel**
```http
POST /api/v1/channels/test
Content-Type: application/json

{
  "type": "email",
  "recipient": "test@example.com"
}
```

## 🧪 Testing

### Manual Testing

#### Test Email Notification
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "title": "Test Email",
    "message": "This is a test email!",
    "recipient": "your-email@example.com"
  }'
```

#### Test Slack Notification
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "slack",
    "title": "Test Slack",
    "message": "This is a test Slack message!",
    "recipient": "#general"
  }'
```

#### Test In-App Notification
```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in_app",
    "title": "Test In-App",
    "message": "This is an in-app notification!",
    "recipient": "user123"
  }'
```

### Automated Testing
```bash
# Run unit tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## 🐳 Docker Deployment

### Production Deployment
```bash
# Build and start services
docker-compose up --build -d

# View logs
docker-compose logs -f notification-service

# Scale services
docker-compose up -d --scale notification-service=3
```

### Development with Hot Reload
```bash
# Start services in development mode
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

## 📊 Monitoring & Logs

### Health Monitoring
- **Health Endpoint:** `GET /health`
- **Database Health:** Automatic connection monitoring
- **Service Status:** Docker health checks

### Logging
```bash
# View service logs
docker-compose logs notification-service

# Follow logs in real-time
docker-compose logs -f notification-service

# View database logs
docker-compose logs postgres
```

### Metrics
- Request/response times
- Database query performance
- Notification delivery status
- Error rates and types

## 🔧 Development

### Local Development Setup
```bash
# Install dependencies
go mod download

# Run locally
go run main.go

# Run with hot reload (requires air)
air
```

### Database Migrations
```bash
# Auto-migration (handled automatically)
# Tables are created on service startup
```

### Code Structure
```
internal/
├── config/          # Configuration management
├── database/        # Database operations
├── models/          # Data structures
├── services/        # Business logic
├── handlers/        # HTTP handlers
└── scheduler/       # Background jobs
```

## 🚨 Troubleshooting

### Common Issues

#### Email Authentication Failed
- **Issue:** "Username and Password not accepted"
- **Solution:** Use App Password instead of regular password for Gmail

#### Slack Token Invalid
- **Issue:** "invalid_auth" error
- **Solution:** Regenerate Slack bot token and update .env

#### Database Connection Failed
- **Issue:** Cannot connect to PostgreSQL
- **Solution:** Check if PostgreSQL container is running: `docker-compose ps`

#### Port Already in Use
- **Issue:** Port 8080 or 5432 already allocated
- **Solution:** Change ports in docker-compose.yml or stop conflicting services

### Debug Mode
```bash
# Enable debug logging
export GIN_MODE=debug
docker-compose up
```

## 🔒 Security Considerations

### Production Security
- Use strong JWT secrets
- Enable HTTPS/TLS
- Implement rate limiting
- Add authentication middleware
- Use environment-specific configurations
- Regular security updates

### Data Protection
- Encrypt sensitive data at rest
- Use secure database connections
- Implement audit logging
- Regular backup procedures

## 📈 Performance

### Optimization Tips
- Use connection pooling for database
- Implement caching for templates
- Batch notification processing
- Monitor and optimize slow queries
- Use appropriate indexes

### Scaling
- Horizontal scaling with load balancer
- Database read replicas
- Message queue for high throughput
- CDN for static assets

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the troubleshooting section
- Review the API documentation

---

**Built with ❤️ using Go, Docker, and PostgreSQL** 