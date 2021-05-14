const expect = require("chai").expect;
const got = require("got");
const utils = require("./utils");

exports.mochaGlobalSetup = async function () {
	// handle "unhandled promise rejection" error
	process.on("unhandledRejection", (error) => {
		console.log("unhandled rejection:", error);
		process.exit(1);
	});

	// refresh token
	accessToken = await refreshToken(process.env.REFRESH_TOKEN);
	console.log(
		"info: new access token '*****%s' received",
		accessToken.slice(-5)
	);

	// delete all files
	files = await getFiles(accessToken);
	console.log("info: found %d files in app data folder", files.length);
	if (files.length > 0) {
		for (const f of files) {
			console.log("info: deleting '%s' from app data folder", f.name);
			await deleteFile(f.id, accessToken);
		}
	}

	// set access token in global scope
	global.accessToken = accessToken;
};

async function refreshToken(refreshToken) {
	var opts = utils.httpOpts("/api/v1/signin/auth/refresh", "POST");
	opts.form = {
		refresh_token: refreshToken,
	};

	const res = await got(utils.baseUrl, opts);
	expect(res.statusCode).to.equal(200);
	return res.body.accessToken;
}

async function getFiles(accessToken) {
	var opts = utils.httpOpts("/drive/v3/files", "GET");
	opts.headers = {
		Authorization: "Bearer " + accessToken,
	};
	opts.searchParams = {
		spaces: "appDataFolder",
	};

	const res = await got(utils.googleApiUrl, opts);
	expect(res.statusCode).to.equal(200);
	return res.body.files.map(function (f) {
		return {
			id: f.id,
			name: f.name,
		};
	});
}

async function deleteFile(id, accessToken) {
	var opts = utils.httpOpts("/drive/v3/files/" + id, "DELETE");
	opts.headers = {
		Authorization: "Bearer " + accessToken,
	};

	const res = await got(utils.googleApiUrl, opts);
	expect(res.statusCode).to.equal(204);
}
