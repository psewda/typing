const expect = require("chai").expect;
const utils = require("./utils");

describe("note suite", function () {
	this.timeout(5000);

	it("verify note list", async function () {
		// create multiple new notes
		let newNote = {
			name: "note",
			desc: "desc",
		};
		for (let i = 0; i < 2; i++) {
			let created = await utils.createNote(newNote);
			expect(created.statusCode).to.equal(201);
		}

		// fetch and verify notes
		let fetched = await utils.getNotes();
		expect(fetched.statusCode).to.equal(200);
		expect(fetched.body.length).to.have.at.least(2);
	});

	it("verify note creation", async function () {
		// create new note
		let newNote = {
			name: "note",
			desc: "desc",
		};
		let created = await utils.createNote(newNote);
		expect(created.statusCode).to.equal(201);
		expect(created.body.id).to.not.be.empty;
		expect(created.body.name).to.equal(newNote.name);
		expect(created.body.desc).to.equal(newNote.desc);

		// fetch and verify the note
		let fetched = await utils.getNote(created.body.id);
		expect(fetched.statusCode).to.equal(200);
		expect(fetched.body.name).to.equal(newNote.name);
		expect(fetched.body.desc).to.equal(newNote.desc);
	});

	it("verify note updation", async function () {
		// create new note
		let newNote = {
			name: "note",
			desc: "desc",
		};
		let created = await utils.createNote(newNote);
		expect(created.statusCode).to.equal(201);

		// update the note
		let upNote = {
			name: newNote.name,
			desc: newNote.desc,
			metadata: {
				meta1: "value1",
				meta2: "value2",
			},
			labels: ["label1", "label2"],
		};
		let updated = await utils.updateNote(created.body.id, upNote);
		expect(updated.statusCode).to.equal(200);
		expect(updated.body.name).to.equal(upNote.name);
		expect(updated.body.desc).to.equal(upNote.desc);
		expect(updated.body.metadata).to.deep.equal(upNote.metadata);
		expect(updated.body.labels).to.deep.equal(upNote.labels);

		// fetch and verify the note
		let fetched = await utils.getNote(created.body.id);
		expect(fetched.statusCode).to.equal(200);
		expect(fetched.body.name).to.equal(upNote.name);
		expect(fetched.body.desc).to.equal(upNote.desc);
		expect(fetched.body.metadata).to.deep.equal(upNote.metadata);
		expect(fetched.body.labels).to.deep.equal(upNote.labels);
	});

	it("verify note deletion", async function () {
		// create new note
		let newNote = {
			name: "note",
			desc: "desc",
		};
		let created = await utils.createNote(newNote);
		expect(created.statusCode).to.equal(201);

		// delete the note
		let deleted = await utils.deleteNote(created.body.id);
		expect(deleted.statusCode).to.equal(204);

		// fetch and verify the note
		let fetched = await utils.getNote(created.body.id);
		expect(fetched.statusCode).to.equal(404);
	});
});
