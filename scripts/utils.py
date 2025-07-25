import subprocess

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