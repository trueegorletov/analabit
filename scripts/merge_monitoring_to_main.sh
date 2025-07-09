#!/bin/bash
# Script to merge the monitoring-setup branch into main after successful validation

set -e

echo "===== Analabit Monitoring Merge to Main ====="
echo "Starting merge process at $(date)"

# 1. Check current branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "monitoring-setup" ]; then
    echo "Error: You are not on the monitoring-setup branch. Please switch to it first."
    exit 1
fi

# 2. Check for any uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo "Error: You have uncommitted changes. Please commit or stash them first."
    exit 1
fi

# 3. Push any remaining changes to origin
echo "Pushing any remaining changes to origin..."
git push origin monitoring-setup

# 4. Check with the user before proceeding
echo ""
echo "IMPORTANT: Before merging to main, please confirm that:"
echo "- The monitoring setup has been deployed to production"
echo "- All monitoring services are working correctly"
echo "- The monitoring endpoints are properly secured"
echo "- The CI/CD workflow has been tested with the new changes"
echo ""
read -p "Have you completed all the validation steps? (y/n): " VALIDATED
if [ "$VALIDATED" != "y" ]; then
    echo "Merge aborted. Please complete the validation steps first."
    exit 1
fi

# 5. Fetch the latest changes from remote
echo "Fetching latest changes from remote..."
git fetch origin

# 6. Checkout main and pull latest changes
echo "Switching to main branch and pulling latest changes..."
git checkout main
git pull origin main

# 7. Merge the monitoring-setup branch
echo "Merging monitoring-setup branch into main..."
git merge --no-ff monitoring-setup -m "Merge monitoring-setup: Add Prometheus and Grafana monitoring stack"

# 8. Push the changes to origin
echo "Pushing merged changes to origin main..."
git push origin main

# 9. Confirm the merge
echo "===== Monitoring setup has been successfully merged to main at $(date) ====="
echo "The CI/CD pipeline should now deploy the changes to production."
echo "Please verify that the monitoring services continue to work after the CI/CD deployment."

# 10. Optional: Switch back to the monitoring-setup branch
git checkout monitoring-setup

echo "Done!"
