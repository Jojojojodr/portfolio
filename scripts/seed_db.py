import os
import sqlite3
import utils
import sys

def seed():
    if len(sys.argv) > 1:
        seed_file_path = os.path.join(utils.seeds_folder, sys.argv[1])
    else:
        print("No seed file specified, please provide a seed file name as an argument.")
        return
    
    conn = sqlite3.connect(utils.db_file)
    cur = conn.cursor()
        
    utils.seed_file(seed_file_path, cur)

    conn.commit()
    conn.close()
    
if __name__ == "__main__":
    seed()