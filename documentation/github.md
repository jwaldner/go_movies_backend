# Git Notes

# Setup
git config --global user.name "jwaldner"
git config --global user.email "joejazzenator@gmail.com"

# new
echo "# goreportService" >> README.md
git init
git add README.md
git commit -m "first commit"
git branch -M main
git remote add origin https://github.com/jwaldner/goreportService.git
git push -u origin main

# existing
git remote add origin https://github.com/jwaldner/goreportService.g

# add
git add .

# commit
git commit -m "commit message"

# updates local with remote repo
git pull origin

# remove example
git rm -r one-of-the-directories 
git commit . -m "Remove directory"
git push origin <your-git-branch> (typically 'master', but not always

ssh -T jwaldner@github.com








