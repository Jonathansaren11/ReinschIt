\# Project Proposal: ReinschIt



\## Overview

\*\*Target Audience:\*\* DIY mechanics and vehicle owners who want to diagnose 

and fix their own cars, trucks, and motorcycles without paying for a shop visit.



\*\*Problem:\*\* When a vehicle breaks down or acts up, most people have no idea 

where to start. Repair manuals are overwhelming, YouTube videos are hit or miss, 

and mechanics are expensive. People who want to fix things themselves have no 

guided, conversational tool to help them diagnose and solve the problem step by step.



\*\*Value Proposition:\*\* My app helps DIY vehicle owners to diagnose and fix 

their own cars, trucks, and motorcycles by walking them through the problem 

with an AI assistant that asks the right questions and gives clear, 

actionable repair guidance.



\---



\## Feature Scope



\### Must-Have (MVP)

1\. Start a new repair session by describing a vehicle problem in plain English

2\. AI asks follow-up questions and provides step-by-step diagnosis and fix instructions

3\. Save repair sessions so users can come back and continue where they left off

4\. View history of past repair sessions

5\. User authentication (register and login)



\### Deferred (Not Building Yet)

\- Parts cost estimator

\- Integration with parts stores (RockAuto, AutoZone)

\- Community forum or shared repairs

\- Mobile app

\- Vehicle VIN lookup and auto-fill

\- Photo or video upload for damage assessment

\- Mechanic marketplace (find a local shop if DIY fails)



\---



\## Pages



\### Public Pages

\*\*Home (/)\*\*

\- Purpose: Landing page explaining what the app does

\- How users arrive: Direct URL or search



\*\*Login (/login)\*\*

\- Purpose: Authenticate existing users

\- How users arrive: Clicking Log In from home or nav



\*\*Register (/register)\*\*

\- Purpose: Create a new account

\- How users arrive: Clicking Sign Up from home or nav



\### Authenticated Pages

\*\*Dashboard (/dashboard)\*\*

\- Purpose: Show all past repair sessions with status and vehicle info

\- How users arrive: After login or clicking logo



\*\*New Session (/sessions/new)\*\*

\- Purpose: Start a new repair conversation by entering vehicle info and problem description

\- How users arrive: Start New Repair button from dashboard



\*\*Session (/sessions/:id)\*\*

\- Purpose: The active AI repair chat

\- How users arrive: Clicking a session from dashboard or after creating a new one



\---



\## Navigation Flow



Home

├── Login → Dashboard

└── Register → Dashboard



Dashboard

├── New Session → Session/:id

└── Click existing session → Session/:id



Session/:id

└── Back → Dashboard



\---



\## Database Schema



\### Table: users

| Column | Type | Notes |

|--------|------|-------|

| id | INT PRIMARY KEY AUTO\_INCREMENT | |

| email | VARCHAR(255) UNIQUE NOT NULL | |

| password\_hash | VARCHAR(255) NOT NULL | |

| name | VARCHAR(100) NOT NULL | |

| created\_at | TIMESTAMP DEFAULT NOW() | |



\### Table: sessions

| Column | Type | Notes |

|--------|------|-------|

| id | INT PRIMARY KEY AUTO\_INCREMENT | |

| user\_id | INT NOT NULL | Foreign key → users.id |

| vehicle\_year | INT | e.g. 2018 |

| vehicle\_make | VARCHAR(100) | e.g. Toyota |

| vehicle\_model | VARCHAR(100) | e.g. Camry |

| problem\_description | TEXT NOT NULL | User's initial description |

| status | ENUM('active','resolved') DEFAULT 'active' | |

| created\_at | TIMESTAMP DEFAULT NOW() | |

| updated\_at | TIMESTAMP DEFAULT NOW() | |



\### Table: messages

| Column | Type | Notes |

|--------|------|-------|

| id | INT PRIMARY KEY AUTO\_INCREMENT | |

| session\_id | INT NOT NULL | Foreign key → sessions.id |

| role | ENUM('user','assistant') NOT NULL | |

| content | TEXT NOT NULL | |

| created\_at | TIMESTAMP DEFAULT NOW() | |



\---



\## AI Review Decisions



\### Suggestion: Add pagination to session list

\- \*\*Decision:\*\* Accepted

\- \*\*Reason:\*\* A user could have many sessions over time. Loading all at once will be slow.



\### Suggestion: Add a parts cost estimator

\- \*\*Decision:\*\* Deferred

\- \*\*Reason:\*\* Scope creep for MVP. Can be added in a later version.



\### Suggestion: Add password reset

\- \*\*Decision:\*\* Deferred

\- \*\*Reason:\*\* Auth is a must-have but password reset can wait for v2.

