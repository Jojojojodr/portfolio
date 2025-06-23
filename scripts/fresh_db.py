import os
import sqlite3
import utils

db_folder = os.path.join(os.path.dirname(__file__), '..', 'database')
db_file = os.path.join(db_folder, 'sqlite.db')
migrations_folder = os.path.join(db_folder, 'migrations')
seeds_folder = os.path.join(db_folder, 'seeds')

def fresh():    
    create_db = not os.path.exists(db_file)
        
    conn = sqlite3.connect(db_file)
    cur = conn.cursor()

    if not create_db:
        utils.drop_all_tables(cur)
    utils.run_migrations(migrations_folder, cur)
    utils.run_seeds(seeds_folder, cur)
        
    conn.commit()
    conn.close()
    
    if create_db:
        print(f"Database created")
    else:
        print(f"Database already exists. Fresh migrations and seeds applied.")

if __name__ == "__main__":
    fresh()