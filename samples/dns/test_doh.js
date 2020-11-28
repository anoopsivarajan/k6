import { sleep } from "k6";
import { buildMessage } from "k6/dns";
import http from "k6/http";

const data = [
	{ domain: "bbc.co.uk", qType: "A" },
	{ domain: "google.com", qType: "AAAA" },
	{ domain: "mathrubhumi.com", qType: "MX" },
	{ domain: "manormaonline.com", qType: "AAAA" },
];

export default function () {
	let index = Math.round(Math.random(Math.random() * data.length) + 1);
	let domain = data[0];
	const currentValue = buildMessage(domain);
	console.log(currentValue);

	const url = "https://dns.google/dns-query";
	let headers = {
		"content-type": "application/dns-message",
		accept: "application/dns-message",
	};
	let response = http.post(url, currentValue, { headers: headers });
	let body = response;
	console.log(JSON.stringify(body));
	// console.log(JSON.stringify(unpackMessage(body)));
	// let resp = sendUDP(domain.domain, domain.qType, server);
	// console.log(JSON.stringify(resp));
	sleep(1);
}
