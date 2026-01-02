# Check for required environment variables
: "${GIT_USER_EMAIL:?Need to set GIT_USER_EMAIL}"
: "${GIT_USER_NAME:?Need to set GIT_USER_NAME}"

# Git config
chmod 600 /root/.ssh/id_rsa
eval $(ssh-agent -s)
ssh-add /root/.ssh/id_rsa

echo "chmod 600 /root/.ssh/id_rsa" >>~/.bashrc
echo "eval \$(ssh-agent -s)" >>~/.bashrc
echo "ssh-add /root/.ssh/id_rsa" >>~/.bashrc

git config --global --add safe.directory /workspaces

git config --global user.email "$GIT_USER_EMAIL" &&
  git config --global user.name "$GIT_USER_NAME"
