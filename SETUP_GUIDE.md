# Setup Guide: Real Email & Slack Notifications

This guide will help you configure the notification service to send real emails and Slack messages.

## Quick Start

1. **Run the startup script:**
```bash
./start.sh
```

2. **Follow the prompts to configure your credentials**

3. **Test your setup**

## Email Configuration

### Option 1: Gmail

1. **Enable 2-Factor Authentication** on your Google account
2. **Generate an App Password:**
   - Go to Google Account settings
   - Security → 2-Step Verification → App passwords
   - Generate a password for "Mail"
3. **Configure in .env:**
```bash
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-16-character-app-password
```

### Option 2: Outlook/Hotmail

1. **Configure in .env:**
```bash
EMAIL_HOST=smtp-mail.outlook.com
EMAIL_PORT=587
EMAIL_USERNAME=your-email@outlook.com
EMAIL_PASSWORD=your-regular-password
```

### Option 3: SendGrid

1. **Create a SendGrid account** at https://sendgrid.com
2. **Generate an API key**
3. **Configure in .env:**
```bash
EMAIL_HOST=smtp.sendgrid.net
EMAIL_PORT=587
EMAIL_USERNAME=apikey
EMAIL_PASSWORD=your-sendgrid-api-key
```

## Slack Configuration

### Step 1: Create a Slack App

1. **Go to https://api.slack.com/apps**
2. **Click "Create New App"**
3. **Choose "From scratch"**
4. **Name your app** (e.g., "Notification Service")
5. **Select your workspace**

### Step 2: Configure Bot Token Scopes

1. **Go to "OAuth & Permissions"** in the left sidebar
2. **Add the following Bot Token Scopes:**
   - `chat:write` - Send messages to channels
   - `chat:write.public` - Send messages to public channels
   - `channels:read` - View basic channel info
   - `groups:read` - View basic private channel info

### Step 3: Install App to Workspace

1. **Go to "Install App"** in the left sidebar
2. **Click "Install to Workspace"**
3. **Authorize the app**

### Step 4: Get Bot Token

1. **Copy the "Bot User OAuth Token"** (starts with `xoxb-`)
2. **Configure in .env:**
```bash
SLACK_TOKEN=xoxb-your-bot-token-here
SLACK_CHANNEL=#your-channel-name
```

### Step 5: Invite Bot to Channel

1. **In your Slack workspace, invite the bot to your channel:**
```
/invite @YourAppName
```

## Testing Your Setup

### Test Email

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "title": "Test Email",
    "message": "This is a test email from the notification service!",
    "recipient": "your-email@example.com"
  }'
```

### Test Slack

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "slack",
    "title": "Test Slack",
    "message": "This is a test message from the notification service!",
    "recipient": "#your-channel"
  }'
```

### Test In-App

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "in_app",
    "title": "Test In-App",
    "message": "This is a test in-app notification!",
    "recipient": "user123"
  }'
```

## Troubleshooting

### Email Issues

**"Authentication failed" error:**
- Check your username and password
- For Gmail: Make sure you're using an app password, not your regular password
- For Outlook: Make sure 2FA is not blocking the connection

**"Connection refused" error:**
- Check your EMAIL_HOST and EMAIL_PORT
- Verify your email provider's SMTP settings

### Slack Issues

**"Invalid token" error:**
- Make sure your SLACK_TOKEN starts with `xoxb-`
- Verify the token is from the correct app
- Check that the app is installed to your workspace

**"Channel not found" error:**
- Make sure the bot is invited to the channel
- Check the channel name format (should start with #)
- Verify the channel exists and is accessible

**"Missing scopes" error:**
- Add the required scopes to your Slack app
- Reinstall the app to your workspace

### General Issues

**Service won't start:**
```bash
# Check logs
docker-compose logs notification-service

# Restart services
docker-compose down
docker-compose up --build -d
```

**Environment variables not loading:**
- Make sure your .env file is in the project root
- Check that variable names match exactly
- Restart the service after changing .env

## Production Considerations

### Security

1. **Use strong JWT secrets**
2. **Store credentials securely** (use Docker secrets or external secret management)
3. **Use HTTPS** in production
4. **Limit API access** with authentication

### Email Providers

**Recommended for production:**
- **SendGrid** - Reliable, good deliverability
- **Mailgun** - Developer-friendly
- **Amazon SES** - Cost-effective for high volume

### Slack Best Practices

1. **Use different apps** for different environments
2. **Monitor rate limits** (Slack has API limits)
3. **Handle errors gracefully**
4. **Use appropriate channels** for different notification types

## Advanced Configuration

### Custom Email Templates

```bash
# Create a template
curl -X POST http://localhost:8080/api/v1/templates \
  -H "Content-Type: application/json" \
  -d '{
    "name": "welcome_email",
    "type": "email",
    "subject": "Welcome {{.Name}}!",
    "content": "Hello {{.Name}},\n\nWelcome to our platform!\n\nBest regards,\nThe Team",
    "variables": {"Name": "string"}
  }'

# Use template
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "type": "email",
    "title": "Welcome",
    "message": "Welcome message",
    "recipient": "user@example.com",
    "template_id": 1,
    "template_data": {"Name": "John Doe"}
  }'
```

### Scheduled Notifications

```bash
curl -X POST http://localhost:8080/api/v1/notifications/schedule \
  -H "Content-Type: application/json" \
  -d '{
    "type": "slack",
    "title": "Scheduled Alert",
    "message": "This will be sent later",
    "recipient": "#alerts",
    "scheduled_at": "2024-01-15T10:00:00Z"
  }'
```

## Support

If you encounter issues:

1. **Check the logs:** `docker-compose logs -f`
2. **Verify your configuration** in the .env file
3. **Test with simple examples** first
4. **Check the main README.md** for API documentation 