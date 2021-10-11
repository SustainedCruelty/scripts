import requests
import json

def get_vote_key(code: str) -> str:

	headers = {
		'authority': 'www.menti.com',
		'sec-ch-ua': '"Chromium";v="93", " Not;A Brand";v="99"',
		'accept': 'application/json',
		'content-type': 'application/json; charset=UTF-8',
		'sec-ch-ua-mobile': '?0',
		'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36',
		'sec-ch-ua-platform': '"Windows"',
		'sec-fetch-site': 'same-origin',
		'sec-fetch-mode': 'cors',
		'sec-fetch-dest': 'empty',
		'referer': 'https://www.menti.com/',
		'accept-language': 'de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7',
	}

	r = requests.get(f"https://www.menti.com/core/vote-ids/{code}/series", headers=headers).json()
	
	return r

def get_public_key(vote_key: str) -> str:

	headers = {
		'authority': 'www.menti.com',
		'sec-ch-ua': '";Not A Brand";v="99", "Chromium";v="94"',
		'accept': 'application/json',
		'content-type': 'application/json; charset=UTF-8',
		'sec-ch-ua-mobile': '?0',
		'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36',
		'sec-ch-ua-platform': '"Windows"',
		'sec-fetch-site': 'same-origin',
		'sec-fetch-mode': 'cors',
		'sec-fetch-dest': 'empty',
		'referer': f'https://www.menti.com/{vote_key}',
		'accept-language': 'de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7',
	}

	r = requests.get(f"https://www.menti.com/core/vote-keys/{vote_key}/series", headers=headers)
	
	return r.json()

def get_identifier() -> str:
	headers = {
	"user-agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 10_0 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) Version/10.0 YaBrowser/17.4.3.195.10 Mobile/14A346 Safari/E7FBAF"
	}
	
	r = requests.post("https://www.menti.com/core/identifier", headers = headers)
	
	return r.json()

def make_vote(vote_key: str, public_key: str, identifier: str, words: str):

	headers = {
    'authority': 'www.menti.com',
    'sec-ch-ua': '";Not A Brand";v="99", "Chromium";v="94"',
    'accept': 'application/json',
    'content-type': 'application/json; charset=UTF-8',
    'sec-ch-ua-mobile': '?0',
    'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36',
    'sec-ch-ua-platform': '"Windows"',
    'origin': 'https://www.menti.com',
    'sec-fetch-site': 'same-origin',
    'sec-fetch-mode': 'cors',
    'sec-fetch-dest': 'empty',
    'referer': f'https://www.menti.com/{vote_key}',
	'x-identifier':identifier,
    'accept-language': 'de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7',
	}
	#vote = " ".join([option1, option2, option3])
	data = '{"question_type":"wordcloud","vote":"%s"}' % words
	
	response = requests.post(f'https://www.menti.com/core/votes/{public_key}', headers=headers, data=data)
	if response.status_code == 200:
		print(f"[+] Successfully added words to the cloud ({words})")
	else:
		print("[-] Failed adding words to the cloud")
		
def add_words(code: str, vote_key: str, public_key: str, words: str) -> None:

	
	identifier = get_identifier()['identifier']
	#print(f"\n[+] Identifier: {identifier}")
	
	make_vote(vote_key, public_key, identifier, words)
	
	
if __name__ == '__main__':
	
	
	code = input("[?] Enter the code: ").replace(' ', '')
	file = input("[?] Filename of the file with the words you want to add: ")
	
	vote_key = get_vote_key(code)['vote_key']
	print(f"\n[*] Vote Key: {vote_key}")
	public_key = get_public_key(vote_key)['questions'][0]['public_key']
	print(f"[*] Public Key: {public_key}\n")
	
	with open(file, 'r') as f:
		total_words = f.read().splitlines()
	for i in range(0, len(total_words)-2, 3):
		words = f"{total_words[i]} {total_words[i+1]} {total_words[i+2]}"
		add_words(code, vote_key, public_key, words)
	
	
	
