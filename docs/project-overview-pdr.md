# fbcli - Product Development Requirements

## Overview
CLI tool for managing Facebook Pages from the terminal. Written in Go, inspired by xurl.

## Problem Statement
No actively maintained Facebook CLI tool exists. Previous tools (fbcmd, facebook-cli) abandoned 10+ years ago.

## Target Users
- Developers managing Facebook Pages
- CI/CD pipelines needing automated posting
- Social media managers preferring terminal workflows

## Core Features (MVP)
- OAuth 2.0 authentication with token storage
- Create posts (text, photo, video, link, scheduled)
- List, read, and delete posts
- JSON output for scripting
- Environment variable support for CI/CD

## Tech Stack
- Go 1.24+ with Cobra CLI framework
- huandu/facebook/v2 SDK
- Facebook Graph API v24.0
- GoReleaser for distribution

## Out of Scope (Phase 2)
- Multi-page management
- Stories
- Multi-photo carousel
- Analytics/insights
- Webhook notifications
