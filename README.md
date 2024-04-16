# Starchasers Nginx Configuration Creator

This tool was created to simplify the management of the Nginx lifecycle on systems requiring frequent configuration updates.
It ensures that the changes are validated before being applied, and that the current configuration is
preserved in case of the server errors.

It was designed to work with an external service that provides continuous configuration updates (S2S), 
although the endpoints can be also accessed directly if necessary.

## Environment Variables

- `API_KEY` - The API key used to authenticate requests to the service (required).
- `NGINX_CONFIG_FILE_PATH` - The path to the Nginx configuration file, also used to access its parent directory for modifications. 
If not set, it will default to `/etc/nginx/nginx.conf`.
