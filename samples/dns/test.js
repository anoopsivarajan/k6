import { sendUDP } from "k6/dns";
import { Trend } from "k6/metrics";
import { sleep } from "k6";
const server = "8.8.8.8:53";
let myTrend = new Trend("Duration");
const data = [
	{ domain: "bbc.co.uk", qType: "A" },
	{ domain: "google.com", qType: "AAAA" },
	{ domain: "mathrubhumi.com", qType: "MX" },
	{ domain: "manormaonline.com", qType: "AAAA" },
];

export default function () {
	let index = Math.round(Math.random(Math.random() * data.length) + 1);
	let domain = data[index];
	console.log(JSON.stringify(domain));
	// let resp = sendUDP(domain.domain, domain.qType, server);
	// console.log(JSON.stringify(resp));
	sleep(1);
	// myTrend.add(resp.duration, { c: "as" });
}
