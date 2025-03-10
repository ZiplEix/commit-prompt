# Clone repo
git clone git@github.com:ZiplEix/commit-prompt.git
cd commit-prompt

# go to the latest release
git checkout $(git describe --tags `git rev-list --tags --max-count=1`)

cp bin/commit-prompt /usr/local/bin/commit-prompt
chmod +x /usr/local/bin/commit-prompt

echo "Commit Prompt installed successfully in /usr/local/bin/commit-prompt"
