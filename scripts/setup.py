import os
import shutil
from sys import platform
import utils
import tailwind

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
        
    sys = platform.system()
    if sys == "Windows":
        tailwind_path = "bin/tailwindcss.exe"
    elif sys in ["Linux", "Darwin"]:
        tailwind_path = "bin/tailwindcss"
        
    if not os.path.exists(tailwind_path):
        tailwind.download_tailwind_binary()
        print(f"Tailwind CSS binary downloaded to {tailwind_path}.")
    else:
        print(f"Tailwind CSS binary already exists at {tailwind_path}.")

if __name__ == "__main__":
    setup()