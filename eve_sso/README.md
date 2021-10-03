Short script made for testing purposes to pull access and refresh tokens allowing you to make authenticated esi requests.
It is required that you have registered your own EVE application. If you haven't done so already you can do it [here](https://developers.eveonline.com/applications/create).
Additionally, the script only works if your application has the default callback url (http://localhost/oauth-callback)-

### Run the script
```bash
python eve_sso.py
```

### Follow the instructions
```console
[?] Enter your applications client_id: <your application's client id>
[?] Enter your applications secret key: <your application's secret key>
[?] Enter the scopes you want to pull (seperated by whitespace): <your application's scopes>

[+] Logon URL copied to your clipboard

[!] Visit the URL in your clipboard and log in with your account
127.0.0.1 - - [03/Oct/2021 14:36:17] "GET /oauth-callback?code=<code> HTTP/1.1" 200 -

=============================================================================

ACCESS TOKEN: <your access token>
REFRESH TOKEN: <your refresh token>
```
## License
[MIT](https://choosealicense.com/licenses/mit/)
