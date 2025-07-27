#!/bin/bash

# Notification Service Startup Script

echo "🚀 Starting Notification Service with Real Email & Slack Support..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker Desktop and try again."
    exit 1
fi

# Check if Docker Compose is available
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose and try again."
    exit 1
fi

echo "✅ Docker is running"

# Check if .env file exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp env.example .env
    echo ""
    echo "⚠️  IMPORTANT: Please edit the .env file with your real credentials:"
    echo "   • EMAIL_USERNAME and EMAIL_PASSWORD for real emails"
    echo "   • SLACK_TOKEN for real Slack notifications"
    echo ""
    echo "📧 Email Setup Options:"
    echo "   • Gmail: Use app password (not regular password)"
    echo "   • Outlook: Use regular password"
    echo "   • SendGrid: Use API key"
    echo ""
    echo "💬 Slack Setup:"
    echo "   • Create app at https://api.slack.com/apps"
    echo "   • Add bot token (starts with xoxb-)"
    echo "   • Invite bot to your channel"
    echo ""
    read -p "Press Enter after you've configured your .env file..."
fi

# Build and start services
echo "🔨 Building and starting services..."
docker-compose up --build -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if services are running
echo "🔍 Checking service status..."
docker-compose ps

echo ""
echo "🎉 Services are starting up!"
echo ""
echo "📋 Service URLs:"
echo "   • Notification Service API: http://localhost:8080"
echo "   • Health Check: http://localhost:8080/health"
echo ""
echo "🧪 Test Commands:"
echo ""
echo "1. Health Check:"
echo "   curl http://localhost:8080/health"
echo ""
echo "2. Send Real Email:"
echo '   curl -X POST http://localhost:8080/api/v1/notifications \
     -H "Content-Type: application/json" \
     -d '"'"'{"type":"email","title":"Test Email","message":"Hello from real email!","recipient":"your-email@example.com"}'"'"
echo ""
echo "3. Send Real Slack Message:"
echo '   curl -X POST http://localhost:8080/api/v1/notifications \
     -H "Content-Type: application/json" \
     -d '"'"'{"type":"slack","title":"Test Slack","message":"Hello from real Slack!","recipient":"#your-channel"}'"'"
echo ""
echo "4. Send In-App Notification:"
echo '   curl -X POST http://localhost:8080/api/v1/notifications \
     -H "Content-Type: application/json" \
     -d '"'"'{"type":"in_app","title":"Test In-App","message":"Hello from in-app!","recipient":"user123"}'"'"
echo ""
echo "🔧 Useful Commands:"
echo "   • View logs: docker-compose logs -f"
echo "   • Stop services: docker-compose down"
echo "   • Restart services: docker-compose restart"
echo ""
echo "📖 For more information, see README.md" 