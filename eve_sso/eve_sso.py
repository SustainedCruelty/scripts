import requests
import base64
import clipboard
from http.server import HTTPServer, BaseHTTPRequestHandler
import threading

client_id = input("[?] Enter your application's client_id: ")
secret_key = input("[?] Enter your application's secret key: ")
scopes = input("[?] Enter the scopes you want to pull (seperated by whitespace): ").split(' ')

auth_header_bytes = (f"{client_id}:{secret_key}").encode('ascii')
auth_header_b64_bytes = base64.b64encode(auth_header_bytes)
auth_header = auth_header_b64_bytes.decode('ascii')

logon_url = f"https://login.eveonline.com/oauth/authorize?response_type=code&redirect_uri=http://localhost/oauth-callback&client_id={client_id}&scope={'+'.join(scopes)}"
clipboard.copy(logon_url)

print(f"\n[+] Login URL copied to your clipboard: {logon_url}\n")
print("[!] Visit the URL in your clipboard and log in with your account")

redirect_url = None

class CallbackHandler(BaseHTTPRequestHandler):
	
	def do_GET(self):
		self.send_response(200)
		self.send_header('content-type', 'text/html')
		self.end_headers()
		self.wfile.write('Thanks for authing; You may close this tab now.'.encode())
		threading.Thread(target=server.shutdown, daemon=True).start()
		global redirect_url
		redirect_url = self.path

server = HTTPServer(('', 80), CallbackHandler)
server.serve_forever()

auth_code = redirect_url.split('oauth-callback?code=')[-1]

headers = {"Content-Type": "application/json","Authorization": 'Basic 'f'{auth_header}',}
body = '{"grant_type":"authorization_code", "code":%s}' % f'"{auth_code}"'
response = requests.post('https://login.eveonline.com/oauth/token', headers=headers, data=body)

if response.status_code != 200:
    raise ValueError("AUTHENTICATION FAILED")

response = response.json()

access_token = response['access_token']
refresh_token = response['refresh_token']

print("\n=============================================================================\n")
print(f"ACCESS TOKEN: {access_token}")
print(f"REFRESH TOKEN: {refresh_token}")