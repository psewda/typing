const got = require("got");

const googleApiUrl = "https://www.googleapis.com";
const baseUrl = "http://localhost:7070";

function httpOpts(url, method) {
	return {
		throwHttpErrors: false,
		responseType: "json",
		method: method,
		pathname: url,
	};
}

async function getVersion() {
	let opts = httpOpts("/api/version", "GET");
	return await send(baseUrl, opts);
}

async function getNote(id) {
	let opts = httpOpts("/api/v1/storage/notes/" + id, "GET");
	opts.headers = {
		Authorization: "Bearer " + global.accessToken,
	};

	return await send(baseUrl, opts);
}

async function getNotes() {
	let opts = httpOpts("/api/v1/storage/notes", "GET");
	opts.headers = {
		Authorization: "Bearer " + global.accessToken,
	};

	return await send(baseUrl, opts);
}

async function createNote(n) {
	let opts = httpOpts("/api/v1/storage/notes", "POST");
	opts.headers = {
		Authorization: "Bearer " + global.accessToken,
	};
	opts.json = n;

	return await send(baseUrl, opts);
}

async function updateNote(id, n) {
	let opts = httpOpts("/api/v1/storage/notes/" + id, "PUT");
	opts.headers = {
		Authorization: "Bearer " + global.accessToken,
	};
	opts.json = n;

	return await send(baseUrl, opts);
}

async function deleteNote(id) {
	let opts = httpOpts("/api/v1/storage/notes/" + id, "DELETE");
	opts.headers = {
		Authorization: "Bearer " + global.accessToken,
	};

	return await send(baseUrl, opts);
}

async function send(baseUrl, opts) {
	const res = await got(baseUrl, opts);
	return {
		statusCode: res.statusCode,
		body: res.body,
	};
}

module.exports = {
	googleApiUrl,
	baseUrl,
	httpOpts,
	getVersion,
	getNote,
	getNotes,
	createNote,
	updateNote,
	deleteNote,
};
