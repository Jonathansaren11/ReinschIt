\# Architecture: ReinschIt



\## Pages



\### Public

\- Home (/) — Landing page, explains the app, links to login and register

\- Login (/login) — Email and password login form

\- Register (/register) — Create a new account



\### Authenticated

\- Dashboard (/dashboard) — List of all repair sessions

\- New Session (/sessions/new) — Form to start a new repair

\- Session (/sessions/:id) — Active AI chat for a repair session



\## Navigation Flow



Home

├── Login → Dashboard

└── Register → Dashboard



Dashboard

├── New Session → Session/:id

└── Click existing session → Session/:id



Session/:id

└── Back → Dashboard



\## API Design



\### Auth

POST /api/register

\- Body: { name, email, password }

\- Response: { token, user }



POST /api/login

\- Body: { email, password }

\- Response: { token, user }



\### Sessions

GET /api/sessions

\- Headers: Authorization: Bearer <token>

\- Response: Array of sessions



POST /api/sessions

\- Body: { vehicle\_year, vehicle\_make, vehicle\_model, problem\_description }

\- Response: Created session



GET /api/sessions/:id

\- Response: Session with all messages



PUT /api/sessions/:id

\- Body: { status }

\- Response: Updated session



DELETE /api/sessions/:id

\- Response: Success confirmation



\### Messages

POST /api/sessions/:id/messages

\- Body: { content }

\- Response: { userMessage, aiMessage }

