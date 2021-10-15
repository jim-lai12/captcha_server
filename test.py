import requests
import time
import random


charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-"

def random_string():
	string = ""
	n = random.randint(1,10)
	for i in range(n):
		string += random.choice(charset)
	return string

def book_take_report_test(key):
	for times in range(20):
		randomt = random.randint(200, 500)

		recaptcha2 = {
			"url": "https://2captcha.com/demo/recaptcha-v2" + random_string(),
			"websiteKey": "6LeIxboZAAAAAFQy7d8GPzgRZu2bV0GwKS8ue_cH",
			"isInvisible": True,
		}

		recaptcha3 = {
			"url": "https://2captcha.com/demo/recaptcha-v3" + random_string(),
			"websiteKey": "6LeIxboZAAAAAFQy7d8GPzgRZu2bV0GwKS8ue_cH",
			"minScore": 0.5
		}

		data2 = {
			"BookingTime": int(time.time()) + randomt,
			"Key": key,
			"Type": "recaptchav2",
			"Number": 100,
			"Cumulative": True,
			"TaskInfo": recaptcha2,
			# "Rate":""

		}
		data3 = {
			"BookingTime": int(time.time()) + randomt,
			"Key": key,
			"Type": "recaptchav3",
			"Number": 100,
			"Cumulative": True,
			"TaskInfo": recaptcha3,
			# "Rate":""

		}
		r1 = requests.post("http://127.0.0.1:8000/api/task/book", json=data2)
		print(r1.json())
		print(data2)
		print("\n\n")
		r2 = requests.post("http://127.0.0.1:8000/api/task/book", json=data3)
		print(r2.json())
		print(data3)

		time.sleep(randomt)

		for i in range(10):
			res = requests.get("http://localhost:8000/api/task/" + r1.json()["TaskId"])
			print(res.json())
			report = {
				"Answer": res.json()["Answer"],
				"Correct": True
			}
			res = requests.post("http://localhost:8000/api/report", json=report)
			print(res.text)

			print "\n\n\n"
			res = requests.get("http://localhost:8000/api/task/" + r2.json()["TaskId"])
			print(res.json())
			report = {
				"Answer": res.json()["Answer"],
				"Correct": True
			}
			res = requests.post("http://localhost:8000/api/report", json=report)
			print(res.text)


def account_test():
	for i in range(10):
		data = {

			"Account": "admin" + random_string(),
			"Password": "admin"
		}
		header = {"Content-Type": "application/x-www-form-urlencoded"}

		r = requests.post("http://127.0.0.1:8000/api/books", data, headers=header)
		print(r.text)
	return r.json()["apikey"]

key = account_test()
book_take_report_test(key)
















