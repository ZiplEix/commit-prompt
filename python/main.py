import time
import os
import subprocess
import logging
import pyperclip

def extract_text(soup):
    text_parts = []
    for element in soup.find("div", class_="markdown").descendants:
        if element.name in ["p", "pre"]:  # Ajouter un saut de ligne après ces éléments
            text_parts.append("\n")
        elif element.name == "span":  # Ajouter un espace après un span
            text_parts.append(" ")
        if isinstance(element, str):
            text_parts.append(element)

    message = "".join(text_parts).strip()

    index = message.find("Copier")
    if index != -1:
        message = message[index + len("Copier") + 1:]

    return message

def ask_gpt_commit_message(prompt: str) -> str:
    chrome_option = Options()

    chrome_option.add_argument("--no-sandbox")
    chrome_option.add_argument("--disable-dev-shm-usage")

    driver = uc.Chrome(service=Service(ChromeDriverManager().install()), options=chrome_option)

    driver.get("https://chat.com")

    try:
        input_element = input_element = WebDriverWait(driver, 10).until(
            EC.element_to_be_clickable((By.ID, "prompt-textarea"))
        )
        time.sleep(1)
        input_element.click()

        splited_prompt = prompt.split("\n")
        for line in splited_prompt:
            input_element.send_keys(line)
            ActionChains(driver).key_down(Keys.SHIFT).key_down(Keys.ENTER).key_up(Keys.SHIFT).key_up(Keys.ENTER).perform()

        # input_element.send_keys(prompt)
        input_element.send_keys(Keys.RETURN)
    except Exception as e:
        print(f"Impossible de trouver l'input\nError: {e}")
        driver.quit()
        os._exit(1)

    # input("Press Enter to continue...")

    try:
        response_element = WebDriverWait(driver, 10).until(
            EC.visibility_of_element_located((By.XPATH, "//div[@data-message-author-role='assistant']"))
        )
        time.sleep(5)
        response_element = WebDriverWait(driver, 10).until(
            EC.visibility_of_element_located((By.XPATH, "//div[@data-message-author-role='assistant']"))
        )
        response_html = response_element.get_attribute("outerHTML")

        soup = BeautifulSoup(response_html, "html.parser")
        response_text = extract_text(soup)

        return response_text
    except Exception as e:
        print(f"Impossible de trouver la réponse\nError: {e}")
        driver.quit()
        os._exit(1)

    driver.quit()

def get_git_diff() -> str:
    try:
        result = subprocess.run(
            ["git", "diff", "--cached"],
            capture_output=True,
            text=True,
            check=True
        )
        return result.stdout
    except subprocess.CalledProcessError as e:
        print(f"Erreur lors de l'exécution de git diff --cached: {e}")
        return None

def get_modified_files():
    try:
        result = subprocess.run(
            ['git', 'diff', '--cached', '--name-only'],
            capture_output=True,   # Capture la sortie standard et d'erreur
            text=True,             # Retourne la sortie sous forme de chaîne de caractères
            check=True             # Lance une exception si la commande échoue
        )
        # Divise la sortie par ligne et renvoie une liste des fichiers modifiés
        modified_files = result.stdout.splitlines()
        return modified_files
    except subprocess.CalledProcessError as e:
        print(f"Erreur lors de l'exécution de git diff --cached: {e}")
        return []


def build_prompt() -> str:
    modified_files = get_modified_files()
    if not modified_files:
        print("Aucun fichier modifié")
        return

    prompt = "Write me a commit message following the conventional commit format for this pending changes:\n\n"

    prompt += "Here are my modified files:\n"
    for file in modified_files:
        prompt += f"- {file}:\n"
        with open(file, "r") as f:
            prompt += f"{f.read()}\n\n"

    prompt += "and here is the result of the command git diff --cached:\n"
    prompt += get_git_diff()

    prompt += "\n\nYour answer should only contain the commit message and the the body of the commit without enything else."

    return prompt


if __name__ == "__main__":
    prompt = build_prompt()
    if not prompt:
        os._exit(0)

    commit_message = ask_gpt_commit_message(prompt)
    pyperclip.copy(commit_message)
    print(f"Commit message copied to clipboard:\n{commit_message}")
