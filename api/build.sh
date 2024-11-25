#!/bin/bash
echo "Starting build process..."

mkdir -p /var/task/models
mkdir -p /var/task/templates
mkdir -p /var/task/static
mkdir -p /var/task/docs

cp -r ../models/* /var/task/models/
cp -r ../templates/* /var/task/templates/
cp -r ../static/* /var/task/static/
cp -r ../docs/* /var/task/docs/

echo "Build process completed"
