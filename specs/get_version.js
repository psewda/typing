const chai = require("chai");
const chaiHttp = require("chai-http");
const utils = require("./utils");

chai.use(chaiHttp);

it("verify version endpoint", function (done) {
	chai
		.request(utils.baseUrl)
		.get("/api/version")
		.end(function (_, res) {
			chai.expect(res).to.have.status(200);
			chai.expect(res.body).to.have.a.property("version");
			chai.expect(res.body.version).to.include("Typing");
			done();
		});
});
