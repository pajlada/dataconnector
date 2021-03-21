# dataconnector

## To run

### Update your domain name and JWT secret
#### Create a new script file in the same Sheet where the Add-on is and put the following function:
```
// head -c 1000 /dev/urandom | tr -dc 'a-zA-Z0-9-' | fold -w 32 | head -n 1
function updateEnvVariables(){
  PropertiesService.getScriptProperties().setProperty('DOMAIN', 'https://api.example.com');
  PropertiesService.getScriptProperties().setProperty('JWT_SECRET', 'my_jwt_secret');
}
```
#### Now, run the `updateEnvVariables` function above with the "Play" icon.

### To deploy the Google Sheets Add-on
npm run deploy
