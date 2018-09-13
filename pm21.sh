#! /bin/bash

# Build web and other services
pm2 restart api scheduler streamserver web
