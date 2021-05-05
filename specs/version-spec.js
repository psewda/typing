const expect = require("chai").expect;
const utils = require("./utils");

describe("version suite", function () {
	it("verify version endpoint", async function () {
		let version = await utils.getVersion();

		expect(version.statusCode).to.equal(200);
		expect(version.body.version).to.include("Typing");
	});
});
