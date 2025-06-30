import os
import sqlite3
import utils

def fresh():    
    create_db = not os.path.exists(utils.db_file)
        
    conn = sqlite3.connect(utils.db_file)
    cur = conn.cursor()

    if not create_db:
        utils.drop_all_tables(cur)
    utils.run_migrations(utils.migrations_folder, cur)
    utils.run_seeds(utils.seeds_folder, cur)
        
    conn.commit()
    conn.close()
    
    if create_db:
        print(f"Database created")
    else:
        print(f"Database already exists. Fresh migrations and seeds applied.")

if __name__ == "__main__":
    fresh()