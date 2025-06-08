# AthenAI

## Prerequisites
1. [Deno](https://deno.land/#installation) - JavaScript/TypeScript runtime
   - Windows: Using PowerShell
     ```powershell
     irm https://deno.land/install.ps1 | iex
     ```
   - Verify installation:
     ```powershell
     deno --version
     ```

2. [PostgreSQL](https://www.postgresql.org/download/windows/)
   - Remember your password during installation
   - Default port: 5432

## Setup
1. Clone the repository
```powershell
git clone https://github.com/alejandro-albiol/AthenAI .git
cd BodyBuilderAI
```

2. Create `.env` file in the root directory
```properties
PG_PASSWORD=your_postgres_password
```

3. Create database
```sql
psql -U postgres
CREATE DATABASE AthenAI;
```

4. Run the application
```powershell
deno task dev
```

The app will be available at http://localhost:8000
