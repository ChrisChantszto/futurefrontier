@echo off
echo Setting vm.max_map_count for Elasticsearch...
wsl -d docker-desktop sysctl -w vm.max_map_count=262144
echo Done! Elasticsearch should now work properly.
pause
