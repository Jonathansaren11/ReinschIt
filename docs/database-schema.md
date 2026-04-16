\# Database Schema: ReinschIt



\## Table: users

| Column | Type | Notes |

|--------|------|-------|

| id | INT PRIMARY KEY AUTO\_INCREMENT | |

| email | VARCHAR(255) UNIQUE NOT NULL | |

| password\_hash | VARCHAR(255) NOT NULL | |

| name | VARCHAR(100) NOT NULL | |

| created\_at | TIMESTAMP DEFAULT NOW() | |



\## Table: sessions

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



\## Table: messages

| Column | Type | Notes |

|--------|------|-------|

| id | INT PRIMARY KEY AUTO\_INCREMENT | |

| session\_id | INT NOT NULL | Foreign key → sessions.id |

| role | ENUM('user','assistant') NOT NULL | |

| content | TEXT NOT NULL | |

| created\_at | TIMESTAMP DEFAULT NOW() | |



\## Relationships

\- users → sessions: One user has many sessions (user\_id foreign key)

\- sessions → messages: One session has many messages (session\_id foreign key)

