import os
import subprocess
import sqlite3
import json

db_folder = os.path.join(os.path.dirname(__file__), '..', 'database')
db_file = os.path.join(db_folder, 'sqlite.db')
migrations_folder = os.path.join(db_folder, 'migrations')
seeds_folder = os.path.join(db_folder, 'seeds')

def install_go_dependency(package: str):
    print(f"Installing Go dependency: {package}...")
    
    try:
        subprocess.run(
            ["go", "install", package],
            check=True
        )
        print(f"Successfully installed {package}.")
    except subprocess.CalledProcessError as e:
        print(f"Failed to install {package}: {e}")
        raise
    
def run_migrations(migrations_folder: str, cur: sqlite3.Cursor):
    # Run all migrations in the migrations folder
    if os.path.isdir(migrations_folder):
        for filename in sorted(os.listdir(migrations_folder)):
            if filename.endswith('.sql'):
                migration_path = os.path.join(migrations_folder, filename)
                with open(migration_path, 'r', encoding='utf-8') as f:
                    sql = f.read()
                    cur.executescript(sql)
                    print(f"Migrated: {filename}")
                    
def run_seeds(seeds_folder: str, cur: sqlite3.Cursor):
    # Seed the table with initial data
    if os.path.isdir(seeds_folder):
        for filename in sorted(os.listdir(seeds_folder)):
            if filename.endswith('.json'):
                table_name = os.path.splitext(filename)[0]
                seed_path = os.path.join(seeds_folder, filename)
                with open(seed_path, 'r', encoding='utf-8') as f:
                    data = json.load(f)
                    if isinstance(data, list) and data:
                        keys = data[0].keys()
                        columns = ', '.join(keys)
                        placeholders = ', '.join(['?'] * len(keys))
                        for row in data:
                            if "password" in row:
                                row["password"] = hash_password(row["password"])
                            values = tuple(row[k] for k in keys)
                            cur.execute(f"INSERT INTO {table_name} ({columns}) VALUES ({placeholders})", values)
                        print(f"Seeded: {filename}")

def seed_file(file_path: str, cur: sqlite3.Cursor):
    if not os.path.isfile(file_path):
        print(f"File not found: {file_path}")
        return

    table_name = os.path.splitext(os.path.basename(file_path))[0]
    with open(file_path, 'r', encoding='utf-8') as f:
        data = json.load(f)
        if isinstance(data, list) and data:
            keys = data[0].keys()
            columns = ', '.join(keys)
            placeholders = ', '.join(['?'] * len(keys))
            for row in data:
                if "password" in row:
                    row["password"] = hash_password(row["password"])
                values = tuple(row[k] for k in keys)
                cur.execute(f"INSERT INTO {table_name} ({columns}) VALUES ({placeholders})", values)
            print(f"Seeded: {file_path}")
                        
def drop_all_tables(cur: sqlite3.Cursor):
    # Drop all tables in the database except sqlite_sequence
    cur.execute("SELECT name FROM sqlite_master WHERE type='table';")
    tables = cur.fetchall()
    for table_name in tables:
        if table_name[0] == "sqlite_sequence":
            continue
        cur.execute(f"DROP TABLE IF EXISTS {table_name[0]}")
        print(f"Dropped table: {table_name[0]}")
        
def hash_password(password: str) -> str:
    result = subprocess.run(
        ["go", "run", "cmd/hash/hashpass.go", password],
        capture_output=True,
        text=True,
        check=True
    )
    return result.stdout.strip()