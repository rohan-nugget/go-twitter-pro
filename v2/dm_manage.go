package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// DMEvents retrieves DM events for the authenticated user
func (c *Client) DMEvents(ctx context.Context, opts DMEventOpts) (*DMEventResponse, error) {
	if opts.MaxResults > dmEventMaxResults {
		return nil, fmt.Errorf("dm events: max results [%d] is greater than max [%d]: %w", opts.MaxResults, dmEventMaxResults, ErrParameter)
	}
	
	ep := dmEventsEndpoint.url(c.Host)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("dm events request: %w", err)
	}
	
	req.Header.Set("Accept", "application/json")
	c.Authorizer.Add(req)
	
	u, err := url.Parse(req.URL.String())
	if err != nil {
		return nil, fmt.Errorf("dm events url parse: %w", err)
	}
	opts.addQuery(u)
	req.URL = u
	
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("dm events response: %w", err)
	}
	defer resp.Body.Close()
	
	decoder := json.NewDecoder(resp.Body)
	rl := rateFromHeader(resp.Header)
	
	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}
	
	respBody := &DMEventResponse{}
	if err := decoder.Decode(respBody); err != nil {
		return nil, fmt.Errorf("dm events decode: %w", err)
	}
	respBody.RateLimit = rl
	
	return respBody, nil
}

// DMConversationEvents retrieves DM events for a specific conversation
func (c *Client) DMConversationEvents(ctx context.Context, conversationID string, opts DMEventOpts) (*DMEventResponse, error) {
	if opts.MaxResults > dmEventMaxResults {
		return nil, fmt.Errorf("dm conversation events: max results [%d] is greater than max [%d]: %w", opts.MaxResults, dmEventMaxResults, ErrParameter)
	}
	
	ep := dmConversationEventsEndpoint.urlID(c.Host, conversationID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("dm conversation events request: %w", err)
	}
	
	req.Header.Set("Accept", "application/json")
	c.Authorizer.Add(req)
	
	u, err := url.Parse(req.URL.String())
	if err != nil {
		return nil, fmt.Errorf("dm conversation events url parse: %w", err)
	}
	opts.addQuery(u)
	req.URL = u
	
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("dm conversation events response: %w", err)
	}
	defer resp.Body.Close()
	
	decoder := json.NewDecoder(resp.Body)
	rl := rateFromHeader(resp.Header)
	
	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}
	
	respBody := &DMEventResponse{}
	if err := decoder.Decode(respBody); err != nil {
		return nil, fmt.Errorf("dm conversation events decode: %w", err)
	}
	respBody.RateLimit = rl
	
	return respBody, nil
}

// DMConversations retrieves DM conversations for the authenticated user
func (c *Client) DMConversations(ctx context.Context, opts DMConversationOpts) (*DMConversationResponse, error) {
	if opts.MaxResults > dmConversationMaxResults {
		return nil, fmt.Errorf("dm conversations: max results [%d] is greater than max [%d]: %w", opts.MaxResults, dmConversationMaxResults, ErrParameter)
	}
	
	ep := dmConversationsEndpoint.url(c.Host)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("dm conversations request: %w", err)
	}
	
	req.Header.Set("Accept", "application/json")
	c.Authorizer.Add(req)
	
	u, err := url.Parse(req.URL.String())
	if err != nil {
		return nil, fmt.Errorf("dm conversations url parse: %w", err)
	}
	opts.addQuery(u)
	req.URL = u
	
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("dm conversations response: %w", err)
	}
	defer resp.Body.Close()
	
	decoder := json.NewDecoder(resp.Body)
	rl := rateFromHeader(resp.Header)
	
	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}
	
	respBody := &DMConversationResponse{}
	if err := decoder.Decode(respBody); err != nil {
		return nil, fmt.Errorf("dm conversations decode: %w", err)
	}
	respBody.RateLimit = rl
	
	return respBody, nil
}

// DMConversationsByParticipant retrieves DM conversations by participant ID
func (c *Client) DMConversationsByParticipant(ctx context.Context, participantID string, opts DMConversationOpts) (*DMConversationResponse, error) {
	if opts.MaxResults > dmConversationMaxResults {
		return nil, fmt.Errorf("dm conversations by participant: max results [%d] is greater than max [%d]: %w", opts.MaxResults, dmConversationMaxResults, ErrParameter)
	}
	
	ep := dmConversationsByParticipantEndpoint.urlParticipantID(c.Host, participantID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("dm conversations by participant request: %w", err)
	}
	
	req.Header.Set("Accept", "application/json")
	c.Authorizer.Add(req)
	
	u, err := url.Parse(req.URL.String())
	if err != nil {
		return nil, fmt.Errorf("dm conversations by participant url parse: %w", err)
	}
	opts.addQuery(u)
	req.URL = u
	
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("dm conversations by participant response: %w", err)
	}
	defer resp.Body.Close()
	
	decoder := json.NewDecoder(resp.Body)
	rl := rateFromHeader(resp.Header)
	
	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}
	
	respBody := &DMConversationResponse{}
	if err := decoder.Decode(respBody); err != nil {
		return nil, fmt.Errorf("dm conversations by participant decode: %w", err)
	}
	respBody.RateLimit = rl
	
	return respBody, nil
}

// CreateDMConversation creates a new DM conversation
func (c *Client) CreateDMConversation(ctx context.Context, req CreateDMConversationRequest) (*CreateDMConversationResponse, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("create dm conversation marshal error: %w", err)
	}
	
	ep := dmConversationCreateEndpoint.url(c.Host)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create dm conversation request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	c.Authorizer.Add(httpReq)
	
	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("create dm conversation response: %w", err)
	}
	defer resp.Body.Close()
	
	decoder := json.NewDecoder(resp.Body)
	rl := rateFromHeader(resp.Header)
	
	if resp.StatusCode != http.StatusCreated {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}
	
	respBody := &CreateDMConversationResponse{}
	if err := decoder.Decode(respBody); err != nil {
		return nil, fmt.Errorf("create dm conversation decode: %w", err)
	}
	respBody.RateLimit = rl
	
	return respBody, nil
}

// SendDM sends a DM to an existing conversation
func (c *Client) SendDM(ctx context.Context, conversationID string, req SendDMRequest) (*SendDMResponse, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("send dm marshal error: %w", err)
	}
	
	ep := dmConversationEventsEndpoint.urlID(c.Host, conversationID)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("send dm request: %w", err)
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	c.Authorizer.Add(httpReq)
	
	resp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send dm response: %w", err)
	}
	defer resp.Body.Close()
	
	decoder := json.NewDecoder(resp.Body)
	rl := rateFromHeader(resp.Header)
	
	if resp.StatusCode != http.StatusCreated {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}
	
	respBody := &SendDMResponse{}
	if err := decoder.Decode(respBody); err != nil {
		return nil, fmt.Errorf("send dm decode: %w", err)
	}
	respBody.RateLimit = rl
	
	return respBody, nil
}

// validate validates the CreateDMConversationRequest
func (r CreateDMConversationRequest) validate() error {
	if len(r.ParticipantIDs) == 0 {
		return fmt.Errorf("create dm conversation: participant IDs are required: %w", ErrParameter)
	}
	if r.Text == "" && r.MediaID == "" {
		return fmt.Errorf("create dm conversation: either text or media ID is required: %w", ErrParameter)
	}
	return nil
}

// validate validates the SendDMRequest
func (r SendDMRequest) validate() error {
	if r.Text == "" && r.MediaID == "" {
		return fmt.Errorf("send dm: either text or media ID is required: %w", ErrParameter)
	}
	return nil
}