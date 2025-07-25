import os
import shutil
import utils

go_deps = [
    "github.com/a-h/templ/cmd/templ@v0.3.865",
    "github.com/air-verse/air@latest",
    "github.com/go-task/task/v3/cmd/task@latest"
]

def setup():
    print("Installing Go dependencies...\n")
    for dep in go_deps:
        utils.install_go_dependency(dep)  
    print("\nAll Go dependencies installed successfully.")
    
    if not os.path.exists(".env"):
        shutil.copy(".env.example", ".env")
        print(".env file created from .env.example.")
    else:
        print(".env file already exists.")
        
    if not os.path.exists("database/sqlite.db"):
        print("Creating SQLite database...")
        with open("database/sqlite.db", "w") as db_file:
            pass
        print("SQLite database created at database/sqlite.db.")
    else:
        print("SQLite database already exists at database/sqlite.db.")

if __name__ == "__main__":
    setup()