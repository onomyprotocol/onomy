/* eslint-disable */
import { Writer, Reader } from "protobufjs/minimal";
export const protobufPackage = "onomyprotocol.onomy.market";
const baseDenomTrace = { index: "", port: "", channel: "", origin: "" };
export const DenomTrace = {
    encode(message, writer = Writer.create()) {
        if (message.index !== "") {
            writer.uint32(18).string(message.index);
        }
        if (message.port !== "") {
            writer.uint32(26).string(message.port);
        }
        if (message.channel !== "") {
            writer.uint32(34).string(message.channel);
        }
        if (message.origin !== "") {
            writer.uint32(42).string(message.origin);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseDenomTrace };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 2:
                    message.index = reader.string();
                    break;
                case 3:
                    message.port = reader.string();
                    break;
                case 4:
                    message.channel = reader.string();
                    break;
                case 5:
                    message.origin = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseDenomTrace };
        if (object.index !== undefined && object.index !== null) {
            message.index = String(object.index);
        }
        else {
            message.index = "";
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
        if (object.origin !== undefined && object.origin !== null) {
            message.origin = String(object.origin);
        }
        else {
            message.origin = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.index !== undefined && (obj.index = message.index);
        message.port !== undefined && (obj.port = message.port);
        message.channel !== undefined && (obj.channel = message.channel);
        message.origin !== undefined && (obj.origin = message.origin);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseDenomTrace };
        if (object.index !== undefined && object.index !== null) {
            message.index = object.index;
        }
        else {
            message.index = "";
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
        if (object.origin !== undefined && object.origin !== null) {
            message.origin = object.origin;
        }
        else {
            message.origin = "";
        }
        return message;
    },
};
