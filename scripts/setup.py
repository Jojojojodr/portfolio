import os
import shutil
import subprocess

go_deps = [
    "github.com/a-h/templ/cmd/templ@v0.3.1020",
    "github.com/air-verse/air@latest",
    "github.com/go-task/task/v3/cmd/task@latest"
]

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

def setup_dependencies():
    print("Installing Go dependencies...\n")
    for dep in go_deps:
        install_go_dependency(dep)

    subprocess.run(["npm", "i"], check=True)
    print("\nAll Go dependencies installed successfully.")
    
def copy_config():    
    if not os.path.exists("config.yaml"):
        shutil.copy("config.example.yaml", "config.yaml")
        print("config.yaml created from config.example.yaml.")
    else:
        print("config file already exists.")

def create_database():
    if not os.path.exists("database/sqlite.db"):
        print("Creating SQLite database...")
        with open("database/sqlite.db", "w") as db_file:
            db_file.write("")
        print("SQLite database created at database/sqlite.db.")
    else:
        print("SQLite database already exists at database/sqlite.db.")

if __name__ == "__main__":
    setup_dependencies()
    copy_config()
    create_database()
