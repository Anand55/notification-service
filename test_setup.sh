#!/bin/bash

# Test script for notification service setup

echo "🧪 Testing Notification Service Setup..."
echo ""

# Check if service is running
echo "1. Checking if service is running..."
if curl -s http://localhost:8080/health > /dev/null; then
    echo "✅ Service is running"
else
    echo "❌ Service is not running. Please start it first:"
    echo "   docker-compose up -d"
    exit 1
fi

echo ""
echo "2. Testing API endpoints..."

# Test health endpoint
echo "   • Health check..."
if curl -s http://localhost:8080/health | grep -q "ok"; then
    echo "   ✅ Health check passed"
else
    echo "   ❌ Health check failed"
fi

# Test email notification
echo "   • Testing email notification..."
EMAIL_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "title": "Test Email",
    "message": "This is a test email from the notification service!",
    "recipient": "test@example.com"
  }')

if echo "$EMAIL_RESPONSE" | grep -q "id"; then
    echo "   ✅ Email notification sent successfully"
    echo "   📧 Check your email or MailHog at http://localhost:8025"
else
    echo "   ❌ Email notification failed"
    echo "   💡 Make sure EMAIL_USERNAME and EMAIL_PASSWORD are set in .env"
fi

# Test Slack notification
echo "   • Testing Slack notification..."
SLACK_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "slack",
    "title": "Test Slack",
    "message": "This is a test message from the notification service!",
    "recipient": "#general"
  }')

if echo "$SLACK_RESPONSE" | grep -q "id"; then
    echo "   ✅ Slack notification sent successfully"
    echo "   💬 Check your Slack channel"
else
    echo "   ❌ Slack notification failed"
    echo "   💡 Make sure SLACK_TOKEN is set in .env and bot is invited to channel"
fi

# Test in-app notification
echo "   • Testing in-app notification..."
INAPP_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in_app",
    "title": "Test In-App",
    "message": "This is a test in-app notification!",
    "recipient": "user123"
  }')

if echo "$INAPP_RESPONSE" | grep -q "id"; then
    echo "   ✅ In-app notification sent successfully"
else
    echo "   ❌ In-app notification failed"
fi

echo ""
echo "3. Testing template functionality..."

# Test template creation
echo "   • Creating test template..."
TEMPLATE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test_template",
    "type": "email",
    "subject": "Hello {{.Name}}!",
    "content": "Hello {{.Name}},\n\nThis is a test template.\n\nBest regards,\nThe Team",
    "variables": {"Name": "string"}
  }')

if echo "$TEMPLATE_RESPONSE" | grep -q "id"; then
    echo "   ✅ Template created successfully"
    TEMPLATE_ID=$(echo "$TEMPLATE_RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
    
    # Test template usage
    echo "   • Testing template usage..."
    TEMPLATE_USAGE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/notifications \
      -H "Content-Type: application/json" \
      -d "{
        \"type\": \"email\",
        \"title\": \"Template Test\",
        \"message\": \"Template message\",
        \"recipient\": \"test@example.com\",
        \"template_id\": $TEMPLATE_ID,
        \"template_data\": {\"Name\": \"John Doe\"}
      }")
    
    if echo "$TEMPLATE_USAGE_RESPONSE" | grep -q "id"; then
        echo "   ✅ Template usage successful"
    else
        echo "   ❌ Template usage failed"
    fi
else
    echo "   ❌ Template creation failed"
fi

echo ""
echo "4. Testing scheduled notifications..."

# Test scheduled notification
echo "   • Creating scheduled notification..."
SCHEDULED_TIME=$(date -d "+2 minutes" -u +"%Y-%m-%dT%H:%M:%SZ")
SCHEDULED_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/notifications/schedule \
  -H "Content-Type: application/json" \
  -d "{
    \"type\": \"slack\",
    \"title\": \"Scheduled Test\",
    \"message\": \"This is a scheduled test message!\",
    \"recipient\": \"#general\",
    \"scheduled_at\": \"$SCHEDULED_TIME\"
  }")

if echo "$SCHEDULED_RESPONSE" | grep -q "id"; then
    echo "   ✅ Scheduled notification created successfully"
    echo "   ⏰ Will be sent at $SCHEDULED_TIME"
else
    echo "   ❌ Scheduled notification creation failed"
fi

echo ""
echo "🎉 Testing completed!"
echo ""
echo "📋 Summary:"
echo "   • Service is running at http://localhost:8080"
echo "   • API documentation available at http://localhost:8080/health"
echo ""
echo "🔧 Next steps:"
echo "   1. Check your email for test messages"
echo "   2. Check your Slack channel for test messages"
echo "   3. View logs: docker-compose logs -f"
echo "   4. Explore the API: curl http://localhost:8080/api/v1/notifications"
echo ""
echo "📖 For more information, see SETUP_GUIDE.md and README.md" 