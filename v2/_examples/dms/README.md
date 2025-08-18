# Direct Messages Examples

This directory contains examples for using the Twitter API v2 Direct Messages endpoints.

## Examples

### Events
- **dm-events-lookup**: Retrieve DM events for the authenticated user

### Conversations  
- **dm-conversations-lookup**: Retrieve DM conversations for the authenticated user

### Manage
- **create-dm-conversation**: Create a new DM conversation
- **send-dm**: Send a message to an existing DM conversation

## OAuth2 Requirements

All DM examples require OAuth2 authentication with appropriate scopes:
- `dm.read` - For reading DM events and conversations
- `dm.write` - For creating conversations and sending messages

## Usage

Each example can be run with:
```bash
go run main.go
```

Make sure to set up your OAuth2 credentials and tokens before running the examples.