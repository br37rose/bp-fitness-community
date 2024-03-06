const twilio = require('twilio');

class TwilioService {
    constructor(accountSid, authToken, fromNumber, toNumber) {
        this.client = twilio(accountSid, authToken);
        this.fromNumber = fromNumber;
        this.toNumber = toNumber;
    }

    sendMessage(body, callback) {
        this.client.messages.create({
            body,
            from: this.fromNumber,
            to: this.toNumber
        }, callback);
    }

    // listCalls(callback, pageSize = 7) {
    //     let i = 0;
    //     this.client.calls.each({
    //         pageSize,
    //         callback: (call, done) => {
    //             callback(call);
    //             i++;
    //             if (i === 10) {
    //                 done();
    //             }
    //         },
    //     });
    // }

    // listMessages(callback) {
    //     this.client.messages.list((err, messages) => {
    //         callback(messages);
    //     });
    // }

    // createTrunk(friendlyName, callback) {
    //     this.client.trunking.v1.trunks.create({ friendlyName }, callback);
    // }

    // fetchTrunk(trunkSid, callback) {
    //     this.client.trunking.v1.trunks(trunkSid).fetch(callback);
    // }

    // updateTrunk(trunkSid, updateParams, callback) {
    //     this.client.trunking.v1.trunks(trunkSid).update(updateParams, callback);
    // }
}

module.exports = TwilioService;
