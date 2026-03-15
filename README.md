# 🛡️ Gin API Abuse Protection Middleware

A lightweight **security middleware for Gin (Go)** that protects APIs from abuse using **rate limiting, bot detection, IP reputation scoring, and temporary IP blocking**.

The middleware tracks behavior per IP and blocks clients that exhibit suspicious patterns like **missing user agents, excessive 404 requests, or rate limit violations**.

---

## 🚀 Features

### ✅ Implemented
- ☑ **IP-Based Rate Limiting** (using in-memory storage)
- ☑ **IP Reputation System**
- ☑ **Automatic IP Blocking** (temporary ban for abusive clients)
- ☑ **Bot Detection**
- ☑ **Geo Blocking** (using MaxMind GeoIP database)

### ⏳ Planned
- ☐ **Path-Based Rate Limiting**
- ☐ **Request Pattern Detection**
- ☐ **Request Fingerprinting**
- ☐ **Dashboard**
