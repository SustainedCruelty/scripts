import requests
import concurrent.futures
from bs4 import BeautifulSoup
import time

def make_vote(url: str, oids: str, proxy: str, timeout: int):
	start = time.perf_counter()
	s = requests.Session()
	s.proxies = {"https":proxy,}
	headers = {
		'authority': 'strawpoll.de',
		'sec-ch-ua': '^\\^Google',
		'accept': '*/*',
		'content-type': 'application/x-www-form-urlencoded; charset=UTF-8',
		'x-requested-with': 'XMLHttpRequest',
		'sec-ch-ua-mobile': '?0',
		'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36',
		'sec-ch-ua-platform': '^\\^Windows^\\^',
		'origin': 'https://strawpoll.de',
		'sec-fetch-site': 'same-origin',
		'sec-fetch-mode': 'cors',
		'sec-fetch-dest': 'empty',
		'referer': url,
		'accept-language': 'de-DE,de;q=0.9,en-US;q=0.8,en;q=0.7',
	}

	data = {
	  'pid': url.split('/')[-1],
	  'oids': oids
	}

	try:
		response = s.post('https://strawpoll.de/vote', headers=headers, data=data, timeout = timeout)
		print(f"[+] Successfully casted vote with proxy {proxy} in {round(time.perf_counter()-start, 2)} seconds")
	except Exception:
		print(f"[-] Failed casting vote with proxy {proxy} (proxy failure / timeout);")	

def get_oids(url: str) -> dict:

	oids = {}
	r = requests.get(url)
	soup = BeautifulSoup(r.text, 'html.parser')
	
	poll = soup.find_all("div", class_="voteanswers")[0]
	for option in poll.find_all("div", class_="checkbox checkbox-danger"):
		id = option.find_all("input", class_="styled check checkvote")[0].get('name')
		name = option.find_all("label")[0].text.strip()
		oids[name] = id
		
	return oids
		
if __name__ == "__main__":

	url = input("[?] Enter the poll's URL: ")
	option = input("[?] What option to vote for: ")
	fname = input("[?] List of proxies to use: ")
	timeout = input("[?] How many seconds until timeout?: ")
	threads = input("[?] How many threads to use: ")
	
	with open(fname, 'r') as f:
		proxies = f.read().splitlines()
	print()
	print(f"[*] Proceeding to vote for option '{option}' with a total of {len(proxies)} proxies, {threads} threads and a timeout of {timeout} seconds")
	
	if input("[?] Do you want to continue? [y/n]: ") == 'y':
		print()
		oids = get_oids(url)
		with concurrent.futures.ThreadPoolExecutor(max_workers = int(threads)) as executor:
			for p in proxies:
				executor.submit(make_vote, url, oids[option], p, int(timeout))
	else:
		quit()