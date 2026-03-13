# 🛡️ Gin API Abuse Protection Middleware

A lightweight **security middleware for Gin (Go)** that protects APIs from abuse using **rate limiting, bot detection, IP reputation scoring, and temporary IP blocking**.

The middleware tracks behavior per IP and blocks clients that exhibit suspicious patterns like **missing user agents, excessive 404 requests, or rate limit violations**.

---

## 🧠 Architecture

```
Client Request
      │
      ▼
Gin Middleware
      │
      ├── User-Agent Inspection
      ├── Rate Limiting
      ├── Reputation Scoring
      ├── 404 Detection
      └── GeoIP Check (only IN/US Traffic)
      │
      ▼
Decision Engine
      │
      ├── Allow Request
      ├── Rate Limited (429)
      └── Temporary Block (403)
```

---

# 🚀 Features

## ✅ Implemented

- ☑ **IP-Based Rate Limiting**
  - Uses Go `rate.Limiter`
  - Default: **1 request/second with burst of 5**
  - Returns **HTTP 429** when exceeded

- ☑ **IP Reputation System**
  - Each suspicious action increases a **reputation score**
  - High score results in automatic IP blocking

- ☑ **Temporary IP Blocking**
  - Suspicious IPs are blocked for **1 minute**
  - Automatically unblocked by background cleanup process

- ☑ **Bot Detection**
  - Detects suspicious **User-Agent strings**
  - Known patterns checked:
    - `curl`
    - `wget`
    - `bot`
    - `spider`
    - `crawler`
    - `python-requests`
    - `axios`
    - `httpclient`
    - `go-http-client`
    - `PostmanRuntime`

- ☑ **Missing User-Agent Detection**
  - Requests without a **User-Agent header** increase the reputation score
  - Helps detect **scripted or automated clients**

- ☑ **404 Abuse Detection**
  - Tracks repeated **404 Not Found requests**
  - If an IP triggers **more than 5 invalid routes**, it is considered scanning activity
  - The IP is temporarily blocked

- ☑ **Reputation Score Recovery**
  - Reputation score slowly decreases over time
  - Prevents permanent punishment for temporary spikes

- ☑ **Visitor Cleanup System**
  - Background goroutine cleans inactive visitors
  - Keeps memory usage low

---

## ⏳ Planned

- ☐ **Path-Based Rate Limiting**
  - Apply different limits per endpoint

- ☐ **Advanced Bot Detection**
  - Detect scraping tools and automated scanners

- ☐ **Request Pattern Detection**
  - Detect sequential ID scraping
  - Detect high 404 scanning patterns

- ☐ **Geo Blocking**
  - Block traffic by country using **MaxMind GeoIP**

- ☐ **Request Fingerprinting**
  - Fingerprint clients using headers + behavior patterns

---

## ⚙️ Example Limits

| Rule | Action |
|-----|------|
| Rate limit exceeded | `429 Too Many Requests` |
| Reputation score > 10 | `403 Forbidden` |
| More than 5 404 requests | Temporary IP block |

---

## 🛠 Tech Stack

- **Go**
- **Gin Web Framework**
- **golang.org/x/time/rate**
- **In-memory visitor tracking**

---

## 📌 Goal

The goal of this middleware is to provide **simple API abuse protection without requiring external services like Redis or API gateways**, making it ideal for:

- small APIs
- microservices
- internal services
- prototype security layers
