import subprocess
import utils

def download_tailwind_binary():
    sys, url = utils.get_tailwind_binary_url()
    
    if sys == "Windows":
        try:
            subprocess.run(
                ["curl", "-L", url, "-o", "bin/tailwindcss.exe"],
                check=True
            )
            print("Tailwind CSS binary downloaded successfully.")
        except subprocess.CalledProcessError as e:
            print(f"Failed to download Tailwind CSS binary: {e}")
            raise
    elif sys in ["Linux", "Darwin"]:
        try:
            subprocess.run(
                ["curl", "-L", url, "-o", "bin/tailwindcss"],
                check=True
            )
            subprocess.run(
                ["chmod", "+x", "bin/tailwindcss"],
                check=True
            )
            print("Tailwind CSS binary downloaded and made executable successfully.")
        except subprocess.CalledProcessError as e:
            print(f"Failed to download or set permissions for Tailwind CSS binary: {e}")
            raise

if __name__ == "__main__":
    download_tailwind_binary()