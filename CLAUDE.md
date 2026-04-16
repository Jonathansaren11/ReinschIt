\# FixIt AI - Claude Context



\## What This Project Is

An AI-powered repair assistant for DIY vehicle owners. Users describe a problem 

with their car, truck, or motorcycle and get conversational step-by-step 

diagnosis and repair guidance powered by the Anthropic Claude API.



\## Folder Structure

/frontend   - React + Vite + Tailwind + Shadcn UI

/backend    - Go + Gin REST API

/docs       - Project documentation



\## How to Run



\### Frontend

cd frontend

npm install

npm run dev  (runs on port 3000)



\### Backend

cd backend

go mod download

go run .  (runs on port 8080)



\## Architecture Decisions

\- JWT tokens for authentication stored in localStorage

\- Anthropic API called from the backend (never from frontend)

\- All /api/\* routes handled by Go backend

\- Frontend is static files served by Nginx in production

\- Vite dev proxy forwards /api/\* to localhost:8080



\## Database

MySQL with three tables: users, sessions, messages

See docs/database-schema.md for full schema



\## Key API Endpoints

POST /api/register - create account

POST /api/login - get auth token

GET /api/sessions - list user's repair sessions

POST /api/sessions - start new repair session

POST /api/sessions/:id/messages - send message and get AI response



\## Environment Variables (backend)

DB\_USER, DB\_PASSWORD, DB\_NAME, DB\_HOST - MySQL connection

ANTHROPIC\_API\_KEY - Claude API key

JWT\_SECRET - for signing tokens

