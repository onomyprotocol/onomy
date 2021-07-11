/* eslint-disable */
import { Reader, util, configure, Writer } from "protobufjs/minimal";
import * as Long from "long";
export const protobufPackage = "onomyprotocol.onomy.market";
const baseMsgCancelOrder = {
    creator: "",
    port: "",
    channel: "",
    amountDenom: "",
    exchRateDenom: "",
    orderID: 0,
};
export const MsgCancelOrder = {
    encode(message, writer = Writer.create()) {
        if (message.creator !== "") {
            writer.uint32(10).string(message.creator);
        }
        if (message.port !== "") {
            writer.uint32(18).string(message.port);
        }
        if (message.channel !== "") {
            writer.uint32(26).string(message.channel);
        }
        if (message.amountDenom !== "") {
            writer.uint32(34).string(message.amountDenom);
        }
        if (message.exchRateDenom !== "") {
            writer.uint32(42).string(message.exchRateDenom);
        }
        if (message.orderID !== 0) {
            writer.uint32(48).int32(message.orderID);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMsgCancelOrder };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.creator = reader.string();
                    break;
                case 2:
                    message.port = reader.string();
                    break;
                case 3:
                    message.channel = reader.string();
                    break;
                case 4:
                    message.amountDenom = reader.string();
                    break;
                case 5:
                    message.exchRateDenom = reader.string();
                    break;
                case 6:
                    message.orderID = reader.int32();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseMsgCancelOrder };
        if (object.creator !== undefined && object.creator !== null) {
            message.creator = String(object.creator);
        }
        else {
            message.creator = "";
        }
        if (object.port !== undefined && object.port !== null) {
            message.port = String(object.port);
        }
        else {
            message.port = "";
        }
        if (object.channel !== undefined && object.channel !== null) {
            message.channel = String(object.channel);
        }
        else {
            message.channel = "";
        }
        if (object.amountDenom !== undefined && object.amountDenom !== null) {
            message.amountDenom = String(object.amountDenom);
        }
        else {
            message.amountDenom = "";
        }
        if (object.exchRateDenom !== undefined && object.exchRateDenom !== null) {
            message.exchRateDenom = String(object.exchRateDenom);
        }
        else {
            message.exchRateDenom = "";
        }
        if (object.orderID !== undefined && object.orderID !== null) {
            message.orderID = Number(object.orderID);
        }
        else {
            message.orderID = 0;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.creator !== undefined && (obj.creator = message.creator);
        message.port !== undefined && (obj.port = message.port);
        message.channel !== undefined && (obj.channel = message.channel);
        message.amountDenom !== undefined &&
            (obj.amountDenom = message.amountDenom);
        message.exchRateDenom !== undefined &&
            (obj.exchRateDenom = message.exchRateDenom);
        message.orderID !== undefined && (obj.orderID = message.orderID);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseMsgCancelOrder };
        if (object.creator !== undefined && object.creator !== null) {
            message.creator = object.creator;
        }
        else {
            message.creator = "";
        }
        if (object.port !== undefined && object.port !== null) {
            message.port = object.port;
        }
        else {
            message.port = "";
        }
        if (object.channel !== undefined && object.channel !== null) {
            message.channel = object.channel;
        }
        else {
            message.channel = "";
        }
        if (object.amountDenom !== undefined && object.amountDenom !== null) {
            message.amountDenom = object.amountDenom;
        }
        else {
            message.amountDenom = "";
        }
        if (object.exchRateDenom !== undefined && object.exchRateDenom !== null) {
            message.exchRateDenom = object.exchRateDenom;
        }
        else {
            message.exchRateDenom = "";
        }
        if (object.orderID !== undefined && object.orderID !== null) {
            message.orderID = object.orderID;
        }
        else {
            message.orderID = 0;
        }
        return message;
    },
};
const baseMsgCancelOrderResponse = {};
export const MsgCancelOrderResponse = {
    encode(_, writer = Writer.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMsgCancelOrderResponse };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(_) {
        const message = { ...baseMsgCancelOrderResponse };
        return message;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = { ...baseMsgCancelOrderResponse };
        return message;
    },
};
const baseMsgSendCreateOrder = {
    sender: "",
    port: "",
    channelID: "",
    timeoutTimestamp: 0,
    amountDenom: "",
    amount: 0,
    sourceCoin: "",
    targetCoin: "",
    exchRateDenom: "",
    exchRate: "",
};
export const MsgSendCreateOrder = {
    encode(message, writer = Writer.create()) {
        if (message.sender !== "") {
            writer.uint32(10).string(message.sender);
        }
        if (message.port !== "") {
            writer.uint32(18).string(message.port);
        }
        if (message.channelID !== "") {
            writer.uint32(26).string(message.channelID);
        }
        if (message.timeoutTimestamp !== 0) {
            writer.uint32(32).uint64(message.timeoutTimestamp);
        }
        if (message.amountDenom !== "") {
            writer.uint32(42).string(message.amountDenom);
        }
        if (message.amount !== 0) {
            writer.uint32(48).int32(message.amount);
        }
        if (message.sourceCoin !== "") {
            writer.uint32(58).string(message.sourceCoin);
        }
        if (message.targetCoin !== "") {
            writer.uint32(66).string(message.targetCoin);
        }
        if (message.exchRateDenom !== "") {
            writer.uint32(74).string(message.exchRateDenom);
        }
        if (message.exchRate !== "") {
            writer.uint32(82).string(message.exchRate);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMsgSendCreateOrder };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.string();
                    break;
                case 2:
                    message.port = reader.string();
                    break;
                case 3:
                    message.channelID = reader.string();
                    break;
                case 4:
                    message.timeoutTimestamp = longToNumber(reader.uint64());
                    break;
                case 5:
                    message.amountDenom = reader.string();
                    break;
                case 6:
                    message.amount = reader.int32();
                    break;
                case 7:
                    message.sourceCoin = reader.string();
                    break;
                case 8:
                    message.targetCoin = reader.string();
                    break;
                case 9:
                    message.exchRateDenom = reader.string();
                    break;
                case 10:
                    message.exchRate = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseMsgSendCreateOrder };
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = String(object.sender);
        }
        else {
            message.sender = "";
        }
        if (object.port !== undefined && object.port !== null) {
            message.port = String(object.port);
        }
        else {
            message.port = "";
        }
        if (object.channelID !== undefined && object.channelID !== null) {
            message.channelID = String(object.channelID);
        }
        else {
            message.channelID = "";
        }
        if (object.timeoutTimestamp !== undefined &&
            object.timeoutTimestamp !== null) {
            message.timeoutTimestamp = Number(object.timeoutTimestamp);
        }
        else {
            message.timeoutTimestamp = 0;
        }
        if (object.amountDenom !== undefined && object.amountDenom !== null) {
            message.amountDenom = String(object.amountDenom);
        }
        else {
            message.amountDenom = "";
        }
        if (object.amount !== undefined && object.amount !== null) {
            message.amount = Number(object.amount);
        }
        else {
            message.amount = 0;
        }
        if (object.sourceCoin !== undefined && object.sourceCoin !== null) {
            message.sourceCoin = String(object.sourceCoin);
        }
        else {
            message.sourceCoin = "";
        }
        if (object.targetCoin !== undefined && object.targetCoin !== null) {
            message.targetCoin = String(object.targetCoin);
        }
        else {
            message.targetCoin = "";
        }
        if (object.exchRateDenom !== undefined && object.exchRateDenom !== null) {
            message.exchRateDenom = String(object.exchRateDenom);
        }
        else {
            message.exchRateDenom = "";
        }
        if (object.exchRate !== undefined && object.exchRate !== null) {
            message.exchRate = String(object.exchRate);
        }
        else {
            message.exchRate = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined && (obj.sender = message.sender);
        message.port !== undefined && (obj.port = message.port);
        message.channelID !== undefined && (obj.channelID = message.channelID);
        message.timeoutTimestamp !== undefined &&
            (obj.timeoutTimestamp = message.timeoutTimestamp);
        message.amountDenom !== undefined &&
            (obj.amountDenom = message.amountDenom);
        message.amount !== undefined && (obj.amount = message.amount);
        message.sourceCoin !== undefined && (obj.sourceCoin = message.sourceCoin);
        message.targetCoin !== undefined && (obj.targetCoin = message.targetCoin);
        message.exchRateDenom !== undefined &&
            (obj.exchRateDenom = message.exchRateDenom);
        message.exchRate !== undefined && (obj.exchRate = message.exchRate);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseMsgSendCreateOrder };
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = object.sender;
        }
        else {
            message.sender = "";
        }
        if (object.port !== undefined && object.port !== null) {
            message.port = object.port;
        }
        else {
            message.port = "";
        }
        if (object.channelID !== undefined && object.channelID !== null) {
            message.channelID = object.channelID;
        }
        else {
            message.channelID = "";
        }
        if (object.timeoutTimestamp !== undefined &&
            object.timeoutTimestamp !== null) {
            message.timeoutTimestamp = object.timeoutTimestamp;
        }
        else {
            message.timeoutTimestamp = 0;
        }
        if (object.amountDenom !== undefined && object.amountDenom !== null) {
            message.amountDenom = object.amountDenom;
        }
        else {
            message.amountDenom = "";
        }
        if (object.amount !== undefined && object.amount !== null) {
            message.amount = object.amount;
        }
        else {
            message.amount = 0;
        }
        if (object.sourceCoin !== undefined && object.sourceCoin !== null) {
            message.sourceCoin = object.sourceCoin;
        }
        else {
            message.sourceCoin = "";
        }
        if (object.targetCoin !== undefined && object.targetCoin !== null) {
            message.targetCoin = object.targetCoin;
        }
        else {
            message.targetCoin = "";
        }
        if (object.exchRateDenom !== undefined && object.exchRateDenom !== null) {
            message.exchRateDenom = object.exchRateDenom;
        }
        else {
            message.exchRateDenom = "";
        }
        if (object.exchRate !== undefined && object.exchRate !== null) {
            message.exchRate = object.exchRate;
        }
        else {
            message.exchRate = "";
        }
        return message;
    },
};
const baseMsgSendCreateOrderResponse = {};
export const MsgSendCreateOrderResponse = {
    encode(_, writer = Writer.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseMsgSendCreateOrderResponse,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(_) {
        const message = {
            ...baseMsgSendCreateOrderResponse,
        };
        return message;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = {
            ...baseMsgSendCreateOrderResponse,
        };
        return message;
    },
};
const baseMsgSendCreatePair = {
    sender: "",
    port: "",
    channelID: "",
    timeoutTimestamp: 0,
    sourceDenom: "",
    targetDenom: "",
};
export const MsgSendCreatePair = {
    encode(message, writer = Writer.create()) {
        if (message.sender !== "") {
            writer.uint32(10).string(message.sender);
        }
        if (message.port !== "") {
            writer.uint32(18).string(message.port);
        }
        if (message.channelID !== "") {
            writer.uint32(26).string(message.channelID);
        }
        if (message.timeoutTimestamp !== 0) {
            writer.uint32(32).uint64(message.timeoutTimestamp);
        }
        if (message.sourceDenom !== "") {
            writer.uint32(42).string(message.sourceDenom);
        }
        if (message.targetDenom !== "") {
            writer.uint32(50).string(message.targetDenom);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMsgSendCreatePair };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.sender = reader.string();
                    break;
                case 2:
                    message.port = reader.string();
                    break;
                case 3:
                    message.channelID = reader.string();
                    break;
                case 4:
                    message.timeoutTimestamp = longToNumber(reader.uint64());
                    break;
                case 5:
                    message.sourceDenom = reader.string();
                    break;
                case 6:
                    message.targetDenom = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseMsgSendCreatePair };
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = String(object.sender);
        }
        else {
            message.sender = "";
        }
        if (object.port !== undefined && object.port !== null) {
            message.port = String(object.port);
        }
        else {
            message.port = "";
        }
        if (object.channelID !== undefined && object.channelID !== null) {
            message.channelID = String(object.channelID);
        }
        else {
            message.channelID = "";
        }
        if (object.timeoutTimestamp !== undefined &&
            object.timeoutTimestamp !== null) {
            message.timeoutTimestamp = Number(object.timeoutTimestamp);
        }
        else {
            message.timeoutTimestamp = 0;
        }
        if (object.sourceDenom !== undefined && object.sourceDenom !== null) {
            message.sourceDenom = String(object.sourceDenom);
        }
        else {
            message.sourceDenom = "";
        }
        if (object.targetDenom !== undefined && object.targetDenom !== null) {
            message.targetDenom = String(object.targetDenom);
        }
        else {
            message.targetDenom = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.sender !== undefined && (obj.sender = message.sender);
        message.port !== undefined && (obj.port = message.port);
        message.channelID !== undefined && (obj.channelID = message.channelID);
        message.timeoutTimestamp !== undefined &&
            (obj.timeoutTimestamp = message.timeoutTimestamp);
        message.sourceDenom !== undefined &&
            (obj.sourceDenom = message.sourceDenom);
        message.targetDenom !== undefined &&
            (obj.targetDenom = message.targetDenom);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseMsgSendCreatePair };
        if (object.sender !== undefined && object.sender !== null) {
            message.sender = object.sender;
        }
        else {
            message.sender = "";
        }
        if (object.port !== undefined && object.port !== null) {
            message.port = object.port;
        }
        else {
            message.port = "";
        }
        if (object.channelID !== undefined && object.channelID !== null) {
            message.channelID = object.channelID;
        }
        else {
            message.channelID = "";
        }
        if (object.timeoutTimestamp !== undefined &&
            object.timeoutTimestamp !== null) {
            message.timeoutTimestamp = object.timeoutTimestamp;
        }
        else {
            message.timeoutTimestamp = 0;
        }
        if (object.sourceDenom !== undefined && object.sourceDenom !== null) {
            message.sourceDenom = object.sourceDenom;
        }
        else {
            message.sourceDenom = "";
        }
        if (object.targetDenom !== undefined && object.targetDenom !== null) {
            message.targetDenom = object.targetDenom;
        }
        else {
            message.targetDenom = "";
        }
        return message;
    },
};
const baseMsgSendCreatePairResponse = {};
export const MsgSendCreatePairResponse = {
    encode(_, writer = Writer.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseMsgSendCreatePairResponse,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(_) {
        const message = {
            ...baseMsgSendCreatePairResponse,
        };
        return message;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = {
            ...baseMsgSendCreatePairResponse,
        };
        return message;
    },
};
export class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
    }
    CancelOrder(request) {
        const data = MsgCancelOrder.encode(request).finish();
        const promise = this.rpc.request("onomyprotocol.onomy.market.Msg", "CancelOrder", data);
        return promise.then((data) => MsgCancelOrderResponse.decode(new Reader(data)));
    }
    SendCreateOrder(request) {
        const data = MsgSendCreateOrder.encode(request).finish();
        const promise = this.rpc.request("onomyprotocol.onomy.market.Msg", "SendCreateOrder", data);
        return promise.then((data) => MsgSendCreateOrderResponse.decode(new Reader(data)));
    }
    SendCreatePair(request) {
        const data = MsgSendCreatePair.encode(request).finish();
        const promise = this.rpc.request("onomyprotocol.onomy.market.Msg", "SendCreatePair", data);
        return promise.then((data) => MsgSendCreatePairResponse.decode(new Reader(data)));
    }
}
var globalThis = (() => {
    if (typeof globalThis !== "undefined")
        return globalThis;
    if (typeof self !== "undefined")
        return self;
    if (typeof window !== "undefined")
        return window;
    if (typeof global !== "undefined")
        return global;
    throw "Unable to locate global object";
})();
function longToNumber(long) {
    if (long.gt(Number.MAX_SAFE_INTEGER)) {
        throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
    }
    return long.toNumber();
}
if (util.Long !== Long) {
    util.Long = Long;
    configure();
}
